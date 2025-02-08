package workflow

import (
	"fmt"
	"time"

	"github.com/cschleiden/go-workflows/workflow"
)

func Workflow1(ctx workflow.Context, input string) error {

	logger := workflow.Logger(ctx)

	logger.Info("WF Init", "input", input)

	var a *activities

	r1, err := workflow.ExecuteActivity[int](ctx, workflow.ActivityOptions{
		RetryOptions: workflow.RetryOptions{
			MaxAttempts:        1,
			FirstRetryInterval: time.Second * 3,
			BackoffCoefficient: 2,
		}}, a.Activity1, 35, 12).Get(ctx)

	if err != nil {
		logger.Error("Error from Activity 1", "err", err)
		return fmt.Errorf("getting result from activity 1: %w", err)
	}

	logger.Info("A1", "result", r1)

	logger.Info("Waiting signal test")

	tctx, cancel := workflow.WithCancel(ctx)

	timerFired := false

	workflow.Select(ctx,
		workflow.Await(workflow.ScheduleTimer(tctx, 60*time.Second), func(ctx workflow.Context, f workflow.Future[any]) {
			if _, err := f.Get(ctx); err != nil {
				logger.Info("Timer canceled")
			} else {
				logger.Info("Timer fired")
				timerFired = true
			}
		}),
		workflow.Receive(workflow.NewSignalChannel[int](ctx, "test"), func(ctx workflow.Context, r int, ok bool) {
			logger.Info("Received signal:", "r", r)

			cancel()
		}),
	)

	if timerFired {
		r2, err := workflow.ExecuteActivity[int](ctx, workflow.DefaultActivityOptions, a.Activity2).Get(ctx)
		if err != nil {
			logger.Error("Error from Activity 2", "err", err)
			return fmt.Errorf("getting result from activity 2: %w", err)
		}

		logger.Info("A2", "result", r2)
	}

	logger.Info("Workflow finished")

	return nil
}
