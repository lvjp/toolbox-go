package function

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConsumer_AndThen(t *testing.T) {
	expectedOrder := []string{"first", "second"}
	var callOrder []string

	value := "consumed value"

	var first Consumer[*string] = func(v *string) {
		callOrder = append(callOrder, expectedOrder[0])
		require.Same(t, &value, v)
	}

	var second Consumer[*string] = func(v *string) {
		callOrder = append(callOrder, expectedOrder[1])
		require.Same(t, &value, v)
	}

	first.AndThen(second)(&value)

	require.Equal(t, expectedOrder, callOrder)
}
