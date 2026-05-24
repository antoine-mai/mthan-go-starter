package services

import (
	"errors"
	"fmt"
	"strings"
)

// ProcessService manages all business logic and request processing.
type ProcessService struct{}

// NewProcessService initializes and returns a new ProcessService instance.
func NewProcessService() *ProcessService {
	return &ProcessService{}
}

// ProcessAPI implements business processing for public /api endpoints.
func (s *ProcessService) ProcessAPI(input string) (string, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", errors.New("input data cannot be empty")
	}
	// Business logic: Transform and return output
	processedVal := strings.ToUpper(trimmed)
	return fmt.Sprintf("Processed via API: %s", processedVal), nil
}

// ProcessPost implements business processing for same-domain /post endpoints.
func (s *ProcessService) ProcessPost(payload map[string]string) (string, error) {
	action, exists := payload["action"]
	if !exists || strings.TrimSpace(action) == "" {
		return "", errors.New("missing action field in payload")
	}
	
	// Business logic: Simulate internal processing
	return fmt.Sprintf("Executed internal action: %s successfully", action), nil
}
