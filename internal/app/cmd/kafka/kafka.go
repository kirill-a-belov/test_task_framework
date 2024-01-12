package kafka

import (
	"context"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/kirill-a-belov/test_task_framework/internal/app/kafka"
	"github.com/kirill-a-belov/test_task_framework/internal/app/kafka/pkg/config"
	"github.com/kirill-a-belov/test_task_framework/pkg/logger"
	"github.com/kirill-a-belov/test_task_framework/pkg/runner"
	"github.com/kirill-a-belov/test_task_framework/pkg/tracer"
)

func New(ctx context.Context) cobra.Command {
	ctx, span := tracer.Start(ctx, "internal.app.cmd.kafka.New")
	defer span.End()

	return cobra.Command{
		Use:   "kafka",
		Short: "Kafka services set",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, span := tracer.Start(ctx, "internal.app.cmd.kafka.New.Run")
			defer span.End()

			cfg := &config.Config{}
			if err := cfg.Load(ctx); err != nil {
				return errors.Wrap(err, "loading config")
			}

			runner.New(kafka.New(ctx, cfg), logger.New("kafka")).Run(ctx)

			return nil
		},
	}
}
