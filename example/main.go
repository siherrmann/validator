package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/siherrmann/validator"
)

type Error struct {
	ID                  int       `json:"id" del:"min1"`
	StatusCode          int       `json:"status_code" vld:"min100" upd:"min100, gr1min1"`
	Message             string    `json:"message" vld:"min1" upd:"min1, gr1min1"`
	UnderlyingException string    `json:"underlying_exception" vld:"min1, gr1min1" upd:"min1, gr1min1"`
	CreatedAt           time.Time `json:"created_at" vld:"-"`
}

func HandleError(w http.ResponseWriter, r *http.Request) {
	// This uses the default `vld` tag. It unmarshals the request body, validates it
	// and updates the `newError` variable. You could use that for creating a new instance
	// of error with all needed parameters.
	var newError Error
	err := validator.UnmarshalValidateAndUpdate(r, &newError)
	if err != nil {
		http.Error(w, fmt.Sprintf("error validating new errror: %v", err), http.StatusBadRequest)
		return
	}

	// For updating an error you could use the `upd` tag (including goups) to make sure
	// that at least one of the values is updated and if so is valid.
	var exisitingErrorFromDb Error
	err = validator.UnmarshalValidateAndUpdate(r, &exisitingErrorFromDb, "upd")
	if err != nil {
		http.Error(w, fmt.Sprintf("error validating error update: %v", err), http.StatusBadRequest)
		return
	}

	// If you would only want to validate some given struct with another tag
	// (for example if you want to check if a given error is valid to delete,
	// containing an id in this case) you could do:
	err = validator.Validate(&newError, "del")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error unmarshaling JSON: %v", err), http.StatusBadRequest)
		return
	}
}
