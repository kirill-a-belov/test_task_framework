package config

import (
	"context"
	"time"

	"github.com/caarlos0/env"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"

	"github.com/kirill-a-belov/test_task_framework/pkg/tracer"
)

type Config struct {
	Address string        `env:"KAFKA_ADDRESS" validate:"hostname_port"`
	Delay   time.Duration `env:"KAFKA_DELAY" validate:"gte=1ms,lte=1s"`
	ConnTTL time.Duration `env:"KAFKA_CONN_TTL" validate:"gte=1ms,lte=1s"`
}

func (c *Config) validate() error {
	return validator.New(validator.WithRequiredStructEnabled()).Struct(c)
}

func (c *Config) Load(ctx context.Context) error {
	_, span := tracer.Start(ctx, "kafka.Config.Load")
	defer span.End()

	if err := env.Parse(c); err != nil {
		return errors.Wrap(err, "config loading")
	}

	if err := c.validate(); err != nil {
		return errors.Wrap(err, "config validation")
	}

	return nil
}
