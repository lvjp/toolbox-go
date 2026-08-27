package function

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func predicateTrue[T any](T) bool  { return true }
func predicateFalse[T any](T) bool { return false }

func TestPredicate_Negate(t *testing.T) {
	testCases := []struct {
		name      string
		predicate Predicate[any]
		expected  bool
	}{
		{"0", predicateFalse[any], true},
		{"1", predicateTrue[any], false},
	}

	for _, tc := range testCases {
		name := fmt.Sprintf("%t", tc.predicate(nil))
		t.Run(name, func(t *testing.T) {
			actual := tc.predicate.Negate()(nil)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestPredicate_And(t *testing.T) {
	testCases := []struct {
		a, b     Predicate[any]
		expected bool
	}{
		{predicateFalse[any], predicateFalse[any], false},
		{predicateFalse[any], predicateTrue[any], false},
		{predicateTrue[any], predicateFalse[any], false},
		{predicateTrue[any], predicateTrue[any], true},
	}

	for _, tc := range testCases {
		name := fmt.Sprintf("%t&&%t", tc.a(nil), tc.b(nil))
		t.Run(name, func(t *testing.T) {
			actual := tc.a.And(tc.b)(nil)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestPredicate_Or(t *testing.T) {
	testCases := []struct {
		a, b     Predicate[any]
		expected bool
	}{
		{predicateFalse[any], predicateFalse[any], false},
		{predicateFalse[any], predicateTrue[any], true},
		{predicateTrue[any], predicateFalse[any], true},
		{predicateTrue[any], predicateTrue[any], true},
	}

	for _, tc := range testCases {
		name := fmt.Sprintf("%t&&%t", tc.a(nil), tc.b(nil))
		t.Run(name, func(t *testing.T) {
			actual := tc.a.Or(tc.b)(nil)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestPredicateIsSame(t *testing.T) {
	testCases := []struct {
		reference any
		input     any
		expected  bool
	}{
		{"foo", "foo", true},
		{"foo", "bar", false},
		{"bar", "foo", false},
		{"bar", "bar", true},
	}

	for _, tc := range testCases {
		name := fmt.Sprintf("%v==%v", tc.reference, tc.input)
		t.Run(name, func(t *testing.T) {
			actual := PredicateIsSame(tc.reference)(tc.input)
			require.Equal(t, tc.expected, actual)
		})
	}
}
