package validator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInject(t *testing.T) {
	t.Run("Several values", func(t *testing.T) {
		ctx := Inject(context.Background(), "user_id", "u1")
		ctx = Inject(ctx, "limit", 10)

		userID, ok := Get[string](ctx, "user_id")
		require.True(t, ok)
		assert.Equal(t, "u1", userID)
		limit, ok := Get[int](ctx, "limit")
		require.True(t, ok)
		assert.Equal(t, 10, limit)
	})

	t.Run("Parent context keeps its value", func(t *testing.T) {
		parent := Inject(context.Background(), "role", "user")
		child := Inject(parent, "role", "admin")

		role, _ := Get[string](parent, "role")
		assert.Equal(t, "user", role)
		role, _ = Get[string](child, "role")
		assert.Equal(t, "admin", role)
	})

	t.Run("Nil context", func(t *testing.T) {
		//nolint:staticcheck // A nil context is handled on purpose.
		ctx := Inject(nil, "a", 1)
		value, ok := Get[int](ctx, "a")
		require.True(t, ok)
		assert.Equal(t, 1, value)
	})
}

func TestGet(t *testing.T) {
	ctx := Inject(context.Background(), "limit", 10)

	value, ok := Get[int](ctx, "limit")
	assert.True(t, ok)
	assert.Equal(t, 10, value)

	_, ok = Get[string](ctx, "limit")
	assert.False(t, ok, "Expected no value for the wrong type")

	_, ok = Get[int](ctx, "missing")
	assert.False(t, ok, "Expected no value for a missing key")

	_, ok = Get[int](context.Background(), "limit")
	assert.False(t, ok, "Expected no value without injected values")

	//nolint:staticcheck // A nil context is handled on purpose.
	_, ok = Get[int](nil, "limit")
	assert.False(t, ok)
}
