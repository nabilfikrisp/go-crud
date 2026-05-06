package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nabilfikrisp/go-crud/internal/dto"
	"github.com/nabilfikrisp/go-crud/internal/entity"
	"github.com/nabilfikrisp/go-crud/internal/usecase/contact"
	"github.com/nabilfikrisp/go-crud/pkg/ptr"
	"go.uber.org/mock/gomock"
)

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockContactRepository(ctrl)
	uc := contact.New(repo)

	t.Run("success creates contact with generated id and timestamps", func(t *testing.T) {
		req := dto.ContactCreate{
			FirstName:    "John",
			LastName:     "Doe",
			Email:        "john@example.com",
			PhoneNumber:  "+1234567890",
		}

		repo.EXPECT().Store(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, c *entity.Contact) error {
				return nil
			},
		)

		result, err := uc.Create(context.Background(), req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if result.ID == "" {
			t.Error("expected non-empty id")
		}

		if result.FirstName != "John" {
			t.Errorf("expected first name 'John', got %q", result.FirstName)
		}

		if result.LastName != "Doe" {
			t.Errorf("expected last name 'Doe', got %q", result.LastName)
		}

		if result.Email != "john@example.com" {
			t.Errorf("expected email 'john@example.com', got %q", result.Email)
		}

		if result.PhoneNumber != "+1234567890" {
			t.Errorf("expected phone number '+1234567890', got %q", result.PhoneNumber)
		}

		if result.Relationship != entity.RelationshipOther {
			t.Errorf("expected relationship 'Other', got %q", result.Relationship)
		}

		if result.CreatedAt.IsZero() {
			t.Error("expected non-zero created at")
		}

		if result.UpdatedAt.IsZero() {
			t.Error("expected non-zero updated at")
		}
	})

	t.Run("success with custom valid relationship", func(t *testing.T) {
		rel := entity.RelationshipFriend
		req := dto.ContactCreate{
			FirstName:    "Jane",
			LastName:     "Smith",
			Email:        "jane@example.com",
			PhoneNumber:  "+0987654321",
			Relationship: &rel,
		}

		repo.EXPECT().Store(gomock.Any(), gomock.Any()).Return(nil)

		result, err := uc.Create(context.Background(), req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if result.Relationship != entity.RelationshipFriend {
			t.Errorf("expected relationship 'Friend', got %q", result.Relationship)
		}
	})

	t.Run("invalid relationship returns error", func(t *testing.T) {
		invalidRel := entity.ContactRelationship("Invalid")
		req := dto.ContactCreate{
			FirstName:    "Bad",
			LastName:     "Relation",
			Email:        "bad@example.com",
			PhoneNumber:  "+1111111111",
			Relationship: &invalidRel,
		}

		_, err := uc.Create(context.Background(), req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, entity.ErrContactRelationshipInvalid) {
			t.Errorf("expected ErrContactRelationshipInvalid, got %v", err)
		}
	})

	t.Run("repo store error propagated", func(t *testing.T) {
		req := dto.ContactCreate{
			FirstName:    "Error",
			LastName:     "Test",
			Email:        "error@example.com",
			PhoneNumber:  "+2222222222",
		}

		repo.EXPECT().Store(gomock.Any(), gomock.Any()).Return(errors.New("repo error"))

		_, err := uc.Create(context.Background(), req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestGetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockContactRepository(ctrl)
	uc := contact.New(repo)

	t.Run("success returns contact from repo", func(t *testing.T) {
		id := "test-id-123"
		now := time.Now().UTC()
		expected := entity.Contact{
			ID:           id,
			FirstName:    "John",
			LastName:     "Doe",
			Email:        "john@example.com",
			PhoneNumber:  "+1234567890",
			Relationship: entity.RelationshipFriend,
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		repo.EXPECT().GetByID(gomock.Any(), id).Return(expected, nil)

		result, err := uc.GetByID(context.Background(), id)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if result.ID != id {
			t.Errorf("expected id %q, got %q", id, result.ID)
		}

		if result.FirstName != "John" {
			t.Errorf("expected first name 'John', got %q", result.FirstName)
		}
	})

	t.Run("not found error propagated from repo", func(t *testing.T) {
		id := "non-existent-id"

		repo.EXPECT().GetByID(gomock.Any(), id).Return(entity.Contact{}, entity.ErrContactNotFound)

		_, err := uc.GetByID(context.Background(), id)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, entity.ErrContactNotFound) {
			t.Errorf("expected ErrContactNotFound, got %v", err)
		}
	})

	t.Run("repo error wrapped", func(t *testing.T) {
		id := "error-id"

		repo.EXPECT().GetByID(gomock.Any(), id).Return(entity.Contact{}, errors.New("db error"))

		_, err := uc.GetByID(context.Background(), id)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockContactRepository(ctrl)
	uc := contact.New(repo)

	t.Run("success returns contacts and total", func(t *testing.T) {
		filter := dto.ContactFilter{
			Limit:  ptr.Uint64(10),
			Offset: ptr.Uint64(0),
		}
		expectedContacts := []entity.Contact{
			{ID: "1", FirstName: "John", LastName: "Doe"},
			{ID: "2", FirstName: "Jane", LastName: "Smith"},
		}

		repo.EXPECT().List(gomock.Any(), filter).Return(expectedContacts, 2, nil)

		contacts, total, err := uc.List(context.Background(), filter)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(contacts) != 2 {
			t.Errorf("expected 2 contacts, got %d", len(contacts))
		}

		if total != 2 {
			t.Errorf("expected total 2, got %d", total)
		}
	})

	t.Run("default limit applied when nil", func(t *testing.T) {
		filter := dto.ContactFilter{}

		repo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]entity.Contact{}, 0, nil)

		_, _, err := uc.List(context.Background(), filter)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("default offset applied when nil", func(t *testing.T) {
		limit := uint64(10)
		filter := dto.ContactFilter{
			Limit: &limit,
		}

		repo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]entity.Contact{}, 0, nil)

		_, _, err := uc.List(context.Background(), filter)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("invalid relationship filter returns error", func(t *testing.T) {
		invalidRel := entity.ContactRelationship("Invalid")
		filter := dto.ContactFilter{
			Relationship: &invalidRel,
		}

		_, _, err := uc.List(context.Background(), filter)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, entity.ErrContactRelationshipInvalid) {
			t.Errorf("expected ErrContactRelationshipInvalid, got %v", err)
		}
	})

	t.Run("repo error wrapped", func(t *testing.T) {
		filter := dto.ContactFilter{
			Limit:  ptr.Uint64(10),
			Offset: ptr.Uint64(0),
		}

		repo.EXPECT().List(gomock.Any(), filter).Return(nil, 0, errors.New("repo error"))

		_, _, err := uc.List(context.Background(), filter)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockContactRepository(ctrl)
	uc := contact.New(repo)

	t.Run("success updates contact", func(t *testing.T) {
		id := "test-id-123"
		now := time.Now().UTC()
		req := dto.ContactUpdate{
			FirstName: ptr.String("Updated"),
		}
		expected := entity.Contact{
			ID:           id,
			FirstName:    "Updated",
			LastName:     "Doe",
			Email:        "john@example.com",
			PhoneNumber:  "+1234567890",
			Relationship: entity.RelationshipFriend,
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		repo.EXPECT().Update(gomock.Any(), id, req).Return(expected, nil)

		result, err := uc.Update(context.Background(), id, req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if result.FirstName != "Updated" {
			t.Errorf("expected first name 'Updated', got %q", result.FirstName)
		}
	})

	t.Run("invalid relationship returns error", func(t *testing.T) {
		id := "test-id"
		invalidRel := entity.ContactRelationship("Invalid")
		req := dto.ContactUpdate{
			Relationship: &invalidRel,
		}

		_, err := uc.Update(context.Background(), id, req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, entity.ErrContactRelationshipInvalid) {
			t.Errorf("expected ErrContactRelationshipInvalid, got %v", err)
		}
	})

	t.Run("not found error propagated directly", func(t *testing.T) {
		id := "non-existent-id"
		req := dto.ContactUpdate{
			FirstName: ptr.String("Test"),
		}

		repo.EXPECT().Update(gomock.Any(), id, req).Return(entity.Contact{}, entity.ErrContactNotFound)

		_, err := uc.Update(context.Background(), id, req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, entity.ErrContactNotFound) {
			t.Errorf("expected ErrContactNotFound, got %v", err)
		}
	})

	t.Run("already exists error propagated directly", func(t *testing.T) {
		id := "test-id"
		req := dto.ContactUpdate{
			Email: ptr.String("existing@example.com"),
		}

		repo.EXPECT().Update(gomock.Any(), id, req).Return(entity.Contact{}, entity.ErrContactAlreadyExists)

		_, err := uc.Update(context.Background(), id, req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, entity.ErrContactAlreadyExists) {
			t.Errorf("expected ErrContactAlreadyExists, got %v", err)
		}
	})

	t.Run("other repo errors wrapped", func(t *testing.T) {
		id := "test-id"
		req := dto.ContactUpdate{
			FirstName: ptr.String("Test"),
		}

		repo.EXPECT().Update(gomock.Any(), id, req).Return(entity.Contact{}, errors.New("db error"))

		_, err := uc.Update(context.Background(), id, req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockContactRepository(ctrl)
	uc := contact.New(repo)

	t.Run("success delegates to repo", func(t *testing.T) {
		id := "test-id-123"

		repo.EXPECT().Delete(gomock.Any(), id).Return(nil)

		err := uc.Delete(context.Background(), id)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("repo error propagated", func(t *testing.T) {
		id := "test-id"

		repo.EXPECT().Delete(gomock.Any(), id).Return(entity.ErrContactNotFound)

		err := uc.Delete(context.Background(), id)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, entity.ErrContactNotFound) {
			t.Errorf("expected ErrContactNotFound, got %v", err)
		}
	})
}
