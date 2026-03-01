package httpTransport

import (
	"errors"

	"github.com/anurag-327/neuron-core/pkg/api"
)

// ValidateExecuteParams validates the execute parameters
func ValidateExecuteParams(params api.ExecuteRequest) error {
	if params.Code == "" {
		return errors.New("code is required")
	}
	if params.Language == "" {
		return errors.New("language is required")
	}
	return nil
}
