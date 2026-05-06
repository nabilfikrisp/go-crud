// Package entity defines contact domain models.
package entity

import "time"

// ContactRelationship represents a relationship type between contacts.
type ContactRelationship string

// Relationship types.
const (
	RelationshipFriend    ContactRelationship = "Friend"
	RelationshipFamily    ContactRelationship = "Family"
	RelationshipColleague ContactRelationship = "Colleague"
	RelationshipOther     ContactRelationship = "Other"
	RelationshipAll       ContactRelationship = ""
)

// Contact represents a contact entity.
type Contact struct {
	ID           string              `json:"id"          example:"550e8400-e29b-41d4-a716-446655440000"`
	FirstName    string              `json:"first_name"  example:"John"`
	LastName     string              `json:"last_name"   example:"Doe"`
	Email        string              `json:"email"       example:"x6CZD@example.com"`
	PhoneNumber  string              `json:"phone_number"       example:"+1234567890"`
	Relationship ContactRelationship `json:"relationship"       example:"Friend"`
	CreatedAt    time.Time           `json:"created_at"  example:"2026-01-01T00:00:00Z"`
	UpdatedAt    time.Time           `json:"updated_at"  example:"2026-01-01T00:00:00Z"`
}

// Valid returns true if the relationship is valid.
func (r ContactRelationship) Valid() bool {
	switch r {
	case RelationshipFriend, RelationshipFamily, RelationshipColleague, RelationshipOther:
		return true
	}
	return false
}
