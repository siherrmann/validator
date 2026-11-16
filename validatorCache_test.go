package validator

import (
	"reflect"
	"testing"

	"github.com/siherrmann/validator/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCachedValidationsFromStruct(t *testing.T) {
	type Cached struct {
		Name string `json:"name" vld:"min1" upd:"max5"`
	}

	t.Run("Same type and tag type is cached", func(t *testing.T) {
		first, err := getCachedValidationsFromStruct(&Cached{}, model.VLD)
		require.NoError(t, err)
		second, err := getCachedValidationsFromStruct(&Cached{Name: "other"}, model.VLD)
		require.NoError(t, err)
		require.Len(t, first, 1)
		assert.Same(t, &first[0], &second[0], "Expected the cached slice to be returned")
	})

	t.Run("Tag types are cached separately", func(t *testing.T) {
		vld, err := getCachedValidationsFromStruct(&Cached{}, model.VLD)
		require.NoError(t, err)
		upd, err := getCachedValidationsFromStruct(&Cached{}, "upd")
		require.NoError(t, err)
		assert.Equal(t, "min1", vld[0].Requirement)
		assert.Equal(t, "max5", upd[0].Requirement)
	})

	t.Run("Errors are not cached", func(t *testing.T) {
		type Invalid struct {
			Name string `vld:"invalid"`
		}
		_, err := getCachedValidationsFromStruct(&Invalid{}, model.VLD)
		require.Error(t, err)
		_, cached := validationCache.Load(validationCacheKey{structType: reflect.TypeOf(&Invalid{}), tagType: model.VLD})
		assert.False(t, cached)
	})
}
