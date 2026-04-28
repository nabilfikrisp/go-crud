package v1

import (
	"github.com/go-playground/validator/v10"
	"github.com/nabilfikrisp/go-crud/internal/usecase"
	"github.com/nabilfikrisp/go-crud/pkg/logger"
)

type V1 struct {
	c usecase.Contact
	l logger.Interface
	v *validator.Validate
}
