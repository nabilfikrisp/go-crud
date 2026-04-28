package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nabilfikrisp/go-crud/internal/usecase"
	"github.com/nabilfikrisp/go-crud/pkg/logger"
)

// NewRoutes -.
func NewRoutes(apiV1Group gin.RouterGroup, c usecase.Contact, l logger.Interface) {
	r := &V1{
		c: c,
		l: l,
		v: validator.New(validator.WithRequiredStructEnabled()),
	}

	// Contact
	contactGroup := apiV1Group.Group("/contacts")
	{
		contactGroup.POST("", r.createContact)
		contactGroup.GET("", r.listContacts)
		contactGroup.GET("/:id", r.getContactByID)
		contactGroup.PUT("/:id", r.updateContact)
		contactGroup.DELETE("/:id", r.deleteContact)
	}
}
