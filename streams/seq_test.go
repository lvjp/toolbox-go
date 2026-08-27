package streams

import (
	"fmt"
	"slices"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStream_Concat(t *testing.T) {
	actual := OfSlice([]int{1, 2, 3}).
		Concat(OfSlice([]int{4, 5, 6})).
		Collect()

	expected := []int{1, 2, 3, 4, 5, 6}

	require.Equal(t, expected, actual)
}

func TestCount(t *testing.T) {
	testCases := []struct {
		expected uint64
		input    []any
	}{
		{0, []any{}},
		{1, []any{nil}},
		{2, []any{nil, nil}},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := OfSlice(tc.input).Count()
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestReduce(t *testing.T) {
	reducer := func(a, b int) int {
		return a * b
	}

	t.Run("normal", func(t *testing.T) {
		// By the fundamental theorem of arithmetic, every positive integer has a unique prime factorization.
		// Thats why we use only prime numbers.
		source := OfSlice([]int{1, 3, 7, 13})
		expected := 1 * 3 * 7 * 13

		value, ok := source.Reduce(reducer)
		require.True(t, ok)
		require.Equal(t, expected, value)
	})

	t.Run("empty", func(t *testing.T) {
		source := Of(slices.Values([]int{}))

		_, ok := source.Reduce(reducer)
		require.False(t, ok)
	})
}

func TestMap(t *testing.T) {
	source := Of(slices.Values([]int{0, 1, 3, 5}))
	expected := []string{"0b0", "0b1", "0b11", "0b101"}

	mapper := func(v int) string {
		return fmt.Sprintf("0b%b", v)
	}

	actual := source.Map(mapper).Collect()
	require.Equal(t, expected, actual)
}

func TestFlatMap(t *testing.T) {
	source := Of(slices.Values([]int{0, 1, 3, 5}))
	expected := []rune{
		// 0 is empty
		'0', 'b', '1', // 1
		'0', 'b', '1', '1', // 3
		'0', 'b', '1', '0', '1', // 5
	}

	mapper := func(v int) *Stream[rune] {
		if v == 0 {
			return Of(slices.Values([]rune{}))
		}

		return Of(slices.Values([]rune(fmt.Sprintf("0b%b", v))))
	}

	actual := source.FlatMap(mapper).Collect()
	require.Equal(t, expected, actual)
}
