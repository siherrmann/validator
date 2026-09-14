package validator

import (
	"fmt"
	"sync"
	"testing"
)

// validateAllocBudget is the ceiling asserted below for a Validate call once its
// type's extraction is already cached. It sits comfortably above the ~23 allocs/op
// this repo measures today (see BenchmarkValidateCachedType) so the test only fails on
// a real regression, not on incidental noise from an unrelated runtime change.
const validateAllocBudget = 30

// TestValidateAllocationBudgetWhenCached is the allocation-budget test called for by
// this package's performance work: once a type's validations are cached,
// Validate must not regress back towards the pre-cache allocation count (~89 for this
// struct shape).
func TestValidateAllocationBudgetWhenCached(t *testing.T) {
	s := &BenchStructFourFields{Field1: "hello", Field2: 50, Field3: 42.5, Field4: "ignored"}

	// Warm the cache and assert success before measuring. An error path formats and
	// allocates error strings freely, and measuring that by accident would pass this
	// test for the wrong reason.
	if err := Validate(s); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	allocs := testing.AllocsPerRun(1000, func() {
		if err := Validate(s); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	if allocs > validateAllocBudget {
		t.Fatalf("expected at most %d allocations for a cached-type Validate call, got %v", validateAllocBudget, allocs)
	}
}

// TestValidateConcurrentSameType exercises the shared validationCache from many
// goroutines validating the same struct type at once. Run with -race: a cache that
// handed out mutable state, or raced on population, would show up here.
func TestValidateConcurrentSameType(t *testing.T) {
	const goroutines = 50
	const iterations = 200

	var wg sync.WaitGroup
	errs := make(chan error, goroutines)

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				// Alternate valid/invalid values across goroutines and iterations so
				// both the success and error paths run concurrently over the same
				// cached []model.Validation.
				valid := (n+j)%2 == 0
				s := &BenchStructFourFields{Field1: "hello", Field2: 50, Field3: 42.5, Field4: "ignored"}
				if !valid {
					s.Field1 = "x" // shorter than min3, should fail validation
				}

				err := Validate(s)
				if valid && err != nil {
					errs <- fmt.Errorf("goroutine %d iteration %d: expected no error, got %v", n, j, err)
					return
				}
				if !valid && err == nil {
					errs <- fmt.Errorf("goroutine %d iteration %d: expected an error, got none", n, j)
					return
				}
			}
		}(i)
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		t.Errorf("unexpected validation result: %v", err)
	}
}
