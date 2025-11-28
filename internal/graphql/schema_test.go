package graphql

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/srt32/ghgraph/internal/github"
)

func TestViewerQuery_AllFields(t *testing.T) {
	// Create mock GitHub API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockUser := github.User{
			ID:        583231,
			Login:     "octocat",
			Name:      stringPtr("The Octocat"),
			Email:     stringPtr("octocat@github.com"),
			AvatarURL: "https://avatars.githubusercontent.com/u/583231",
			Bio:       stringPtr("GitHub mascot"),
			Company:   stringPtr("@github"),
			Location:  stringPtr("San Francisco"),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockUser)
	}))
	defer server.Close()

	// Create GitHub client pointing to mock server
	client := github.NewClient("test-token")
	client.SetBaseURL(server.URL)

	// Create schema
	schema, err := NewSchema(client)
	if err != nil {
		t.Fatalf("Failed to create schema: %v", err)
	}

	// Execute query with all fields
	query := `
		query {
			viewer {
				id
				login
				name
				email
				avatarUrl
				bio
				company
				location
			}
		}
	`

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})

	// Check for errors
	if len(result.Errors) > 0 {
		t.Fatalf("GraphQL errors: %v", result.Errors)
	}

	// Verify response structure
	data, ok := result.Data.(map[string]interface{})
	if !ok {
		t.Fatal("Expected data to be a map")
	}

	viewer, ok := data["viewer"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected viewer to be a map")
	}

	// Verify all fields
	if viewer["login"] != "octocat" {
		t.Errorf("Expected login 'octocat', got %v", viewer["login"])
	}
	if viewer["name"] != "The Octocat" {
		t.Errorf("Expected name 'The Octocat', got %v", viewer["name"])
	}
	if viewer["email"] != "octocat@github.com" {
		t.Errorf("Expected email 'octocat@github.com', got %v", viewer["email"])
	}
	if viewer["avatarUrl"] != "https://avatars.githubusercontent.com/u/583231" {
		t.Errorf("Expected avatarUrl, got %v", viewer["avatarUrl"])
	}
	if viewer["bio"] != "GitHub mascot" {
		t.Errorf("Expected bio 'GitHub mascot', got %v", viewer["bio"])
	}
	if viewer["company"] != "@github" {
		t.Errorf("Expected company '@github', got %v", viewer["company"])
	}
	if viewer["location"] != "San Francisco" {
		t.Errorf("Expected location 'San Francisco', got %v", viewer["location"])
	}
}

func TestViewerQuery_PartialFields(t *testing.T) {
	// Create mock GitHub API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockUser := github.User{
			ID:        583231,
			Login:     "octocat",
			Name:      stringPtr("The Octocat"),
			Email:     stringPtr("octocat@github.com"),
			AvatarURL: "https://avatars.githubusercontent.com/u/583231",
			Bio:       nil,
			Company:   nil,
			Location:  nil,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockUser)
	}))
	defer server.Close()

	// Create GitHub client pointing to mock server
	client := github.NewClient("test-token")
	client.SetBaseURL(server.URL)

	// Create schema
	schema, err := NewSchema(client)
	if err != nil {
		t.Fatalf("Failed to create schema: %v", err)
	}

	// Execute query with only some fields
	query := `
		query {
			viewer {
				login
				name
			}
		}
	`

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})

	// Check for errors
	if len(result.Errors) > 0 {
		t.Fatalf("GraphQL errors: %v", result.Errors)
	}

	// Verify response structure
	data, ok := result.Data.(map[string]interface{})
	if !ok {
		t.Fatal("Expected data to be a map")
	}

	viewer, ok := data["viewer"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected viewer to be a map")
	}

	// Verify requested fields are present
	if viewer["login"] != "octocat" {
		t.Errorf("Expected login 'octocat', got %v", viewer["login"])
	}
	if viewer["name"] != "The Octocat" {
		t.Errorf("Expected name 'The Octocat', got %v", viewer["name"])
	}

	// Verify unrequested fields are not present
	if _, exists := viewer["email"]; exists {
		t.Error("Expected email to not be in response")
	}
	if _, exists := viewer["bio"]; exists {
		t.Error("Expected bio to not be in response")
	}
}

