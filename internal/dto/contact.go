package dto

import "github.com/nabilfikrisp/go-crud/internal/entity"

type (
	// ContactFilter defines the structure for filtering contacts.
	ContactFilter struct {
		FirstName    *string
		LastName     *string
		Email        *string
		PhoneNumber  *string
		Relationship *entity.ContactRelationship
		Limit        *uint64
		Offset       *uint64
	}

	// CreateContactRequest defines the structure for creating a new contact.
	ContactCreate struct {
		FirstName    string
		LastName     string
		Email        string
		PhoneNumber  string
		Relationship *entity.ContactRelationship
	}

	// UpdateContactRequest defines the structure for updating an existing contact.
	ContactUpdate struct {
		FirstName    *string
		LastName     *string
		Email        *string
		PhoneNumber  *string
		Relationship *entity.ContactRelationship
	}
)
