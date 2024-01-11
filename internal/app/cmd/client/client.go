package client

import (
	"context"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/kirill-a-belov/test_task_framework/internal/app/client"
	"github.com/kirill-a-belov/test_task_framework/internal/app/client/pkg/config"
	"github.com/kirill-a-belov/test_task_framework/pkg/logger"
	"github.com/kirill-a-belov/test_task_framework/pkg/runner"
	"github.com/kirill-a-belov/test_task_framework/pkg/tracer"
)

func New(ctx context.Context) cobra.Command {
	ctx, span := tracer.Start(ctx, "internal.app.cmd.client.New")
	defer span.End()

	return cobra.Command{
		Use:   "client",
		Short: "TCP client",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, span := tracer.Start(ctx, "internal.app.cmd.client.New.Run")
			defer span.End()

			cfg := &config.Config{}
			if err := cfg.Load(ctx); err != nil {
				return errors.Wrap(err, "loading config")
			}

			runner.New(client.New(ctx, cfg), logger.New("client")).Run(ctx)

			return nil
		},
	}
}
