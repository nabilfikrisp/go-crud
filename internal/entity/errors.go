// Package entity defines contact domain errors.
package entity

import "errors"

// Contact errors.
var (
	ErrContactAlreadyExists       = errors.New("contact already exists")
	ErrContactNotFound            = errors.New("contact not found")
	ErrContactRelationshipInvalid = errors.New("invalid contact relationship")
)
