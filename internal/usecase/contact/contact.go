// Package contact provides contact use case logic.
package contact

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nabilfikrisp/go-crud/internal/dto"
	"github.com/nabilfikrisp/go-crud/internal/entity"
	"github.com/nabilfikrisp/go-crud/internal/repo"
)

// UseCase -.
type UseCase struct {
	repo repo.ContactRepository
}

// New -.
func New(r repo.ContactRepository) *UseCase {
	return &UseCase{repo: r}
}

// Create handles contact creation.
func (uc *UseCase) Create(ctx context.Context, req dto.ContactCreate) (entity.Contact, error) {
	now := time.Now().UTC()

	relationship := entity.RelationshipOther
	if req.Relationship != nil {
		if !req.Relationship.Valid() {
			return entity.Contact{}, fmt.Errorf("ContactUseCase - Create - invalid relationship: %w", entity.ErrContactRelationshipInvalid)
		}
		relationship = *req.Relationship
	}

	contact := entity.Contact{
		ID:           uuid.New().String(),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		PhoneNumber:  req.PhoneNumber,
		Relationship: relationship,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err := uc.repo.Store(ctx, &contact)
	if err != nil {
		return entity.Contact{}, fmt.Errorf("ContactUseCase - Create - uc.repo.Store: %w", err)
	}

	return contact, nil
}

// GetByID -.
func (uc *UseCase) GetByID(ctx context.Context, id string) (entity.Contact, error) {
	contact, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return entity.Contact{}, fmt.Errorf("ContactUseCase - GetByID - uc.repo.GetByID: %w", err)
	}

	return contact, nil
}

// List -.
func (uc *UseCase) List(ctx context.Context, filter dto.ContactFilter) ([]entity.Contact, int, error) {
	// Apply defaults
	if filter.Limit == nil {
		filter.Limit = new(uint64(10))
	}
	if filter.Offset == nil {
		filter.Offset = new(uint64(0))
	}

	if filter.Relationship != nil && !filter.Relationship.Valid() {
		return nil, 0, fmt.Errorf("ContactUseCase - Create - invalid relationship: %w", entity.ErrContactRelationshipInvalid)
	}

	contacts, total, err := uc.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("ContactUseCase - List - uc.repo.List: %w", err)
	}

	return contacts, total, nil
}

// Update -.
func (uc *UseCase) Update(ctx context.Context, id string, req dto.ContactUpdate) (entity.Contact, error) {
	if req.Relationship != nil && !req.Relationship.Valid() {
		return entity.Contact{}, fmt.Errorf("ContactUseCase - Update - invalid relationship: %w", entity.ErrContactRelationshipInvalid)
	}

	contact, err := uc.repo.Update(ctx, id, req)
	if err != nil {
		if errors.Is(err, entity.ErrContactNotFound) || errors.Is(err, entity.ErrContactAlreadyExists) {
			return entity.Contact{}, err
		}
		return entity.Contact{}, fmt.Errorf("ContactUseCase - Update - uc.repo.Update: %w", err)
	}

	return contact, nil
}

// Delete handles contact deletion.
func (uc *UseCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}
