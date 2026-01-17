package engine

import (
	"errors"
)

// validateExecuteParams validates the execute parameters
func validateExecuteParams(params ExecuteSpec) error {
	if params.Code == "" {
		return errors.New("code is required")
	}
	if params.Language == "" {
		return errors.New("language is required")
	}
	return nil
}
