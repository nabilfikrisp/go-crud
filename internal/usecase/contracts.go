// Package usecase contains business logic.
package usecase

import (
	"context"

	"github.com/nabilfikrisp/go-crud/internal/dto"
	"github.com/nabilfikrisp/go-crud/internal/entity"
)

type (
	// Contact defines the interface for contact use cases.
	Contact interface {
		Create(ctx context.Context, req dto.ContactCreate) (entity.Contact, error)
		GetByID(ctx context.Context, id string) (entity.Contact, error)
		List(ctx context.Context, filter dto.ContactFilter) ([]entity.Contact, int, error)
		Update(ctx context.Context, id string, req dto.ContactUpdate) (entity.Contact, error)
		Delete(ctx context.Context, id string) error
	}
)
