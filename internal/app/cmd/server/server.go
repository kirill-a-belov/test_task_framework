package server

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/kirill-a-belov/test_task_framework/pkg/tracer"
)

func New(ctx context.Context) cobra.Command {
	ctx, span := tracer.Start(ctx, "internal.app.cmd.server.New")
	defer span.End()

	//server

	// runner Run (server)

	return cobra.Command{}
}
