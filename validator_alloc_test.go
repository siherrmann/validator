package validator

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// validateAllocBudget sits a little above the ~8 allocs/op measured today, so only a real regression fails.
const validateAllocBudget = 12

func TestValidateAllocationBudgetWhenCached(t *testing.T) {
	s := &BenchStructFourFields{Field1: "hello", Field2: 50, Field3: 42.5, Field4: "ignored"}

	// Warm the cache and check success first, an error path would allocate for the error message.
	require.NoError(t, Validate(s))

	allocs := testing.AllocsPerRun(1000, func() {
		_ = Validate(s)
	})
	assert.LessOrEqual(t, allocs, float64(validateAllocBudget))
}

func TestValidateConcurrentSameType(t *testing.T) {
	const goroutines = 50
	const iterations = 200

	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				valid := (n+j)%2 == 0
				s := &BenchStructFourFields{Field1: "hello", Field2: 50, Field3: 42.5, Field4: "ignored"}
				if !valid {
					s.Field1 = "x"
				}

				err := Validate(s)
				if valid {
					assert.NoError(t, err)
				} else {
					assert.Error(t, err)
				}
			}
		}(i)
	}
	wg.Wait()
}
