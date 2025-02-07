package workflow

import (
	"github.com/cschleiden/go-workflows/workflow"
)

func Workflow1(ctx workflow.Context, input string) error {

	logger := workflow.Logger(ctx)

	logger.Info("WF Init", "input", input)

	r1, err := workflow.ExecuteActivity[int](ctx, workflow.DefaultActivityOptions, Activity1, 35, 12).Get(ctx)
	if err != nil {
		panic("error getting activity 1 result")
	}

	logger.Info("A1", "result", r1)

	r2, err := workflow.ExecuteActivity[int](ctx, workflow.DefaultActivityOptions, Activity2).Get(ctx)
	if err != nil {
		panic("error getting activity 1 result")
	}

	logger.Info("A2", "result", r2)

	return nil
}
