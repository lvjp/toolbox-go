package function

type Consumer[T any] func(T)

func (c Consumer[T]) AndThen(after Consumer[T]) Consumer[T] {
	return func(t T) {
		c(t)
		after(t)
	}
}
