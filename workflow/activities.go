package workflow

import (
	"context"
	"time"

	"github.com/cschleiden/go-workflows/activity"
)

func Activity1(ctx context.Context, a, b int) (int, error) {

	logger := activity.Logger(ctx)
	logger.Info("Entering Activity1")

	time.Sleep(10 * time.Second)

	return a + b, nil
}

func Activity2(ctx context.Context) (int, error) {

	logger := activity.Logger(ctx)
	logger.Info("Entering Activity2")

	time.Sleep(3 * time.Second)

	return 12, nil
}
