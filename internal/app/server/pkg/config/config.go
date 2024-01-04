package config

import (
	"context"
	"fmt"
	"github.com/caarlos0/env"
	"github.com/kirill-a-belov/test_task_framework/pkg/tracer"
	"github.com/pkg/errors"
	"time"
)

type Config struct {
	Port         int           `env:"SERVER_PORT"`
	ConnPoolSize int           `env:"SERVER_CONN_POOL_SIZE"`
	ConnTTL      time.Duration `env:"SERVER_CONN_TTL"`
}

func (c *Config) validate() error {
	if c.Port < 0 || c.Port > 65535 {
		return fmt.Errorf("invalid Port number: %d", c.Port)
	}

	if c.ConnPoolSize < 1 || c.ConnPoolSize > 1024 {
		return fmt.Errorf("invalid connection pool size: %d", c.ConnPoolSize)
	}

	if c.ConnTTL < time.Millisecond || c.ConnTTL > time.Second {
		return fmt.Errorf("invalid connextion TTL: %v", c.ConnPoolSize)
	}

	return nil
}

func (c *Config) Load(ctx context.Context) error {
	_, span := tracer.Start(ctx, "server.Config.Load")
	defer span.End()

	if err := env.Parse(c); err != nil {
		return errors.Wrap(err, "config loading")
	}

	if err := c.validate(); err != nil {
		return errors.Wrap(err, "config validation")
	}

	return nil
}
