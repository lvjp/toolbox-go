package streams

import (
	"iter"
	"slices"

	"github.com/lvjp/toolbox-go/streams/function"
)

func Of[T any](v iter.Seq[T]) *Stream[T] {
	return &Stream[T]{
		seq: v,
	}
}

func OfSlice[S ~[]E, E any](v S) *Stream[E] {
	return Of(slices.Values(v))
}

type Stream[T any] struct {
	seq iter.Seq[T]
}

func (s *Stream[T]) Iter() iter.Seq[T] {
	return s.seq
}

func (s *Stream[T]) Collect() []T {
	return slices.Collect(s.seq)
}

func (s *Stream[T]) Concat(v *Stream[T]) *Stream[T] {
	return &Stream[T]{
		seq: func(yield func(T) bool) {
			for _, it := range []iter.Seq[T]{s.seq, v.Iter()} {
				for v := range it {
					if !yield(v) {
						return
					}
				}
			}
		},
	}
}

func (s *Stream[T]) Count() uint64 {
	var count uint64

	for range s.seq {
		count++
		if count == 0 {
			panic("count overflow")
		}
	}

	return count
}

func (s *Stream[T]) ForEach(consumer function.Consumer[T]) {
	for v := range s.seq {
		consumer(v)
	}
}

func (s *Stream[T]) AllMatch(predicate function.Predicate[T]) bool {
	for v := range s.seq {
		if !predicate(v) {
			return false
		}
	}

	return true
}

func (s *Stream[T]) AnyMatch(predicate function.Predicate[T]) bool {
	for v := range s.seq {
		if predicate(v) {
			return true
		}
	}
	return false
}

func (s *Stream[T]) NoneMatch(predicate function.Predicate[T]) bool {
	for v := range s.seq {
		if predicate(v) {
			return false
		}
	}
	return true
}

func (s *Stream[T]) FindAny() (T, bool) {
	return s.FindFirst()
}

func (s *Stream[T]) FindFirst() (T, bool) {
	for v := range s.seq {
		return v, true
	}

	var zeroValue T
	return zeroValue, false
}

func (s *Stream[T]) Max(comparator function.Comparator[T]) (T, bool) {
	var maxValue T
	var found bool
	for v := range s.seq {
		if !found || comparator(v, maxValue) > 0 {
			maxValue = v
			found = true
		}
	}
	return maxValue, found
}

func (s *Stream[T]) Min(comparator function.Comparator[T]) (T, bool) {
	var minValue T
	var found bool
	for v := range s.seq {
		if !found || comparator(v, minValue) < 0 {
			minValue = v
			found = true
		}
	}
	return minValue, found
}

func (s *Stream[T]) Reduce(accumulator function.BinaryOperator[T]) (T, bool) {
	var result T
	var found bool
	for v := range s.seq {
		if !found {
			result = v
			found = true
		} else {
			result = accumulator(result, v)
		}
	}
	return result, found
}

func (s *Stream[T]) ReduceWithIdentity(
	identity T,
	accumulator function.BinaryOperator[T],
) T {
	result := identity
	for v := range s.seq {
		result = accumulator(result, v)
	}
	return result
}

func (s *Stream[T]) Map[U any](mapper function.Function[T, U]) *Stream[U] {
	return &Stream[U]{
		seq: func(yield func(U) bool) {
			for t := range s.seq {
				if !yield(mapper(t)) {
					break
				}
			}
		},
	}
}

func (s *Stream[T]) FlatMap[U any](mapper function.Function[T, *Stream[U]]) *Stream[U] {
	return &Stream[U]{
		func(yield func(U) bool) {
			for t := range s.seq {
				for u := range mapper(t).Iter() {
					if !yield(u) {
						break
					}
				}
			}
		},
	}
}

func (s *Stream[T]) Filter(predicate function.Predicate[T]) *Stream[T] {
	return &Stream[T]{
		seq: func(yield func(T) bool) {
			for t := range s.seq {
				if predicate(t) && !yield(t) {
					break
				}
			}
		},
	}
}

func (s *Stream[T]) Limit(maxSize uint64) *Stream[T] {
	return &Stream[T]{
		seq: func(yield func(T) bool) {
			var count uint64
			for t := range s.seq {
				if count >= maxSize {
					break
				}
				if !yield(t) {
					break
				}
				count++
			}
		},
	}
}

func (s *Stream[T]) TakeWhile(predicate function.Predicate[T]) *Stream[T] {
	return &Stream[T]{
		seq: func(yield func(T) bool) {
			for t := range s.seq {
				if !predicate(t) || !yield(t) {
					break
				}
			}
		},
	}
}

func (s *Stream[T]) DropWhile(predicate function.Predicate[T]) *Stream[T] {
	return &Stream[T]{
		seq: func(yield func(T) bool) {
			drop := true
			for t := range s.seq {
				if drop {
					if predicate(t) {
						continue
					}
					drop = false
				}

				if !yield(t) {
					break
				}
			}
		},
	}
}

func (s *Stream[T]) Peek(consumer function.Consumer[T]) *Stream[T] {
	return &Stream[T]{
		seq: func(yield func(T) bool) {
			for t := range s.seq {
				consumer(t)
				if !yield(t) {
					break
				}
			}
		},
	}
}

func (s *Stream[T]) Skip(n uint64) *Stream[T] {
	return &Stream[T]{
		seq: func(yield func(T) bool) {
			var count uint64
			for t := range s.seq {
				if count < n {
					count++
					continue
				}
				if !yield(t) {
					break
				}
			}
		},
	}
}
