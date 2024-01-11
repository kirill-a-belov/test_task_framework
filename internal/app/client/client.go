// Package client implements general TCP client
package client

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/pkg/errors"

	"github.com/kirill-a-belov/test_task_framework/internal/app/client/pkg/config"
	"github.com/kirill-a-belov/test_task_framework/internal/pkg/network"
	"github.com/kirill-a-belov/test_task_framework/internal/pkg/protocol"
	"github.com/kirill-a-belov/test_task_framework/pkg/context_helper"
	"github.com/kirill-a-belov/test_task_framework/pkg/logger"
	"github.com/kirill-a-belov/test_task_framework/pkg/rand"
	"github.com/kirill-a-belov/test_task_framework/pkg/tracer"
)

func New(ctx context.Context, config *config.Config) *Client {
	_, span := tracer.Start(ctx, "internal.app.client.New")
	defer span.End()

	return &Client{
		config:   config,
		stopChan: make(chan struct{}),
		logger:   logger.New("client"),

		dialler: func() (net.Conn, error) {
			return net.Dial(protocol.NetworkType, config.Address)
		},
	}
}

type Client struct {
	config   *config.Config
	stopChan chan struct{}
	logger   logger.Logger

	dialler func() (net.Conn, error)
}

func (c *Client) Start(ctx context.Context) error {
	_, span := tracer.Start(ctx, "internal.app.client.Client.Start")
	defer span.End()

	c.logger.Info(fmt.Sprintf("config: %+v", *c.config))

	go c.processor(ctx, handle)

	return nil
}

func (c *Client) processor(ctx context.Context, handler func(io.ReadWriter) error) {
	_, span := tracer.Start(ctx, "internal.app.client.Client.processor")
	defer span.End()

	for {
		select {
		case <-c.stopChan:
			c.logger.Info("processor terminated")

			return
		default:
			conn, err := c.dialler()
			if err != nil {
				c.logger.Error(err, "connection dialing")

				continue
			}

			go func() {
				defer func(conn net.Conn) {
					_ = conn.Close()
				}(conn)
				if err := context_helper.RunWithTimeout(c.config.ConnTTL, func() error {
					return handler(conn)
				}); err != nil {
					c.logger.Error(err, "connection handling")
				}
			}()

			time.Sleep(c.config.Delay)
		}
	}
}

func (c *Client) Stop(ctx context.Context) {
	_, span := tracer.Start(ctx, "internal.app.client.Client.Stop")
	defer span.End()

	close(c.stopChan)
}

func handle(conn io.ReadWriter) error {
	ctx, span := tracer.Start(context.Background(), "internal.app.client.Client.handle")
	defer span.End()

	const (
		payloadSize     = 3
		payloadMaxDigit = 1024
	)
	payload, err := rand.Rand(payloadSize, payloadMaxDigit)
	if err != nil {
		return errors.Wrap(err, "generating payload")
	}

	if err := network.Send(ctx, conn, protocol.Request{
		Message: protocol.Message{
			Type: protocol.MessageTypeRequest,
		},
		Payload: payload,
	}); err != nil {
		return errors.Wrap(err, "sending request")
	}

	response, err := network.Receive[protocol.Response](ctx, conn)
	if err != nil {
		return errors.Wrap(err, "receiving server response")
	}
	if response.Type != protocol.MessageTypeResponse {
		return errors.Errorf("server response: received wrong message (%v)", response)
	}

	logger.New("client.handle").Info("handling connection",
		"sent", payload,
		"received", response.Payload,
	)

	return nil
}
