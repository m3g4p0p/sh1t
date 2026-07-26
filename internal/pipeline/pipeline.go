package pipeline

type Step[T any] func(state T) error

func Pipe[T any](steps ...Step[T]) Step[T] {
	return func(state T) error {
		return Run(state, steps...)
	}
}

func Run[T any](state T, steps ...Step[T]) error {
	for _, s := range steps {
		if err := s(state); err != nil {
			return err
		}
	}
	return nil
}
