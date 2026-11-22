package validator

import (
	"context"
	"maps"
)

// This file holds the values injected for a single validation call,
// used by `${name}` in requirements and by functions added with AddValidationFuncContext.

// injectedKey is the context key of the injected values.
type injectedKey struct{}

// Inject returns a copy of ctx with the value injected under the given key.
// The value can be used in requirements with `${key}` (eg. `equ${api_key}` or `rel(${role} equadmin) => min1`)
// and read with [Get] in functions added with AddValidationFuncContext.
//
// The values belong to the context, so they only exist for the validation calls with this context
// (eg. a single request). The given context is not changed, so Inject is safe for concurrent use.
//
// The unpack functions (eg. UnmarshalAndValidate) use the context of the request,
// all other functions have a Context variant (eg. ValidateContext).
func Inject(ctx context.Context, key string, value any) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	// All values are kept in one map, so a lookup is one ctx.Value call, regardless of the number of values.
	// The map is copied, because the map of the given context may be used by other contexts.
	current, _ := ctx.Value(injectedKey{}).(map[string]any)
	values := make(map[string]any, len(current)+1)
	maps.Copy(values, current)
	values[key] = value
	return context.WithValue(ctx, injectedKey{}, values)
}

// Get returns the value injected under the given key with [Inject].
// ok is false if no value was injected under the key or if it is not of type T.
func Get[T any](ctx context.Context, key string) (value T, ok bool) {
	injected, found := injectedValue(ctx, key)
	if !found {
		return value, false
	}
	value, ok = injected.(T)
	return value, ok
}

// injectedValue returns the value injected under the given key.
func injectedValue(ctx context.Context, key string) (any, bool) {
	if ctx == nil {
		return nil, false
	}
	values, _ := ctx.Value(injectedKey{}).(map[string]any)
	value, ok := values[key]
	return value, ok
}
