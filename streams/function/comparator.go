package function

type Comparator[T any] func(a, b T) int

func (c Comparator[T]) Reversed() Comparator[T] {
	return func(a, b T) int {
		return -c(a, b)
	}
}

func (c Comparator[T]) AndThen(other Comparator[T]) Comparator[T] {
	return func(a, b T) int {
		result := c(a, b)
		if result != 0 {
			return result
		}
		return other(a, b)
	}
}
