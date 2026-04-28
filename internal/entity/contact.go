package entity

import "time"

type ContactRelationship string

const (
	RelationshipFriend    ContactRelationship = "Friend"
	RelationshipFamily    ContactRelationship = "Family"
	RelationshipColleague ContactRelationship = "Colleague"
	RelationshipOther     ContactRelationship = "Other"
	RelationshipAll       ContactRelationship = ""
)

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

func (r ContactRelationship) Valid() bool {
	switch r {
	case RelationshipFriend, RelationshipFamily, RelationshipColleague, RelationshipOther:
		return true
	}
	return false
}

type ContactFilter struct {
	FirstName    *string              `json:"first_name"`
	LastName     *string              `json:"last_name"`
	Email        *string              `json:"email"`
	PhoneNumber  *string              `json:"phone_number"`
	Relationship *ContactRelationship `json:"relationship"`
	Limit        *uint64              `json:"limit"`
	Offset       *uint64              `json:"offset"`
}
