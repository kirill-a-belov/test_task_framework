package main

import (
	"context"

	"github.com/kirill-a-belov/test_task_framework/internal/app/cmd/client"
	"github.com/kirill-a-belov/test_task_framework/pkg/logger"
)

func main() {
	log := logger.New("cmd.client")
	cmd := client.New(context.Background())

	if err := cmd.Execute(); err != nil {
		log.Error(err, "running client")
	}
}
