package container

import (
	"context"
	"path/filepath"

	"github.com/cschleiden/go-workflows/backend"
	"github.com/cschleiden/go-workflows/backend/sqlite"
	"github.com/cschleiden/go-workflows/client"
	"github.com/dorianlgs/photo-cifu/pkg/config"
	"github.com/dorianlgs/photo-cifu/tools"
	"github.com/dorianlgs/photo-cifu/workflow"
	"github.com/pocketbase/pocketbase"
)

// Container holds all application dependencies
type Container struct {
	App            *pocketbase.PocketBase
	Config         *config.Config
	WorkflowClient *client.Client
	Services       *ServiceContainer
}

// ServiceContainer holds all service implementations
type ServiceContainer struct {
	Gallery  GalleryService
	Workflow WorkflowService
	Signal   SignalService
	Settings SettingsService
}

// New creates a new dependency injection container
func New(app *pocketbase.PocketBase) *Container {
	cfg := config.New()

	// Create workflow client
	workflowClient := createWorkflowClient(app, cfg.WorkflowDB.Name)

	// Create services
	services := &ServiceContainer{
		Gallery:  NewGalleryService(app, cfg),
		Workflow: NewWorkflowService(workflowClient),
		Signal:   NewSignalService(workflowClient),
		Settings: NewSettingsService(app),
	}

	return &Container{
		App:            app,
		Config:         cfg,
		WorkflowClient: workflowClient,
		Services:       services,
	}
}

func createWorkflowClient(app *pocketbase.PocketBase, workflowDbName string) *client.Client {
	baseDir, _ := tools.InspectRuntime()
	workflowDBPath := filepath.Join(baseDir, "pb_data", workflowDbName)

	workflowBackend := sqlite.NewSqliteBackend(workflowDBPath, sqlite.WithBackendOptions(backend.WithLogger(app.Logger())))
	workflowClient := client.New(workflowBackend)

	// Start workflow worker
	ctx := context.Background()
	go workflow.RunWorker(ctx, workflowBackend, app)

	return workflowClient
}

// Service interfaces for better testability
type GalleryService interface {
	CreateGallery(name, location string, imagesZip, thumbnail []byte, zipHeader, thumbHeader string) (string, error)
}

type WorkflowService interface {
	CreateWorkflow(workflowType string, input interface{}) (string, error)
}

type SignalService interface {
	SendSignal(instanceID, signalName string, data interface{}) error
}

type SettingsService interface {
	UpdateSettings(settings map[string]interface{}) error
}