package pipeline

type Step[T any] func(state T) (T, error)

func Run[T any](state T, steps ...Step[T]) error {
	for _, s := range steps {
		var err error
		state, err = s(state)
		if err != nil {
			return err
		}
	}
	return nil
}
