package contact

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nabilfikrisp/go-crud/internal/entity"
	"github.com/nabilfikrisp/go-crud/internal/repo"
	"github.com/nabilfikrisp/go-crud/internal/usecase"
	"github.com/nabilfikrisp/go-crud/pkg/ptr"
)

// UseCase -.
type UseCase struct {
	repo repo.ContactRepository
}

// New -.
func New(r repo.ContactRepository) *UseCase {
	return &UseCase{repo: r}
}

// Create
func (uc *UseCase) Create(ctx context.Context, req usecase.CreateContact) (entity.Contact, error) {
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
func (uc *UseCase) List(ctx context.Context, filter entity.ContactFilter) ([]entity.Contact, int, error) {
	// Apply defaults
	if filter.Limit == nil {
		filter.Limit = ptr.Uint64(10)
	}
	if filter.Offset == nil {
		filter.Offset = ptr.Uint64(0)
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
func (uc *UseCase) Update(ctx context.Context, id string, req usecase.UpdateContact) (entity.Contact, error) {
	foundContact, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, entity.ErrContactNotFound) {
			return entity.Contact{}, fmt.Errorf("ContactUseCase - Update - entity.ErrContactNotFound: %w", err)
		}
		return entity.Contact{}, fmt.Errorf("ContactUseCase - Update - uc.repo.GetByID: %w", err)
	}

	if req.FirstName != nil {
		foundContact.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		foundContact.LastName = *req.LastName
	}
	if req.Email != nil {
		foundContact.Email = *req.Email
	}
	if req.PhoneNumber != nil {
		foundContact.PhoneNumber = *req.PhoneNumber
	}
	if req.Relationship != nil {
		if !req.Relationship.Valid() {
			return entity.Contact{}, fmt.Errorf("ContactUseCase - Create - invalid relationship: %w", entity.ErrContactRelationshipInvalid)
		}
		foundContact.Relationship = *req.Relationship
	}

	foundContact.UpdatedAt = time.Now().UTC()

	err = uc.repo.Update(ctx, &foundContact)
	if err != nil {
		return entity.Contact{}, fmt.Errorf("ContactUseCase - Update - uc.repo.Update: %w", err)
	}

	return foundContact, nil
}

func (uc *UseCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}
