package pipeline

type Step[T any] func(state T) (Step[T], error)

func Sequence[T any](steps ...Step[T]) Step[T] {
	return func(state T) (Step[T], error) {
		for _, step := range steps {
			if err := Run(state, step); err != nil {
				return nil, err
			}
		}
		return nil, nil
	}
}

func Run[T any](state T, step Step[T]) error {
	for step != nil {
		var err error
		step, err = step(state)
		if err != nil {
			return err
		}
	}

	return nil
}
