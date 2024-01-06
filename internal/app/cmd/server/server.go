package server

import (
	"context"
	"github.com/kirill-a-belov/test_task_framework/internal/app/server"
	"github.com/kirill-a-belov/test_task_framework/internal/app/server/pkg/config"
	"github.com/kirill-a-belov/test_task_framework/pkg/logger"
	"github.com/kirill-a-belov/test_task_framework/pkg/runner"
	"github.com/pkg/errors"

	"github.com/spf13/cobra"

	"github.com/kirill-a-belov/test_task_framework/pkg/tracer"
)

func New(ctx context.Context) cobra.Command {
	ctx, span := tracer.Start(ctx, "internal.app.cmd.server.New")
	defer span.End()

	return cobra.Command{
		Use:   "server",
		Short: "TCP server",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, span := tracer.Start(ctx, "internal.app.cmd.server.New.Run")
			defer span.End()

			cfg := &config.Config{}
			if err := cfg.Load(ctx); err != nil {
				return errors.Wrap(err, "loading config")
			}

			runner.New(server.New(ctx, cfg), logger.New("server")).Run(ctx)

			return nil
		},
	}
}
