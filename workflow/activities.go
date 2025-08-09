package workflow

import (
	"context"
	"fmt"
	"time"

	"github.com/cschleiden/go-workflows/activity"
	"github.com/pocketbase/pocketbase"
)

type activities struct {
	pb *pocketbase.PocketBase
}

func (act *activities) ProcessGalleryImages(ctx context.Context, galleryID string) (int, error) {
	logger := activity.Logger(ctx)
	logger.Info("Processing gallery images", "galleryID", galleryID)

	record, err := act.pb.FindRecordById("galleries", galleryID)
	if err != nil {
		logger.Error("Failed to find gallery", "galleryID", galleryID, "error", err.Error())
		return 0, fmt.Errorf("gallery not found: %s", galleryID)
	}

	logger.Info("Found gallery", "name", record.Get("name"), "location", record.Get("location"))

	// Simulate image processing work
	time.Sleep(10 * time.Second)

	// In a real implementation, this would process the gallery images
	// For now, return the number of images processed
	images := record.Get("images")
	if imageList, ok := images.([]interface{}); ok {
		return len(imageList), nil
	}

	return 0, nil
}

func (act *activities) SendNotificationEmail(ctx context.Context, galleryName, userEmail string) error {
	logger := activity.Logger(ctx)
	logger.Info("Sending notification email", "galleryName", galleryName, "userEmail", userEmail)

	// Simulate email sending
	time.Sleep(3 * time.Second)

	// In a real implementation, this would use the mailer service
	// to send a notification email about the gallery processing completion

	logger.Info("Notification email sent successfully")
	return nil
}
