package workflow

import (
	"fmt"
	"time"

	"github.com/cschleiden/go-workflows/workflow"
)

// GalleryProcessingInput represents the input for gallery processing workflow
type GalleryProcessingInput struct {
	GalleryID   string `json:"gallery_id"`
	GalleryName string `json:"gallery_name"`
	UserEmail   string `json:"user_email"`
}

func Workflow1(ctx workflow.Context, input GalleryProcessingInput) error {
	logger := workflow.Logger(ctx)
	logger.Info("Starting gallery processing workflow", "galleryID", input.GalleryID, "galleryName", input.GalleryName)

	var a *activities

	// Process gallery images
	processedCount, err := workflow.ExecuteActivity[int](ctx, workflow.ActivityOptions{
		RetryOptions: workflow.RetryOptions{
			MaxAttempts:        3,
			FirstRetryInterval: time.Second * 5,
			BackoffCoefficient: 2,
		},
	}, a.ProcessGalleryImages, input.GalleryID).Get(ctx)

	if err != nil {
		logger.Error("Failed to process gallery images", "error", err)
		return fmt.Errorf("failed to process gallery images: %w", err)
	}

	logger.Info("Gallery images processed", "count", processedCount)

	// Wait for processing completion signal or timeout
	logger.Info("Waiting for processing completion signal")
	
	tctx, cancel := workflow.WithCancel(ctx)
	timerFired := false

	workflow.Select(ctx,
		workflow.Await(workflow.ScheduleTimer(tctx, 5*time.Minute), func(ctx workflow.Context, f workflow.Future[any]) {
			if _, err := f.Get(ctx); err != nil {
				logger.Info("Processing timer canceled")
			} else {
				logger.Info("Processing timeout reached")
				timerFired = true
			}
		}),
		workflow.Receive(workflow.NewSignalChannel[map[string]interface{}](ctx, "processing_complete"), func(ctx workflow.Context, data map[string]interface{}, ok bool) {
			logger.Info("Received processing completion signal", "data", data)
			cancel()
		}),
	)

	// Send notification email if timeout occurred
	if timerFired {
		_, err := workflow.ExecuteActivity[any](ctx, workflow.DefaultActivityOptions, a.SendNotificationEmail, input.GalleryName, input.UserEmail).Get(ctx)
		if err != nil {
			logger.Error("Failed to send notification email", "error", err)
			return fmt.Errorf("failed to send notification email: %w", err)
		}
		logger.Info("Notification email sent due to processing timeout")
	}

	logger.Info("Gallery processing workflow completed", "galleryID", input.GalleryID)
	return nil
}
