package function

type Predicate[T any] func(T) bool

func (p Predicate[T]) Negate() Predicate[T] {
	return func(v T) bool {
		return !p(v)
	}
}

func (p Predicate[T]) And(other Predicate[T]) Predicate[T] {
	return func(v T) bool {
		return p(v) && other(v)
	}
}

func (p Predicate[T]) Or(other Predicate[T]) Predicate[T] {
	return func(v T) bool {
		return p(v) || other(v)
	}
}

func PredicateIsSame[T comparable](reference T) Predicate[T] {
	return func(v T) bool {
		return v == reference
	}
}