func TestViewerQuery_NullableFields(t *testing.T) {
	// Create mock GitHub API server with null fields
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockUser := github.User{
			ID:        583231,
			Login:     "octocat",
			Name:      nil,
			Email:     nil,
			AvatarURL: "https://avatars.githubusercontent.com/u/583231",
			Bio:       nil,
			Company:   nil,
			Location:  nil,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockUser)
	}))
	defer server.Close()

	// Create GitHub client pointing to mock server
	client := github.NewClient("test-token")
	client.SetBaseURL(server.URL)

	// Create schema
	schema, err := NewSchema(client)
	if err != nil {
		t.Fatalf("Failed to create schema: %v", err)
	}

	// Execute query
	query := `
		query {
			viewer {
				login
				name
				bio
				company
				location
			}
		}
	`

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})

	// Check for errors
	if len(result.Errors) > 0 {
		t.Fatalf("GraphQL errors: %v", result.Errors)
	}

	// Verify response structure
	data, ok := result.Data.(map[string]interface{})
	if !ok {
		t.Fatal("Expected data to be a map")
	}

	viewer, ok := data["viewer"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected viewer to be a map")
	}

	// Verify nullable fields are null
	if viewer["name"] != nil {
		t.Errorf("Expected name to be null, got %v", viewer["name"])
	}
	if viewer["bio"] != nil {
		t.Errorf("Expected bio to be null, got %v", viewer["bio"])
	}
	if viewer["company"] != nil {
		t.Errorf("Expected company to be null, got %v", viewer["company"])
	}
	if viewer["location"] != nil {
		t.Errorf("Expected location to be null, got %v", viewer["location"])
	}
}

func TestUserQuery_Success(t *testing.T) {
	// Create mock GitHub API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the correct endpoint is called
		if r.URL.Path != "/users/octocat" {
			t.Errorf("Expected path '/users/octocat', got '%s'", r.URL.Path)
		}

		mockUser := github.User{
			ID:        583231,
			Login:     "octocat",
			Name:      stringPtr("The Octocat"),
			Email:     stringPtr("octocat@github.com"),
			AvatarURL: "https://avatars.githubusercontent.com/u/583231",
			Bio:       stringPtr("GitHub mascot"),
			Company:   stringPtr("@github"),
			Location:  stringPtr("San Francisco"),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockUser)
	}))
	defer server.Close()

	// Create GitHub client pointing to mock server
	client := github.NewClient("test-token")
	client.SetBaseURL(server.URL)

	// Create schema
	schema, err := NewSchema(client)
	if err != nil {
		t.Fatalf("Failed to create schema: %v", err)
	}

	// Execute query
	query := `
		query {
			user(login: "octocat") {
				login
				name
				email
				bio
			}
		}
	`

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})

	// Check for errors
	if len(result.Errors) > 0 {
		t.Fatalf("GraphQL errors: %v", result.Errors)
	}

	// Verify response structure
	data, ok := result.Data.(map[string]interface{})
	if !ok {
		t.Fatal("Expected data to be a map")
	}

	user, ok := data["user"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected user to be a map")
	}

	// Verify fields
	if user["login"] != "octocat" {
		t.Errorf("Expected login 'octocat', got %v", user["login"])
	}
	if user["name"] != "The Octocat" {
		t.Errorf("Expected name 'The Octocat', got %v", user["name"])
	}
	if user["email"] != "octocat@github.com" {
		t.Errorf("Expected email 'octocat@github.com', got %v", user["email"])
	}
}

func TestUserQuery_NotFound(t *testing.T) {
	// Create mock GitHub API server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Not Found",
		})
	}))
	defer server.Close()

	// Create GitHub client pointing to mock server
	client := github.NewClient("test-token")
	client.SetBaseURL(server.URL)

	// Create schema
	schema, err := NewSchema(client)
	if err != nil {
		t.Fatalf("Failed to create schema: %v", err)
	}

	// Execute query
	query := `
		query {
			user(login: "nonexistent") {
				login
			}
		}
	`

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})

	// Should have errors
	if len(result.Errors) == 0 {
		t.Fatal("Expected GraphQL errors for non-existent user")
	}

	// Verify error message
	errorMsg := result.Errors[0].Message
	if errorMsg != "user not found: nonexistent" {
		t.Errorf("Expected error 'user not found: nonexistent', got '%s'", errorMsg)
	}
}

func TestUserQuery_DifferentUsers(t *testing.T) {
	// Test that different login arguments fetch different users
	callCount := 0
	expectedLogins := []string{"user1", "user2"}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var login string
		if r.URL.Path == "/users/user1" {
			login = "user1"
		} else if r.URL.Path == "/users/user2" {
			login = "user2"
		} else {
			t.Errorf("Unexpected path: %s", r.URL.Path)
		}

		mockUser := github.User{
			ID:        int64(callCount + 1),
			Login:     login,
			Name:      stringPtr(login + " name"),
			Email:     stringPtr(login + "@example.com"),
			AvatarURL: "https://avatars.githubusercontent.com/u/" + login,
		}
		callCount++
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockUser)
	}))
	defer server.Close()

	client := github.NewClient("test-token")
	client.SetBaseURL(server.URL)

	schema, err := NewSchema(client)
	if err != nil {
		t.Fatalf("Failed to create schema: %v", err)
	}

	// Query for two different users
	for _, expectedLogin := range expectedLogins {
		query := `
			query {
				user(login: "` + expectedLogin + `") {
					login
				}
			}
		`

		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: query,
		})

		if len(result.Errors) > 0 {
			t.Fatalf("GraphQL errors: %v", result.Errors)
		}

		data := result.Data.(map[string]interface{})
		user := data["user"].(map[string]interface{})

		if user["login"] != expectedLogin {
			t.Errorf("Expected login '%s', got %v", expectedLogin, user["login"])
		}
	}
}

