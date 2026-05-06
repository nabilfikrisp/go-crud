// Package dto provides data transfer objects.
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

	// ContactCreate represents a contact creation request.
	ContactCreate struct {
		FirstName    string
		LastName     string
		Email        string
		PhoneNumber  string
		Relationship *entity.ContactRelationship
	}

	// ContactUpdate represents a contact update request.
	ContactUpdate struct {
		FirstName    *string
		LastName     *string
		Email        *string
		PhoneNumber  *string
		Relationship *entity.ContactRelationship
	}
)
