package handlers

import (
	"io"
	"net/http"

	"github.com/dorianlgs/photo-cifu/pkg/container"
	"github.com/dorianlgs/photo-cifu/pkg/errors"
	"github.com/dorianlgs/photo-cifu/pkg/validation"
	"github.com/pocketbase/pocketbase/core"
)

// Handlers contains all HTTP handlers
type Handlers struct {
	container *container.Container
}

// New creates a new handlers instance
func New(container *container.Container) *Handlers {
	return &Handlers{container: container}
}

// CreateGallery handles gallery creation requests
func (h *Handlers) CreateGallery(e *core.RequestEvent) error {
	// Parse and validate request
	req := &validation.GalleryCreateRequest{
		Name:     e.Request.FormValue("name"),
		Location: e.Request.FormValue("location"),
	}

	// Get images zip file
	imagesFile, imagesHeader, err := e.Request.FormFile("imagesZip")
	if err != nil {
		return errors.HandleError(e, errors.BadRequest("Failed to read images zip file", err))
	}
	defer imagesFile.Close()
	req.ImagesZip = imagesHeader

	// Get thumbnail file
	thumbnailFile, thumbnailHeader, err := e.Request.FormFile("thumbnail")
	if err != nil {
		return errors.HandleError(e, errors.BadRequest("Failed to read thumbnail file", err))
	}
	defer thumbnailFile.Close()
	req.Thumbnail = thumbnailHeader

	// Validate request
	if err := req.Validate(); err != nil {
		return errors.HandleError(e, err)
	}

	// Read file data
	imagesData, err := io.ReadAll(imagesFile)
	if err != nil {
		return errors.HandleError(e, errors.InternalError("Failed to read images file", err))
	}

	thumbnailData, err := io.ReadAll(thumbnailFile)
	if err != nil {
		return errors.HandleError(e, errors.InternalError("Failed to read thumbnail file", err))
	}

	// Create gallery using service
	galleryID, err := h.container.Services.Gallery.CreateGallery(
		req.Name,
		req.Location,
		imagesData,
		thumbnailData,
		imagesHeader.Filename,
		thumbnailHeader.Filename,
	)
	if err != nil {
		return errors.HandleError(e, err)
	}

	return e.JSON(http.StatusOK, map[string]string{
		"gallery_id": galleryID,
		"message":    "Gallery created successfully",
	})
}

// CreateWorkflow handles workflow creation requests
func (h *Handlers) CreateWorkflow(e *core.RequestEvent) error {
	// Parse request
	info, err := e.RequestInfo()
	if err != nil {
		return errors.HandleError(e, errors.BadRequest("Failed to parse request", err))
	}

	req := &validation.WorkflowCreateRequest{
		WorkflowType: getStringFromBody(info.Body, "workflow_type"),
		Input:        info.Body["input"],
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return errors.HandleError(e, err)
	}

	// Create workflow using service
	instanceID, err := h.container.Services.Workflow.CreateWorkflow(req.WorkflowType, req.Input)
	if err != nil {
		return errors.HandleError(e, err)
	}

	return e.JSON(http.StatusOK, map[string]string{
		"instance_id": instanceID,
		"message":     "Workflow created successfully",
	})
}

// SendSignal handles workflow signal requests
func (h *Handlers) SendSignal(e *core.RequestEvent) error {
	// Parse request
	info, err := e.RequestInfo()
	if err != nil {
		return errors.HandleError(e, errors.BadRequest("Failed to parse request", err))
	}

	req := &validation.SignalRequest{
		InstanceID: getStringFromBody(info.Body, "instance_id"),
		SignalName: getStringFromBody(info.Body, "signal_name"),
		Data:       info.Body["data"],
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return errors.HandleError(e, err)
	}

	// Send signal using service
	err = h.container.Services.Signal.SendSignal(req.InstanceID, req.SignalName, req.Data)
	if err != nil {
		return errors.HandleError(e, err)
	}

	return e.JSON(http.StatusOK, map[string]string{
		"message": "Signal sent successfully",
	})
}

// UpdateSettings handles settings update requests
func (h *Handlers) UpdateSettings(e *core.RequestEvent) error {
	// Parse request
	info, err := e.RequestInfo()
	if err != nil {
		return errors.HandleError(e, errors.BadRequest("Failed to parse request", err))
	}

	// Update settings using service
	err = h.container.Services.Settings.UpdateSettings(info.Body)
	if err != nil {
		return errors.HandleError(e, err)
	}

	return e.JSON(http.StatusOK, map[string]string{
		"message": "Settings updated successfully",
	})
}

// Helper functions
func getStringFromBody(body map[string]any, key string) string {
	if value, ok := body[key].(string); ok {
		return value
	}
	return ""
}