package integrationtest

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"testing"
)

const (
	relationshipFriend    = "Friend"
	relationshipFamily    = "Family"
	relationshipColleague = "Colleague"
	relationshipOther     = "Other"
)

type contactResponse struct {
	ID           string `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	PhoneNumber  string `json:"phone_number"`
	Relationship string `json:"relationship"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type createContactRequest struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	PhoneNumber  string `json:"phone_number"`
	Relationship string `json:"relationship"`
}

type updateContactRequest struct {
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
	Email        string `json:"email,omitempty"`
	PhoneNumber  string `json:"phone_number,omitempty"`
	Relationship string `json:"relationship,omitempty"`
}

func httpCreateContact(t *testing.T, req createContactRequest) contactResponse {
	t.Helper()

	createBody := fmt.Sprintf(`{
		"first_name": "%s",
		"last_name": "%s",
		"email": "%s",
		"phone_number": "%s",
		"relationship": "%s"
	}`, req.FirstName, req.LastName, req.Email, req.PhoneNumber, req.Relationship)

	ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
	defer cancel()

	resp, err := doRequest(ctx, http.MethodPost, basePathV1+"/contacts/", bytes.NewBufferString(createBody))
	if err != nil {
		t.Fatalf("Create contact: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Create contact: expected 201, got %d", resp.StatusCode)
	}

	return parseJSON[contactResponse](t, resp)
}

func httpGetContact(t *testing.T, id string) contactResponse {
	t.Helper()

	ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
	defer cancel()

	resp, err := doRequest(ctx, http.MethodGet, basePathV1+"/contacts/"+id, nil)
	if err != nil {
		t.Fatalf("Get contact: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Get contact: expected 200, got %d", resp.StatusCode)
	}

	return parseJSON[contactResponse](t, resp)
}

func httpListContacts(t *testing.T, query string) (contacts []contactResponse, total int) {
	t.Helper()

	ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
	defer cancel()

	url := basePathV1 + "/contacts/"
	if query != "" {
		url += "?" + query
	}

	resp, err := doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("List contacts: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("List contacts: expected 200, got %d", resp.StatusCode)
	}

	type listResponse struct {
		Contacts []contactResponse `json:"contacts"`
		Total    int               `json:"total"`
	}

	result := parseJSON[listResponse](t, resp)
	return result.Contacts, result.Total
}

func httpUpdateContact(t *testing.T, id string, req updateContactRequest) contactResponse {
	t.Helper()

	createBody := fmt.Sprintf(`{
		"first_name": "%s",
		"last_name": "%s",
		"email": "%s",
		"phone_number": "%s",
		"relationship": "%s"
	}`, req.FirstName, req.LastName, req.Email, req.PhoneNumber, req.Relationship)

	ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
	defer cancel()

	resp, err := doRequest(ctx, http.MethodPut, basePathV1+"/contacts/"+id, bytes.NewBufferString(createBody))
	if err != nil {
		t.Fatalf("Update contact: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Update contact: expected 200, got %d", resp.StatusCode)
	}

	return parseJSON[contactResponse](t, resp)
}

func httpDeleteContact(t *testing.T, id string) {
	t.Helper()

	ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
	defer cancel()

	resp, err := doRequest(ctx, http.MethodDelete, basePathV1+"/contacts/"+id, nil)
	if err != nil {
		t.Fatalf("Delete contact: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Delete contact: expected 204, got %d", resp.StatusCode)
	}
}

func TestHTTPContactCreateV1(t *testing.T) {
	created := httpCreateContact(t, createContactRequest{
		FirstName:    "John",
		LastName:     "Doe",
		Email:        "john@example.com",
		PhoneNumber:  "+1234567890",
		Relationship: relationshipFriend,
	})

	if created.ID == "" {
		t.Fatal("expected non-empty id")
	}

	if created.FirstName != "John" {
		t.Errorf("expected first_name 'John', got %q", created.FirstName)
	}

	if created.Relationship != relationshipFriend {
		t.Errorf("expected relationship 'Friend', got %q", created.Relationship)
	}
}

func TestHTTPContactGetV1(t *testing.T) {
	created := httpCreateContact(t, createContactRequest{
		FirstName:    "Jane",
		LastName:     "Smith",
		Email:        "jane@example.com",
		PhoneNumber:  "+0987654321",
		Relationship: relationshipFamily,
	})

	got := httpGetContact(t, created.ID)

	if got.ID != created.ID {
		t.Errorf("expected id %q, got %q", created.ID, got.ID)
	}

	if got.Email != "jane@example.com" {
		t.Errorf("expected email 'jane@example.com', got %q", got.Email)
	}
}

func TestHTTPContactListV1(t *testing.T) {
	httpCreateContact(t, createContactRequest{
		FirstName:    "List",
		LastName:     "Test",
		Email:        "list@example.com",
		PhoneNumber:  "+1111111111",
		Relationship: relationshipColleague,
	})

	_, total := httpListContacts(t, "limit=10&offset=0")

	if total < 1 {
		t.Errorf("expected total >= 1, got %d", total)
	}
}

