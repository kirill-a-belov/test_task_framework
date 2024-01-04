package config

import (
	"context"
	"fmt"
	"github.com/kirill-a-belov/temp_test_task/utils/tracing"
	"github.com/pkg/errors"
	"os"
	"strconv"
	"time"
)

type Config struct {
	port         int
	connPoolSize int
	connTTL      time.Duration
}

func (c *Config) validate() error {
	if c.port < 0 || c.port > 65535 {
		return fmt.Errorf("invalid port number: %d", c.port)
	}

	if c.connPoolSize < 0 || c.connPoolSize > 1024 {
		return fmt.Errorf("invalid connection pool size: %d", c.connPoolSize)
	}

	if c.connTTL < time.Millisecond || c.connTTL > time.Second {
		return fmt.Errorf("invalid connextion TTL: %v", c.connPoolSize)
	}

	return nil
}

const (
	portEnvVarName         = "SERVER_PORT"
	connPoolSizeEnvVarName = "SERVER_CONN_POOL_SIZE"
	connTTLEnvVarName      = "SERVER_CONN_TTL_SEC"
)

func (c *Config) Load(ctx context.Context) error {
	span, _ := tracing.NewSpan(ctx, "server.Config.Load")
	defer span.Close()

	port := os.Getenv(portEnvVarName)
	if port == "" {
		return errors.Errorf("empty env var (%s) value", portEnvVarName)
	}
	var err error
	if c.port, err = strconv.Atoi(port); err != nil {
		return errors.Wrapf(err, "parsing value (%s) from env var (%s)", port, portEnvVarName)
	}

	connPoolSize := os.Getenv(connPoolSizeEnvVarName)
	if connPoolSize == "" {
		return errors.Errorf("empty env var (%s) value", connPoolSize)
	}
	if c.connPoolSize, err = strconv.Atoi(connPoolSize); err != nil {
		return errors.Wrapf(err, "parsing value (%s) from env var (%s)", connPoolSize, connPoolSizeEnvVarName)
	}

	connTTLSec := os.Getenv(connTTLEnvVarName)
	if connTTLSec == "" {
		return errors.Errorf("empty env var (%s) value", connTTLSec)
	}

	connTTLSecParsed, err := strconv.Atoi(connTTLSec)
	if err != nil {
		return errors.Wrapf(err, "parsing value (%s) from env var (%s)", connTTLSec, connTTLEnvVarName)
	}
	c.connTTL = time.Second * time.Duration(connTTLSecParsed)

	return c.validate()
}
