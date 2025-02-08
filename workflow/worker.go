package workflow

import (
	"context"

	"github.com/cschleiden/go-workflows/backend"
	"github.com/cschleiden/go-workflows/worker"
	"github.com/pocketbase/pocketbase"
)

func RunWorker(ctx context.Context, mb backend.Backend, pb *pocketbase.PocketBase) {
	w := worker.New(mb, nil)

	w.RegisterWorkflow(Workflow1)

	w.RegisterActivity(&activities{pb: pb})
	w.RegisterActivity(&activities{pb: pb})

	if err := w.Start(ctx); err != nil {
		panic("could not start worker")
	}
}
