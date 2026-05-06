package v1

import (
	"github.com/go-playground/validator/v10"
	"github.com/nabilfikrisp/go-crud/internal/usecase"
	"github.com/nabilfikrisp/go-crud/pkg/logger"
)

// V1 contains dependencies for version 1 REST API handlers.
type V1 struct {
	c usecase.Contact
	l logger.Interface
	v *validator.Validate
}
