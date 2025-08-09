package container

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/cschleiden/go-workflows/client"
	"github.com/dorianlgs/photo-cifu/pkg/config"
	"github.com/dorianlgs/photo-cifu/pkg/errors"
	"github.com/dorianlgs/photo-cifu/workflow"
	"github.com/google/uuid"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

// GalleryServiceImpl implements GalleryService
type GalleryServiceImpl struct {
	app *pocketbase.PocketBase
	cfg *config.Config
}

func NewGalleryService(app *pocketbase.PocketBase, cfg *config.Config) GalleryService {
	return &GalleryServiceImpl{app: app, cfg: cfg}
}

func (s *GalleryServiceImpl) CreateGallery(name, location string, imagesZip, thumbnail []byte, zipHeader, thumbHeader string) (string, error) {
	// Validate file size
	if int64(len(imagesZip)) > s.cfg.Gallery.MaxFileSize {
		return "", errors.ValidationError("Images zip file exceeds maximum size limit", nil)
	}

	// Parse zip file
	zipReader, err := zip.NewReader(bytes.NewReader(imagesZip), int64(len(imagesZip)))
	if err != nil {
		return "", errors.BadRequest("Invalid zip file format", err)
	}

	// Check image count
	if len(zipReader.File) > s.cfg.Gallery.MaxImages {
		return "", errors.ValidationError(
			fmt.Sprintf("Gallery cannot contain more than %d images", s.cfg.Gallery.MaxImages),
			nil,
		)
	}

	imagesCollection, err := s.app.FindCollectionByNameOrId("images")
	if err != nil {
		return "", errors.InternalError("Failed to find images collection", err)
	}

	galleriesCollection, err := s.app.FindCollectionByNameOrId("galleries")
	if err != nil {
		return "", errors.InternalError("Failed to find galleries collection", err)
	}

	galleryRecord := core.NewRecord(galleriesCollection)
	var galleryID string

	transactErr := s.app.RunInTransaction(func(txApp core.App) error {
		var imageIDs []string

		// Process each image in the zip
		for _, file := range zipReader.File {
			imageID, err := s.processImageFile(txApp, imagesCollection, file)
			if err != nil {
				return fmt.Errorf("failed to process image %s: %w", file.Name, err)
			}
			imageIDs = append(imageIDs, imageID)
		}

		// Create gallery record
		galleryRecord.Set("name", name)
		galleryRecord.Set("location", location)
		galleryRecord.Set("images", imageIDs)

		// Set thumbnail
		thumbnailFile, err := filesystem.NewFileFromBytes(thumbnail, thumbHeader)
		if err != nil {
			return fmt.Errorf("failed to create thumbnail file: %w", err)
		}
		galleryRecord.Set("thumbnail", thumbnailFile)

		if err := txApp.Save(galleryRecord); err != nil {
			return fmt.Errorf("failed to save gallery: %w", err)
		}

		galleryID = galleryRecord.Id
		return nil
	})

	if transactErr != nil {
		return "", errors.InternalError("Failed to create gallery", transactErr)
	}

	return galleryID, nil
}

func (s *GalleryServiceImpl) processImageFile(txApp core.App, collection *core.Collection, file *zip.File) (string, error) {
	fileReader, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file in archive: %w", err)
	}
	defer fileReader.Close()

	// Read file data
	dataBuffer := new(bytes.Buffer)
	if _, err := io.Copy(dataBuffer, fileReader); err != nil {
		return "", fmt.Errorf("failed to read file data: %w", err)
	}

	// Create image record
	imageRecord := core.NewRecord(collection)
	imageRecord.Set("likes", 0)

	imageFile, err := filesystem.NewFileFromBytes(dataBuffer.Bytes(), file.Name)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	imageRecord.Set("image", imageFile)

	if err := txApp.Save(imageRecord); err != nil {
		return "", fmt.Errorf("failed to save image: %w", err)
	}

	return imageRecord.Id, nil
}

// WorkflowServiceImpl implements WorkflowService
type WorkflowServiceImpl struct {
	client *client.Client
}

func NewWorkflowService(client *client.Client) WorkflowService {
	return &WorkflowServiceImpl{client: client}
}

func (s *WorkflowServiceImpl) CreateWorkflow(workflowType string, input interface{}) (string, error) {
	instanceID := uuid.NewString()
	ctx := context.Background()

	switch workflowType {
	case "gallery_process":
		// Convert input to proper workflow input structure
		galleryInput, err := s.convertToGalleryInput(input)
		if err != nil {
			return "", errors.ValidationError("Invalid input for gallery processing workflow", err)
		}

		_, err = s.client.CreateWorkflowInstance(ctx, client.WorkflowInstanceOptions{
			InstanceID: instanceID,
		}, workflow.Workflow1, galleryInput)

		if err != nil {
			return "", errors.InternalError("Failed to create workflow instance", err)
		}
	default:
		return "", errors.ValidationError(fmt.Sprintf("Unknown workflow type: %s", workflowType), nil)
	}

	return instanceID, nil
}

func (s *WorkflowServiceImpl) convertToGalleryInput(input interface{}) (workflow.GalleryProcessingInput, error) {
	var galleryInput workflow.GalleryProcessingInput

	if inputMap, ok := input.(map[string]interface{}); ok {
		if galleryID, ok := inputMap["gallery_id"].(string); ok {
			galleryInput.GalleryID = galleryID
		}
		if galleryName, ok := inputMap["gallery_name"].(string); ok {
			galleryInput.GalleryName = galleryName
		}
		if userEmail, ok := inputMap["user_email"].(string); ok {
			galleryInput.UserEmail = userEmail
		}
	}

	if galleryInput.GalleryID == "" {
		return galleryInput, fmt.Errorf("gallery_id is required")
	}

	return galleryInput, nil
}

// SignalServiceImpl implements SignalService
type SignalServiceImpl struct {
	client *client.Client
}

func NewSignalService(client *client.Client) SignalService {
	return &SignalServiceImpl{client: client}
}

func (s *SignalServiceImpl) SendSignal(instanceID, signalName string, data interface{}) error {
	ctx := context.Background()

	err := s.client.SignalWorkflow(ctx, instanceID, signalName, data)
	if err != nil {
		return errors.InternalError("Failed to send workflow signal", err)
	}

	return nil
}

// SettingsServiceImpl implements SettingsService
type SettingsServiceImpl struct {
	app *pocketbase.PocketBase
}

func NewSettingsService(app *pocketbase.PocketBase) SettingsService {
	return &SettingsServiceImpl{app: app}
}

func (s *SettingsServiceImpl) UpdateSettings(settings map[string]interface{}) error {
	// TODO: Implement actual settings persistence
	// For now, this is a placeholder that validates the input
	
	if settings == nil {
		return errors.ValidationError("Settings data is required", nil)
	}

	// Here you would typically save to a settings collection or configuration file
	// For example:
	// settingsCollection, err := s.app.FindCollectionByNameOrId("settings")
	// if err != nil {
	//     return errors.InternalError("Failed to find settings collection", err)
	// }
	// ... implement settings persistence logic

	return nil
}