package network

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"

	"playGround/utils/tracing"
)

type conn interface {
	Write(b []byte) (int, error)
	Read(b []byte) (int, error)
}

func Send[A any](ctx context.Context, c conn, msg A) error {
	span, _ := tracing.NewSpan(ctx, "utils.network.Conn.Send")
	defer span.Close()

	body, err := json.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "marshalling initial request")
	}
	if _, err := c.Write(body); err != nil {
		return errors.Wrap(err, "sending initial request")
	}

	return nil
}

func Receive[A any](ctx context.Context, c conn) (A, error) {
	span, _ := tracing.NewSpan(ctx, "utils.network.Conn.Receive")
	defer span.Close()

	var result A

	buffer := make([]byte, 1024)
	n, err := c.Read(buffer)
	if err != nil {
		return result, errors.Wrap(err, "reading server question")
	}
	buffer = buffer[:n]
	if err := json.Unmarshal(buffer, &result); err != nil {
		return result, errors.Wrap(err, "unmarshalling server question request")
	}

	return result, nil
}
