package services

import (
	"context"
	"net/http"

	"github.com/cschleiden/go-workflows/client"
	"github.com/pocketbase/pocketbase/core"
)

func SendSignal(e *core.RequestEvent, workflowClient *client.Client) error {

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
}
