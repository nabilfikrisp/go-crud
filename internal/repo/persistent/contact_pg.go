package persistent

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/nabilfikrisp/go-crud/internal/dto"
	"github.com/nabilfikrisp/go-crud/internal/entity"
	"github.com/nabilfikrisp/go-crud/pkg/postgres"
)

// ContactPGRepo -.
type ContactPGRepo struct {
	*postgres.Postgres
}

// NewContactPGRepo -.
func NewContactPGRepo(pg *postgres.Postgres) *ContactPGRepo {
	return &ContactPGRepo{pg}
}

// Store -.
func (r *ContactPGRepo) Store(ctx context.Context, contact *entity.Contact) error {
	sql, args, err := r.Builder.
		Insert("contacts").
		Columns("id, first_name, last_name, email, phone_number, relationship, created_at, updated_at").
		Values(contact.ID, contact.FirstName, contact.LastName, contact.Email, contact.PhoneNumber, contact.Relationship, contact.CreatedAt, contact.UpdatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("ContactRepo - Store - r.Builder: %w", err)
	}
	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("ContactRepo - Store - r.Exec: %w", err)
	}

	return nil
}

// GetByID -.
func (r *ContactPGRepo) GetByID(ctx context.Context, id string) (entity.Contact, error) {
	sql, args, err := r.Builder.
		Select("id, first_name, last_name, email, phone_number, relationship, created_at, updated_at").
		From("contacts").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return entity.Contact{}, fmt.Errorf("ContactRepo - GetByID - r.Builder: %w", err)
	}

	var contact entity.Contact

	err = r.Pool.QueryRow(ctx, sql, args...).
		Scan(&contact.ID, &contact.FirstName, &contact.LastName, &contact.Email, &contact.PhoneNumber, &contact.Relationship, &contact.CreatedAt, &contact.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Contact{}, entity.ErrContactNotFound
		}

		return entity.Contact{}, fmt.Errorf("ContactRepo - GetByID - r.Pool.QueryRow: %w", err)
	}

	return contact, nil
}

// List -.
func (r *ContactPGRepo) List(ctx context.Context, filter dto.ContactFilter) ([]entity.Contact, int, error) {
	conditions := sq.And{}

	if filter.FirstName != nil {
		conditions = append(conditions, sq.ILike{"first_name": "%" + *filter.FirstName + "%"})
	}
	if filter.LastName != nil {
		conditions = append(conditions, sq.ILike{"last_name": "%" + *filter.LastName + "%"})
	}
	if filter.Email != nil {
		conditions = append(conditions, sq.ILike{"email": "%" + *filter.Email + "%"})
	}
	if filter.PhoneNumber != nil {
		conditions = append(conditions, sq.Like{"phone_number": "%" + *filter.PhoneNumber + "%"})
	}
	if filter.Relationship != nil {
		conditions = append(conditions, sq.Eq{"relationship": *filter.Relationship})
	}

	countSQL, countArgs, err := r.Builder.
		Select("COUNT(*)").
		From("contacts").
		Where(conditions).
		ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("ContactRepo - List - r.Builder count: %w", err)
	}

	var totalMatches int
	err = r.Pool.QueryRow(ctx, countSQL, countArgs...).Scan(&totalMatches)
	if err != nil {
		return nil, 0, fmt.Errorf("ContactRepo - List - r.Pool.QueryRow count: %w", err)
	}

	if totalMatches == 0 {
		return []entity.Contact{}, 0, nil
	}

	if filter.Offset != nil {
		if totalMatches < 0 {
			return nil, 0, fmt.Errorf("ContactRepo - List - unexpected negative total matches: %d", totalMatches)
		}

		if *filter.Offset >= uint64(totalMatches) {
			return []entity.Contact{}, totalMatches, nil
		}
	}

	selectBuilder := r.Builder.
		Select("id, first_name, last_name, email, phone_number, relationship, created_at, updated_at").
		From("contacts").
		Where(conditions).
		OrderBy("created_at DESC")

	if filter.Limit != nil {
		selectBuilder = selectBuilder.Limit(*filter.Limit)
	}
	if filter.Offset != nil {
		selectBuilder = selectBuilder.Offset(*filter.Offset)
	}

	sql, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("ContactRepo - List - r.Builder select: %w", err)
	}
	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("ContactRepo - List - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	contacts := make([]entity.Contact, 0, totalMatches)
	for rows.Next() {
		var c entity.Contact
		if err := rows.Scan(
			&c.ID,
			&c.FirstName,
			&c.LastName,
			&c.Email,
			&c.PhoneNumber,
			&c.Relationship,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("ContactRepo - List - rows.Scan: %w", err)
		}
		contacts = append(contacts, c)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("ContactRepo - List - rows.Err: %w", err)
	}

	return contacts, totalMatches, nil
}

// Update atomically applies a partial update to the contact identified by id.
// Returns entity.ErrContactNotFound if no row matches, or
// entity.ErrContactAlreadyExists on email unique violation.
func (r *ContactPGRepo) Update(ctx context.Context, id string, patch dto.ContactUpdate) (entity.Contact, error) {
	builder := r.Builder.Update("contacts")

	hasUpdate := false
	if patch.FirstName != nil {
		builder = builder.Set("first_name", *patch.FirstName)
		hasUpdate = true
	}
	if patch.LastName != nil {
		builder = builder.Set("last_name", *patch.LastName)
		hasUpdate = true
	}
	if patch.Email != nil {
		builder = builder.Set("email", *patch.Email)
		hasUpdate = true
	}
	if patch.PhoneNumber != nil {
		builder = builder.Set("phone_number", *patch.PhoneNumber)
		hasUpdate = true
	}
	if patch.Relationship != nil {
		builder = builder.Set("relationship", *patch.Relationship)
		hasUpdate = true
	}

	// If no fields to update, just fetch the current row.
	if !hasUpdate {
		return r.GetByID(ctx, id)
	}

	builder = builder.
		Set("updated_at", time.Now().UTC()).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING id, first_name, last_name, email, phone_number, relationship, created_at, updated_at")

	sql, args, err := builder.ToSql()
	if err != nil {
		return entity.Contact{}, fmt.Errorf("ContactRepo - Update - r.Builder: %w", err)
	}

	var c entity.Contact
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&c.ID,
		&c.FirstName,
		&c.LastName,
		&c.Email,
		&c.PhoneNumber,
		&c.Relationship,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Contact{}, entity.ErrContactNotFound
		}
		if IsUniqueViolation(err) {
			return entity.Contact{}, entity.ErrContactAlreadyExists
		}
		return entity.Contact{}, fmt.Errorf("ContactRepo - Update - r.Pool.QueryRow: %w", err)
	}

	return c, nil
}

// Delete removes a contact by ID, returning entity.ErrContactNotFound if no match is found.
func (r *ContactPGRepo) Delete(ctx context.Context, id string) error {
	sql, args, err := r.Builder.
		Delete("contacts").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("ContactRepo - Delete - r.Builder: %w", err)
	}

	cmdTag, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("ContactRepo - Delete - r.Pool.Exec: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return entity.ErrContactNotFound
	}

	return nil
}
