package network

import (
	"context"
	"encoding/gob"
	"io"

	"github.com/pkg/errors"

	"github.com/kirill-a-belov/test_task_framework/pkg/tracer"
)

func Send[A any](ctx context.Context, c io.ReadWriter, msg A) error {
	_, span := tracer.Start(ctx, "internal.pkg.network.Send")
	defer span.End()

	enc := gob.NewEncoder(c)

	if err := enc.Encode(msg); err != nil {
		return errors.Wrapf(err, "encoding msg (%v)", msg)
	}

	return nil
}

func Receive[A any](ctx context.Context, c io.ReadWriter) (A, error) {
	_, span := tracer.Start(ctx, "internal.pkg.network.Receive")
	defer span.End()

	var result A

	dec := gob.NewDecoder(c)

	if err := dec.Decode(&result); err != nil {
		return result, errors.Wrapf(err, "decoding msg (%v)", result)
	}

	return result, nil
}
