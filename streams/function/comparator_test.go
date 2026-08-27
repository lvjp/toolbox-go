package function

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func newComparator(ret int) Comparator[any] {
	return func(_, _ any) int {
		return ret
	}
}

func TestComparatorReversed(t *testing.T) {
	testCases := []struct {
		comparator Comparator[any]
		expected   int
	}{
		{newComparator(-1), 1},
		{newComparator(0), 0},
		{newComparator(1), -1},
	}

	for _, tc := range testCases {
		name := fmt.Sprintf("%d", tc.comparator(nil, nil))
		t.Run(name, func(t *testing.T) {
			actual := tc.comparator.Reversed()(nil, nil)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestComparatorAndThen(t *testing.T) {
	const otherValue = 42

	testCases := []struct {
		name       string
		comparator Comparator[any]
		expected   int
	}{
		{"negative", newComparator(-1), -1},
		{"zero", newComparator(0), otherValue},
		{"positive", newComparator(1), -1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.NotPanics(t, func() {
				tc.comparator.AndThen(newComparator(otherValue))(nil, nil)
			})
		})
	}
}
