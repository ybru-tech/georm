package testutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func CheckErr(t testing.TB, expect, actual error) (finished bool) {
	t.Helper()

	if expect == nil && actual == nil {
		return false
	}

	finished = true

	if expect == nil && actual != nil {
		t.Fatalf("unexpected error: %v", actual)
		return
	}

	if expect != nil && actual == nil {
		t.Fatalf("expected error: %v", expect)
		return
	}

	require.ErrorIs(t, actual, expect)

	return
}
