package usecase

import (
	"context"

	"github.com/nabilfikrisp/go-crud/internal/entity"
)

type (
	// Contact defines the interface for contact use cases.
	Contact interface {
		Create(ctx context.Context, req CreateContact) (entity.Contact, error)
		GetByID(ctx context.Context, id string) (entity.Contact, error)
		List(ctx context.Context, filter entity.ContactFilter) ([]entity.Contact, int, error)
		Update(ctx context.Context, id string, req UpdateContact) (entity.Contact, error)
		Delete(ctx context.Context, id string) error
	}

	// CreateContactRequest defines the structure for creating a new contact.
	CreateContact struct {
		FirstName    string                      `json:"first_name"`
		LastName     string                      `json:"last_name"`
		Email        string                      `json:"email"`
		PhoneNumber  string                      `json:"phone_number"`
		Relationship *entity.ContactRelationship `json:"relationship"`
	}

	// UpdateContactRequest defines the structure for updating an existing contact.
	UpdateContact struct {
		FirstName    *string                     `json:"first_name"`
		LastName     *string                     `json:"last_name"`
		Email        *string                     `json:"email"`
		PhoneNumber  *string                     `json:"phone_number"`
		Relationship *entity.ContactRelationship `json:"relationship"`
	}
)
