package bobuild

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func TestGetUser(t *testing.T) {
	mockUser := testUser{
		ID:    "123",
		Name:  "Alice",
		Email: "alice@example.com",
	}

	// Create a test server that returns the mock user
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/_api/users/123" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockUser)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-api-key")
	user, err := Get[testUser](client, "/users/123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user.Name != mockUser.Name {
		t.Errorf("Expected name %s, got %s", mockUser.Name, user.Name)
	}
}

func TestCreateUser(t *testing.T) {
	inputUser := &testUser{
		Name:  "Bob",
		Email: "bob@example.com",
	}

	// Create a test server that echos the request body back as the response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/_api/insert" {
			http.NotFound(w, r)
			return
		}

		r.Body.Close()

		// Just echo the request as response for testing
		w.Header().Set("Content-Type", "application/json")
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-api-key")
	_, err := Insert(client, "/insert", inputUser)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
