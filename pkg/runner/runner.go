package runner

import (
	"context"
	"fmt"
	"github.com/kirill-a-belov/test_task_framework/pkg/logger"
	"github.com/kirill-a-belov/test_task_framework/pkg/tracer"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
)

type appController interface {
	Start(context.Context) error
	Stop(context.Context)
}

type Runner struct {
	app     appController
	log     logger.Logger
	sigChan chan os.Signal
}

func New(app appController, log logger.Logger) *Runner {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	return &Runner{
		app:     app,
		log:     log,
		sigChan: sigChan,
	}
}

func (r *Runner) Run(ctx context.Context) {
	ctx, span := tracer.Start(ctx, "pkg.runner.Run")
	defer span.End()

	defer func() {
		if pnc := recover(); r != nil {
			debug.PrintStack()
			r.log.Error(fmt.Errorf("panic while command executionn (%v)", pnc))
		}

		return
	}()

	if err := r.app.Start(ctx); err != nil {
		r.log.Error(err, "error while command execution")

		return
	}

	select {
	case <-r.sigChan:
		r.app.Stop(ctx)
	}

}