func TestRepositoryQuery_Success(t *testing.T) {
	// Create mock GitHub API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the correct endpoint is called
		if r.URL.Path != "/repos/octocat/hello-world" {
			t.Errorf("Expected path '/repos/octocat/hello-world', got '%s'", r.URL.Path)
		}

		mockRepo := github.Repository{
			ID:              12345,
			Name:            "hello-world",
			FullName:        "octocat/hello-world",
			Description:     stringPtr("My first repository"),
			Private:         false,
			HTMLURL:         "https://github.com/octocat/hello-world",
			StargazersCount: 100,
			ForksCount:      25,
			DefaultBranch:   "main",
			Owner: github.User{
				ID:    583231,
				Login: "octocat",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockRepo)
	}))
	defer server.Close()

	// Create GitHub client pointing to mock server
	client := github.NewClient("test-token")
	client.SetBaseURL(server.URL)

	// Create schema
	schema, err := NewSchema(client)
	if err != nil {
		t.Fatalf("Failed to create schema: %v", err)
	}

	// Execute query
	query := `
		query {
			repository(owner: "octocat", name: "hello-world") {
				name
				nameWithOwner
				description
				isPrivate
				url
				stargazerCount
				forkCount
				defaultBranchRef
			}
		}
	`

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})

	// Check for errors
	if len(result.Errors) > 0 {
		t.Fatalf("GraphQL errors: %v", result.Errors)
	}

	// Verify response structure
	data, ok := result.Data.(map[string]interface{})
	if !ok {
		t.Fatal("Expected data to be a map")
	}

	repo, ok := data["repository"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected repository to be a map")
	}

	// Verify fields
	if repo["name"] != "hello-world" {
		t.Errorf("Expected name 'hello-world', got %v", repo["name"])
	}
	if repo["nameWithOwner"] != "octocat/hello-world" {
		t.Errorf("Expected nameWithOwner 'octocat/hello-world', got %v", repo["nameWithOwner"])
	}
	if repo["description"] != "My first repository" {
		t.Errorf("Expected description 'My first repository', got %v", repo["description"])
	}
	if repo["isPrivate"] != false {
		t.Errorf("Expected isPrivate false, got %v", repo["isPrivate"])
	}
	if repo["stargazerCount"] != 100 {
		t.Errorf("Expected stargazerCount 100, got %v", repo["stargazerCount"])
	}
	if repo["defaultBranchRef"] != "main" {
		t.Errorf("Expected defaultBranchRef 'main', got %v", repo["defaultBranchRef"])
	}
}

func TestRepositoryQuery_WithOwner(t *testing.T) {
	// Test querying nested owner field
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockRepo := github.Repository{
			ID:       12345,
			Name:     "hello-world",
			FullName: "octocat/hello-world",
			Owner: github.User{
				ID:    583231,
				Login: "octocat",
				Name:  stringPtr("The Octocat"),
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockRepo)
	}))
	defer server.Close()

	client := github.NewClient("test-token")
	client.SetBaseURL(server.URL)

	schema, err := NewSchema(client)
	if err != nil {
		t.Fatalf("Failed to create schema: %v", err)
	}

	// Query with nested owner fields
	query := `
		query {
			repository(owner: "octocat", name: "hello-world") {
				name
				owner {
					login
					name
				}
			}
		}
	`

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})

	if len(result.Errors) > 0 {
		t.Fatalf("GraphQL errors: %v", result.Errors)
	}

	data := result.Data.(map[string]interface{})
	repo := data["repository"].(map[string]interface{})
	owner := repo["owner"].(map[string]interface{})

	if owner["login"] != "octocat" {
		t.Errorf("Expected owner login 'octocat', got %v", owner["login"])
	}
	if owner["name"] != "The Octocat" {
		t.Errorf("Expected owner name 'The Octocat', got %v", owner["name"])
	}
}

func TestRepositoryQuery_NotFound(t *testing.T) {
	// Create mock GitHub API server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Not Found",
		})
	}))
	defer server.Close()

	client := github.NewClient("test-token")
	client.SetBaseURL(server.URL)

	schema, err := NewSchema(client)
	if err != nil {
		t.Fatalf("Failed to create schema: %v", err)
	}

	query := `
		query {
			repository(owner: "nonexistent", name: "repo") {
				name
			}
		}
	`

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})

	// Should have errors
	if len(result.Errors) == 0 {
		t.Fatal("Expected GraphQL errors for non-existent repository")
	}

	errorMsg := result.Errors[0].Message
	if errorMsg != "repository not found: nonexistent/repo" {
		t.Errorf("Expected error 'repository not found: nonexistent/repo', got '%s'", errorMsg)
	}
}

func stringPtr(s string) *string {
	return &s
}
