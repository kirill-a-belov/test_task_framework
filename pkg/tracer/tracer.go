package tracer

import "context"

type exampleSpan struct{}

func (span *exampleSpan) End() {}

type Span interface {
	End()
}

func Start(ctx context.Context, name string) (context.Context, Span) {
	_ = name

	return ctx, &exampleSpan{}
}
