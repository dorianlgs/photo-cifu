package services

import (
	"context"
	"net/http"
	"path/filepath"

	"github.com/cschleiden/go-workflows/backend"
	"github.com/cschleiden/go-workflows/backend/sqlite"
	"github.com/cschleiden/go-workflows/client"
	"github.com/dorianlgs/photo-cifu/tools"
	"github.com/dorianlgs/photo-cifu/workflow"
	"github.com/google/uuid"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func New(app *pocketbase.PocketBase, workflowDbName string) *client.Client {

	baseDir, _ := tools.InspectRuntime()
	workflowDBPath := filepath.Join(baseDir, "pb_data", workflowDbName)

	workflowBackend := sqlite.NewSqliteBackend(workflowDBPath, sqlite.WithBackendOptions(backend.WithLogger(app.Logger())))

	workflowClient := client.New(workflowBackend)

	ctx := context.Background()

	go workflow.RunWorker(ctx, workflowBackend, app)

	return workflowClient
}

func CreateWorkflow(e *core.RequestEvent, workflowClient *client.Client) error {

	instanceId := uuid.NewString()

	ctx := context.Background()

	_, err := workflowClient.CreateWorkflowInstance(ctx, client.WorkflowInstanceOptions{
		InstanceID: instanceId,
	}, workflow.Workflow1, "input-for-workflow")
	if err != nil {
		return e.BadRequestError("Could not start workflow", err)
	}

	return e.JSON(http.StatusOK, map[string]any{"success": true, "instanceId": instanceId})
}
