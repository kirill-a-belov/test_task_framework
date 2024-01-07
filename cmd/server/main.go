package main

import (
	"context"

	"github.com/kirill-a-belov/test_task_framework/internal/app/cmd/server"
	"github.com/kirill-a-belov/test_task_framework/pkg/logger"
)

func main() {
	log := logger.New("cmd.server")
	cmd := server.New(context.Background())

	if err := cmd.Execute(); err != nil {
		log.Error(err, "running server")
	}
}