func TestHTTPContactUpdateV1(t *testing.T) {
	created := httpCreateContact(t, createContactRequest{
		FirstName:    "Update",
		LastName:     "Me",
		Email:        "update@example.com",
		PhoneNumber:  "+2222222222",
		Relationship: relationshipOther,
	})

	updated := httpUpdateContact(t, created.ID, updateContactRequest{
		FirstName:    "Updated",
		LastName:     "Name",
		Email:        "updated@example.com",
		PhoneNumber:  "+3333333333",
		Relationship: relationshipFriend,
	})

	if updated.FirstName != "Updated" {
		t.Errorf("expected first_name 'Updated', got %q", updated.FirstName)
	}

	if updated.Email != "updated@example.com" {
		t.Errorf("expected email 'updated@example.com', got %q", updated.Email)
	}
}

func TestHTTPContactPartialUpdateV1(t *testing.T) {
	created := httpCreateContact(t, createContactRequest{
		FirstName:    "Original",
		LastName:     "Name",
		Email:        "original@example.com",
		PhoneNumber:  "+1111111111",
		Relationship: relationshipFriend,
	})

	ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
	defer cancel()

	partialBody := `{"first_name": "PartialUpdated"}`
	resp, err := doRequest(ctx, http.MethodPut, basePathV1+"/contacts/"+created.ID, bytes.NewBufferString(partialBody))
	if err != nil {
		t.Fatalf("Partial update contact: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Partial update contact: expected 200, got %d", resp.StatusCode)
	}

	updated := parseJSON[contactResponse](t, resp)

	if updated.FirstName != "PartialUpdated" {
		t.Errorf("expected first_name 'PartialUpdated', got %q", updated.FirstName)
	}

	if updated.LastName != "Name" {
		t.Errorf("expected last_name 'Name' to remain unchanged, got %q", updated.LastName)
	}

	if updated.Email != "original@example.com" {
		t.Errorf("expected email 'original@example.com' to remain unchanged, got %q", updated.Email)
	}

	if updated.PhoneNumber != "+1111111111" {
		t.Errorf("expected phone_number '+1111111111' to remain unchanged, got %q", updated.PhoneNumber)
	}

	if updated.Relationship != relationshipFriend {
		t.Errorf("expected relationship 'Friend' to remain unchanged, got %q", updated.Relationship)
	}
}

func TestHTTPContactDeleteV1(t *testing.T) {
	created := httpCreateContact(t, createContactRequest{
		FirstName:    "Delete",
		LastName:     "Me",
		Email:        "delete@example.com",
		PhoneNumber:  "+4444444444",
		Relationship: relationshipFriend,
	})

	httpDeleteContact(t, created.ID)

	ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
	defer cancel()

	resp, err := doRequest(ctx, http.MethodDelete, basePathV1+"/contacts/"+created.ID, nil)
	if err != nil {
		t.Fatalf("Get deleted contact: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}

func TestHTTPContactErrorsV1(t *testing.T) {
	t.Run("create with missing required fields returns 400", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
		defer cancel()

		body := `{"first_name": "Missing"}`
		resp, err := doRequest(ctx, http.MethodPost, basePathV1+"/contacts/", bytes.NewBufferString(body))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected 400, got %d", resp.StatusCode)
		}
	})

	t.Run("get non-existent contact returns 404", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
		defer cancel()

		resp, err := doRequest(ctx, http.MethodGet, basePathV1+"/contacts/00000000-0000-0000-0000-000000000000", nil)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected 404, got %d", resp.StatusCode)
		}
	})

	t.Run("update non-existent contact returns 404", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
		defer cancel()

		body := `{"first_name": "NotFound"}`
		resp, err := doRequest(ctx, http.MethodPut, basePathV1+"/contacts/00000000-0000-0000-0000-000000000000", bytes.NewBufferString(body))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected 404, got %d", resp.StatusCode)
		}
	})

	t.Run("delete non-existent contact returns 404", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
		defer cancel()

		resp, err := doRequest(ctx, http.MethodDelete, basePathV1+"/contacts/00000000-0000-0000-0000-000000000000", nil)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected 404, got %d", resp.StatusCode)
		}
	})

	t.Run("create with invalid email returns 400", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
		defer cancel()

		body := `{"first_name": "Bad","last_name":"Email","email":"not-an-email","phone_number":"+1234567890"}`
		resp, err := doRequest(ctx, http.MethodPost, basePathV1+"/contacts/", bytes.NewBufferString(body))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected 400, got %d", resp.StatusCode)
		}
	})

	t.Run("filter by relationship", func(t *testing.T) {
		httpCreateContact(t, createContactRequest{
			FirstName:    "Filter",
			LastName:     "Test",
			Email:        "filter@example.com",
			PhoneNumber:  "+5555555555",
			Relationship: relationshipFamily,
		})

		contacts, _ := httpListContacts(t, "relationship=Family&limit=10&offset=0")

		if len(contacts) < 1 {
			t.Errorf("expected contacts filtered by relationship, got %d", len(contacts))
		}
	})
}
