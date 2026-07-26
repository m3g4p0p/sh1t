package pipeline

import (
	"context"
	"sync"
)

type Step[T any] func(ctx context.Context) (Step[T], error)

type AnyStep = Step[any]

func Parallel[T any](steps ...Step[T]) Step[T] {
	return func(ctx context.Context) (Step[T], error) {
		ctx, cancel := context.WithCancelCause(ctx)
		defer cancel(nil)

		var wg sync.WaitGroup

		for _, step := range steps {
			wg.Go(func() {
				if err := Run(ctx, step); err != nil {
					cancel(err)
				}
			})
		}

		wg.Wait()
		return nil, context.Cause(ctx)
	}
}

func Sequence[T any](steps ...Step[T]) Step[T] {
	return func(ctx context.Context) (Step[T], error) {
		for _, step := range steps {
			if err := Run(ctx, step); err != nil {
				return nil, err
			}
		}
		return nil, nil
	}
}

func Run[T any](ctx context.Context, step Step[T]) error {
	for step != nil {
		var err error
		step, err = step(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
