package request

import "github.com/nabilfikrisp/go-crud/internal/entity"

// CreateContact -.
type CreateContact struct {
	FirstName    string                      `json:"first_name"    validate:"required,max=255"                              example:"John"`
	LastName     string                      `json:"last_name"     validate:"required,max=255"                              example:"Doe"`
	Email        string                      `json:"email"         validate:"required,email,max=255"                        example:"x6CZD@example.com"`
	PhoneNumber  string                      `json:"phone_number"  validate:"required,max=20"                               example:"+1234567890"`
	Relationship *entity.ContactRelationship `json:"relationship"  validate:"omitempty,oneof=Friend Family Colleague Other"  example:"Friend"`
} // @name v1.CreateContact

// UpdateContact -.
type UpdateContact struct {
	FirstName    *string                     `json:"first_name"    validate:"required,max=255"                              example:"Jane"`
	LastName     *string                     `json:"last_name"     validate:"required,max=255"                              example:"Smith"`
	Email        *string                     `json:"email"         validate:"required,email,max=255"                        example:"jane@example.com"`
	PhoneNumber  *string                     `json:"phone_number"  validate:"required,max=20"                               example:"+0987654321"`
	Relationship *entity.ContactRelationship `json:"relationship"  validate:"omitempty,oneof=Friend Family Colleague Other"  example:"Family"`
} // @name v1.UpdateContact
