package validation

import (
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/dorianlgs/photo-cifu/pkg/errors"
)

// GalleryCreateRequest represents gallery creation input
type GalleryCreateRequest struct {
	Name      string                `json:"name"`
	Location  string                `json:"location"`
	ImagesZip *multipart.FileHeader `json:"-"`
	Thumbnail *multipart.FileHeader `json:"-"`
}

// Validate validates the gallery creation request
func (r *GalleryCreateRequest) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return errors.ValidationError("Gallery name is required", nil)
	}

	if len(r.Name) > 100 {
		return errors.ValidationError("Gallery name must be less than 100 characters", nil)
	}

	if r.ImagesZip == nil {
		return errors.ValidationError("Images zip file is required", nil)
	}

	if r.Thumbnail == nil {
		return errors.ValidationError("Thumbnail image is required", nil)
	}

	// Validate file extensions
	if !isValidZipFile(r.ImagesZip.Filename) {
		return errors.ValidationError("Images file must be a zip archive", nil)
	}

	if !isValidImageFile(r.Thumbnail.Filename) {
		return errors.ValidationError("Thumbnail must be a valid image file", nil)
	}

	return nil
}

// WorkflowCreateRequest represents workflow creation input
type WorkflowCreateRequest struct {
	WorkflowType string      `json:"workflow_type"`
	Input        interface{} `json:"input"`
}

// Validate validates the workflow creation request
func (r *WorkflowCreateRequest) Validate() error {
	if strings.TrimSpace(r.WorkflowType) == "" {
		return errors.ValidationError("Workflow type is required", nil)
	}

	validTypes := []string{"gallery_process", "image_enhancement", "cleanup"}
	if !contains(validTypes, r.WorkflowType) {
		return errors.ValidationError(
			fmt.Sprintf("Invalid workflow type. Valid types: %s", strings.Join(validTypes, ", ")),
			nil,
		)
	}

	return nil
}

// SignalRequest represents workflow signal input
type SignalRequest struct {
	InstanceID string      `json:"instance_id"`
	SignalName string      `json:"signal_name"`
	Data       interface{} `json:"data"`
}

// Validate validates the signal request
func (r *SignalRequest) Validate() error {
	if strings.TrimSpace(r.InstanceID) == "" {
		return errors.ValidationError("Instance ID is required", nil)
	}

	if strings.TrimSpace(r.SignalName) == "" {
		return errors.ValidationError("Signal name is required", nil)
	}

	return nil
}

// Helper functions
func isValidZipFile(filename string) bool {
	return strings.HasSuffix(strings.ToLower(filename), ".zip")
}

func isValidImageFile(filename string) bool {
	validExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	lower := strings.ToLower(filename)
	for _, ext := range validExts {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}
	return false
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}