package main

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/cschleiden/go-workflows/backend"
	"github.com/cschleiden/go-workflows/backend/sqlite"
	"github.com/cschleiden/go-workflows/client"
	"github.com/google/uuid"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/pocketbase/pocketbase/tools/hook"

	"github.com/shujink0/photo-cifu/ui"
	"github.com/shujink0/photo-cifu/workflow"
)

const StaticWildcardParam = "path"

func main() {
	app := pocketbase.New()

	// ---------------------------------------------------------------
	// Optional plugin flags:
	// ---------------------------------------------------------------

	var hooksDir string
	app.RootCmd.PersistentFlags().StringVar(
		&hooksDir,
		"hooksDir",
		"",
		"the directory with the JS app hooks",
	)

	var hooksWatch bool
	app.RootCmd.PersistentFlags().BoolVar(
		&hooksWatch,
		"hooksWatch",
		true,
		"auto restart the app on pb_hooks file change; it has no effect on Windows",
	)

	var hooksPool int
	app.RootCmd.PersistentFlags().IntVar(
		&hooksPool,
		"hooksPool",
		15,
		"the total prewarm goja.Runtime instances for the JS app hooks execution",
	)

	var migrationsDir string
	app.RootCmd.PersistentFlags().StringVar(
		&migrationsDir,
		"migrationsDir",
		"",
		"the directory with the user defined migrations",
	)

	var automigrate bool
	app.RootCmd.PersistentFlags().BoolVar(
		&automigrate,
		"automigrate",
		true,
		"enable/disable auto migrations",
	)

	var publicDir string
	app.RootCmd.PersistentFlags().StringVar(
		&publicDir,
		"publicDir",
		defaultPublicDir(),
		"the directory to serve static files",
	)

	var indexFallback bool
	app.RootCmd.PersistentFlags().BoolVar(
		&indexFallback,
		"indexFallback",
		true,
		"fallback the request to index.html on missing static path, e.g. when pretty urls are used with SPA",
	)

	app.RootCmd.ParseFlags(os.Args[1:])

	// ---------------------------------------------------------------
	// Plugins and hooks:
	// ---------------------------------------------------------------

	// load jsvm (pb_hooks and pb_migrations)
	// jsvm.MustRegister(app, jsvm.Config{
	// 	MigrationsDir: migrationsDir,
	// 	HooksDir:      hooksDir,
	// 	HooksWatch:    hooksWatch,
	// 	HooksPoolSize: hooksPool,
	// })

	// migrate command (with js templates)
	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		TemplateLang: migratecmd.TemplateLangJS,
		Automigrate:  automigrate,
		Dir:          migrationsDir,
	})

	app.OnServe().Bind(&hook.Handler[*core.ServeEvent]{
		Func: func(e *core.ServeEvent) error {

			if !e.Router.HasRoute(http.MethodGet, "/{path...}") {
				e.Router.GET("/{path...}", apis.Static(ui.DistDirFS, indexFallback)).
					Bind(apis.Gzip())
			}

			baseDir, _ := inspectRuntime()
			workflowDBPath := filepath.Join(baseDir, "pb_data", "wofkflow.db")

			workflowBackend := sqlite.NewSqliteBackend(workflowDBPath, sqlite.WithBackendOptions(backend.WithLogger(app.Logger())))

			workflowClient := client.New(workflowBackend)

			ctx := context.Background()

			go workflow.RunWorker(ctx, workflowBackend, app)

			e.Router.POST("/api/photocifu/settings", func(e *core.RequestEvent) error {
				return e.JSON(http.StatusOK, map[string]bool{"success": true})
			}).Bind(apis.RequireAuth())

			e.Router.POST("/api/photocifu/gallery/create", func(e *core.RequestEvent) error {

				name := e.Request.FormValue("name")

				location := e.Request.FormValue("location")

				mf, mh, err := e.Request.FormFile("imagesZip")

				if err != nil {
					return e.BadRequestError("Failed to read zip file", err)
				}

				tmf, tmh, err := e.Request.FormFile("thumbnail")

				if err != nil {
					return e.BadRequestError("Failed to read zip file", err)
				}

				archive, err := zip.NewReader(mf, int64(mh.Size))
				if err != nil {
					return e.BadRequestError("Failed to open the zip file", err)
				}

				var imagesNames []string

				for _, f := range archive.File {

					fileInArchive, err := f.Open()
					if err != nil {
						return e.BadRequestError("Failed to open the zip file "+f.Name, err)
					}

					collection, err := app.FindCollectionByNameOrId("images")
					if err != nil {
						return err
					}

					record := core.NewRecord(collection)

					record.Set("likes", 0)
					dataFile := new(bytes.Buffer)
					if _, err := io.Copy(dataFile, fileInArchive); err != nil {
						return e.BadRequestError("Failed to copy the file "+f.Name, err)
					}

					fileInArchive.Close()

					f2, _ := filesystem.NewFileFromBytes(dataFile.Bytes(), f.Name)

					record.Set("image", f2)

					err = app.Save(record)
					if err != nil {
						return e.BadRequestError("Failed to save the image "+f.Name, err)
					}

					imagesNames = append(imagesNames, string(record.Id))

				}

				collection, err := app.FindCollectionByNameOrId("galleries")
				if err != nil {
					return err
				}

				record1 := core.NewRecord(collection)

				record1.Set("name", name)
				record1.Set("location", location)

				thumbnailData, err := io.ReadAll(tmf)
				if err != nil {
					return e.BadRequestError("Failed to read thumbnail file", err)
				}

				f3, _ := filesystem.NewFileFromBytes(thumbnailData, tmh.Filename)

				record1.Set("images", imagesNames)

				record1.Set("thumbnail", f3)

				err = app.Save(record1)
				if err != nil {
					return e.BadRequestError("Failed to save the gallery "+name, err)
				}

				return e.JSON(http.StatusOK, map[string]any{"galleryId": record1.Id})
			}).Bind(apis.RequireAuth())

			e.Router.POST("/api/photocifu/signal/send", func(e *core.RequestEvent) error {

				// alternatively, read the body via the parsed request info
				info, err := e.RequestInfo()
				if err != nil {
					return e.BadRequestError("Failed to read request data", err)
				}

				InstanceID, ok := info.Body["instanceId"].(string)

				if !ok {
					return e.BadRequestError("Failed to read instanceId", nil)
				}

				ctx := context.Background()

				workflowClient.SignalWorkflow(ctx, InstanceID, "test", 42)

				return e.JSON(http.StatusOK, map[string]bool{"success": true})
			}).Bind(apis.RequireAuth())

			e.Router.POST("/api/photocifu/workflow/create", func(e *core.RequestEvent) error {

				instanceId := uuid.NewString()

				_, err := workflowClient.CreateWorkflowInstance(ctx, client.WorkflowInstanceOptions{
					InstanceID: instanceId,
				}, workflow.Workflow1, "input-for-workflow")
				if err != nil {
					return e.BadRequestError("Could not start workflow", err)
				}

				return e.JSON(http.StatusOK, map[string]any{"success": true, "instanceId": instanceId})
			}).Bind(apis.RequireAuth())

			return e.Next()
		},
		Priority: 999, // execute as latest as possible to allow users to provide their own route
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}

}

// the default pb_public dir location is relative to the executable
func defaultPublicDir() string {
	if strings.HasPrefix(os.Args[0], os.TempDir()) {
		// most likely ran with go run
		return "./pb_public"
	}

	return filepath.Join(os.Args[0], "../pb_public")
}

func inspectRuntime() (baseDir string, withGoRun bool) {
	if strings.HasPrefix(os.Args[0], os.TempDir()) {
		// probably ran with go run
		withGoRun = true
		baseDir, _ = os.Getwd()
	} else {
		// probably ran with go build
		withGoRun = false
		baseDir = filepath.Dir(os.Args[0])
	}
	return
}
