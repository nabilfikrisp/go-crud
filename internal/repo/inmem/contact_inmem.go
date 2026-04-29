package inmem

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/nabilfikrisp/go-crud/internal/entity"
)

// ContactInMemRepo is an in-memory implementation of the Contact repository.
type ContactInMemRepo struct {
	mu       sync.RWMutex
	contacts map[string]entity.Contact
	emails   map[string]bool
}

// NewContactInMemRepo creates a new instance of ContactInMemRepo.
func NewContactInMemRepo() *ContactInMemRepo {
	return &ContactInMemRepo{
		contacts: make(map[string]entity.Contact),
		emails:   make(map[string]bool),
	}
}

// Store adds a new contact to the in-memory repository. It checks for email uniqueness before adding.
func (r *ContactInMemRepo) Store(ctx context.Context, contact *entity.Contact) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.emails[contact.Email] {
		return fmt.Errorf("ContactInMemRepo - Store - uniqueness check: %w", entity.ErrContactAlreadyExists)
	}

	r.contacts[contact.ID] = *contact
	return nil
}

// GetByID retrieves a contact by its ID. It returns an error if the contact is not found.
func (r *ContactInMemRepo) GetByID(ctx context.Context, id string) (entity.Contact, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	contact, ok := r.contacts[id]
	if !ok {
		return entity.Contact{}, fmt.Errorf("ContatctInMemRepoo - GetByID - contact not found: %w", entity.ErrContactNotFound)
	}

	return contact, nil
}

// List retrieves a list of contacts based on the provided filter criteria. It supports pagination and returns the total number of matches.
func (r *ContactInMemRepo) List(ctx context.Context, filter entity.ContactFilter) ([]entity.Contact, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filteredContacts []entity.Contact

	for _, c := range r.contacts {
		if filter.FirstName != nil {
			if !strings.Contains(strings.ToLower(c.FirstName), strings.ToLower(*filter.FirstName)) {
				continue
			}
		}

		if filter.LastName != nil {
			if !strings.Contains(strings.ToLower(c.LastName), strings.ToLower(*filter.LastName)) {
				continue
			}
		}

		if filter.Email != nil {
			if !strings.Contains(strings.ToLower(c.Email), strings.ToLower(*filter.Email)) {
				continue
			}
		}

		if filter.PhoneNumber != nil {
			if !strings.Contains(c.PhoneNumber, *filter.PhoneNumber) {
				continue
			}
		}

		if filter.Relationship != nil && c.Relationship != *filter.Relationship {
			continue
		}

		filteredContacts = append(filteredContacts, c)
	}

	totalMatches := len(filteredContacts)
	total := uint64(totalMatches)
	offset := *filter.Offset
	limit := *filter.Limit

	if offset >= total {
		return []entity.Contact{}, totalMatches, nil
	}

	end := min(offset+limit, total)

	return filteredContacts[offset:end], totalMatches, nil
}

// Update modifies an existing contact by ID in the in-memory store, returning entity.ErrContactNotFound if no match is found.
func (r *ContactInMemRepo) Update(ctx context.Context, contact *entity.Contact) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.contacts[contact.ID]; !ok {
		return fmt.Errorf("ContatctInMemRepoo - Update - contact not found: %w", entity.ErrContactNotFound)
	}

	r.contacts[contact.ID] = *contact

	return nil
}

// Delete removes a contact by ID from the in-memory store, returning entity.ErrContactNotFound if the specified ID does not exist.
func (r *ContactInMemRepo) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.contacts[id]; !ok {
		return fmt.Errorf("ContactInMemRepo - Delete - contact not found: %w", entity.ErrContactNotFound)
	}

	delete(r.contacts, id)
	return nil
}
