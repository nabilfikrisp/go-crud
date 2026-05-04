package repo

import (
	"context"

	"github.com/nabilfikrisp/go-crud/internal/dto"
	"github.com/nabilfikrisp/go-crud/internal/entity"
)

type (
	// ContactRepository defines the interface for contact data storage and retrieval.
	ContactRepository interface {
		Store(ctx context.Context, contact *entity.Contact) error
		GetByID(ctx context.Context, id string) (entity.Contact, error)
		List(ctx context.Context, filter dto.ContactFilter) ([]entity.Contact, int, error)
		Update(ctx context.Context, id string, patch dto.ContactUpdate) (entity.Contact, error)
		Delete(ctx context.Context, id string) error
	}
)
