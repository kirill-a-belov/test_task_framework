package network

import (
	"context"
	"encoding/gob"
	"github.com/kirill-a-belov/test_task_framework/pkg/tracer"

	"github.com/pkg/errors"
)

type conn interface {
	Write(b []byte) (int, error)
	Read(b []byte) (int, error)
}

func Send[A any](ctx context.Context, c conn, msg A) error {
	_, span := tracer.Start(ctx, "internal.pkg.network.Send")
	defer span.End()

	enc := gob.NewEncoder(c)

	if err := enc.Encode(msg); err != nil {
		return errors.Wrapf(err, "encoding msg (%v)", msg)
	}

	return nil
}

func Receive[A any](ctx context.Context, c conn) (A, error) {
	_, span := tracer.Start(ctx, "internal.pkg.network.Receive")
	defer span.End()

	var result A

	dec := gob.NewDecoder(c)

	if err := dec.Decode(&result); err != nil {
		return result, errors.Wrapf(err, "decoding msg (%v)", result)
	}

	return result, nil
}
