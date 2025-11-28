package github

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAuthenticatedUser_Success(t *testing.T) {
	// Create mock user data
	mockUser := User{
		ID:        12345,
		Login:     "octocat",
		Name:      stringPtr("The Octocat"),
		Email:     stringPtr("octocat@github.com"),
		AvatarURL: "https://avatars.githubusercontent.com/u/583231",
		Bio:       stringPtr("GitHub mascot"),
		Company:   stringPtr("@github"),
		Location:  stringPtr("San Francisco"),
	}

	// Create a test server that returns the mock user
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request headers
		if auth := r.Header.Get("Authorization"); auth != "Bearer test-token" {
			t.Errorf("Expected Authorization header 'Bearer test-token', got '%s'", auth)
		}
		if accept := r.Header.Get("Accept"); accept != "application/vnd.github+json" {
			t.Errorf("Expected Accept header 'application/vnd.github+json', got '%s'", accept)
		}
		if version := r.Header.Get("X-GitHub-Api-Version"); version != "2022-11-28" {
			t.Errorf("Expected X-GitHub-Api-Version header '2022-11-28', got '%s'", version)
		}

		// Verify endpoint
		if r.URL.Path != "/user" {
			t.Errorf("Expected path '/user', got '%s'", r.URL.Path)
		}

		// Return mock user
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockUser)
	}))
	defer server.Close()

	// Create client pointing to test server
	client := NewClient("test-token")
	client.baseURL = server.URL

	// Call GetAuthenticatedUser
	user, err := client.GetAuthenticatedUser()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify response
	if user.ID != mockUser.ID {
		t.Errorf("Expected ID %d, got %d", mockUser.ID, user.ID)
	}
	if user.Login != mockUser.Login {
		t.Errorf("Expected Login %s, got %s", mockUser.Login, user.Login)
	}
	if *user.Name != *mockUser.Name {
		t.Errorf("Expected Name %s, got %s", *mockUser.Name, *user.Name)
	}
	if *user.Email != *mockUser.Email {
		t.Errorf("Expected Email %s, got %s", *mockUser.Email, *user.Email)
	}
	if user.AvatarURL != mockUser.AvatarURL {
		t.Errorf("Expected AvatarURL %s, got %s", mockUser.AvatarURL, user.AvatarURL)
	}
}

func TestGetAuthenticatedUser_Unauthorized(t *testing.T) {
	// Create a test server that returns 401
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Bad credentials",
		})
	}))
	defer server.Close()

	// Create client pointing to test server
	client := NewClient("invalid-token")
	client.baseURL = server.URL

	// Call GetAuthenticatedUser
	_, err := client.GetAuthenticatedUser()
	if err == nil {
		t.Fatal("Expected error for unauthorized request, got nil")
	}
}

func TestGetAuthenticatedUser_NullableFields(t *testing.T) {
	// Create mock user with null optional fields
	mockUser := User{
		ID:        12345,
		Login:     "octocat",
		Name:      nil,
		Email:     nil,
		AvatarURL: "https://avatars.githubusercontent.com/u/583231",
		Bio:       nil,
		Company:   nil,
		Location:  nil,
	}

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockUser)
	}))
	defer server.Close()

	// Create client pointing to test server
	client := NewClient("test-token")
	client.baseURL = server.URL

	// Call GetAuthenticatedUser
	user, err := client.GetAuthenticatedUser()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify nullable fields are nil
	if user.Name != nil {
		t.Errorf("Expected Name to be nil, got %v", user.Name)
	}
	if user.Email != nil {
		t.Errorf("Expected Email to be nil, got %v", user.Email)
	}
	if user.Bio != nil {
		t.Errorf("Expected Bio to be nil, got %v", user.Bio)
	}
}

func TestGetUser_Success(t *testing.T) {
	// Create mock user data
	mockUser := User{
		ID:        583231,
		Login:     "octocat",
		Name:      stringPtr("The Octocat"),
		Email:     stringPtr("octocat@github.com"),
		AvatarURL: "https://avatars.githubusercontent.com/u/583231",
		Bio:       stringPtr("GitHub mascot"),
		Company:   stringPtr("@github"),
		Location:  stringPtr("San Francisco"),
	}

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request headers
		if auth := r.Header.Get("Authorization"); auth != "Bearer test-token" {
			t.Errorf("Expected Authorization header 'Bearer test-token', got '%s'", auth)
		}

		// Verify endpoint
		expectedPath := "/users/octocat"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path '%s', got '%s'", expectedPath, r.URL.Path)
		}

		// Return mock user
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockUser)
	}))
	defer server.Close()

	// Create client pointing to test server
	client := NewClient("test-token")
	client.baseURL = server.URL

	// Call GetUser
	user, err := client.GetUser("octocat")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify response
	if user.Login != mockUser.Login {
		t.Errorf("Expected Login %s, got %s", mockUser.Login, user.Login)
	}
	if *user.Name != *mockUser.Name {
		t.Errorf("Expected Name %s, got %s", *mockUser.Name, *user.Name)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	// Create a test server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Not Found",
		})
	}))
	defer server.Close()

	// Create client pointing to test server
	client := NewClient("test-token")
	client.baseURL = server.URL

	// Call GetUser
	_, err := client.GetUser("nonexistent")
	if err == nil {
		t.Fatal("Expected error for non-existent user, got nil")
	}

	// Verify error message contains "not found"
	expectedMsg := "user not found: nonexistent"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func stringPtr(s string) *string {
	return &s
}
