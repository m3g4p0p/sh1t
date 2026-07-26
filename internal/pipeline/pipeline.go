package pipeline

import (
	"context"
	"sync"
)

type Step func(ctx context.Context) (Step, error)

func Parallel(steps ...Step) Step {
	return func(ctx context.Context) (Step, error) {
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

func ParallelN(n int, steps ...Step) Step {
	return func(ctx context.Context) (Step, error) {
		ctx, cancel := context.WithCancelCause(ctx)
		defer cancel(nil)

		var wg sync.WaitGroup
		ch := make(chan Step)

		go func() {
			defer close(ch)

			for _, step := range steps {
				ch <- step
			}
		}()

		for range n {
			wg.Go(func() {
				if err := stepWorker(ctx, ch); err != nil {
					cancel(err)
				}
			})
		}

		wg.Wait()
		return nil, context.Cause(ctx)
	}
}

func Sequence(steps ...Step) Step {
	return func(ctx context.Context) (Step, error) {
		for _, step := range steps {
			if err := Run(ctx, step); err != nil {
				return nil, err
			}
		}
		return nil, nil
	}
}

func Run(ctx context.Context, step Step) error {
	for step != nil {
		var err error
		step, err = step(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func stepWorker(ctx context.Context, ch <-chan Step) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case step, ok := <-ch:
			if !ok {
				return nil
			}
			if err := Run(ctx, step); err != nil {
				return err
			}
		}
	}
}
