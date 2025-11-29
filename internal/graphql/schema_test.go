package graphql

import (
	"context"
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

	ctx := context.WithValue(context.Background(), githubClientKey, client)
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
		Context:       ctx,
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

	ctx := context.WithValue(context.Background(), githubClientKey, client)
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
		Context:       ctx,
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

	ctx := context.WithValue(context.Background(), githubClientKey, client)
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
		Context:       ctx,
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

	ctx := context.WithValue(context.Background(), githubClientKey, client)
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
		Context:       ctx,
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

	ctx := context.WithValue(context.Background(), githubClientKey, client)
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
		Context:       ctx,
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

	ctx := context.WithValue(context.Background(), githubClientKey, client)
		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: query,
		Context:       ctx,
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

	ctx := context.WithValue(context.Background(), githubClientKey, client)
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
		Context:       ctx,
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

	ctx := context.WithValue(context.Background(), githubClientKey, client)
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
		Context:       ctx,
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

	ctx := context.WithValue(context.Background(), githubClientKey, client)
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
		Context:       ctx,
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

func TestOrganizationQuery_Success(t *testing.T) {
	// Create mock GitHub API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the correct endpoint is called
		if r.URL.Path != "/orgs/github" {
			t.Errorf("Expected path '/orgs/github', got '%s'", r.URL.Path)
		}

		mockOrg := github.Organization{
			ID:          1,
			Login:       "github",
			Name:        stringPtr("GitHub"),
			Description: stringPtr("How people build software"),
			AvatarURL:   "https://avatars.githubusercontent.com/u/9919",
			HTMLURL:     "https://github.com/github",
			Email:       stringPtr("support@github.com"),
			Location:    stringPtr("San Francisco, CA"),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockOrg)
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
			organization(login: "github") {
				login
				name
				description
				avatarUrl
				url
				email
				location
			}
		}
	`

	ctx := context.WithValue(context.Background(), githubClientKey, client)
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
		Context:       ctx,
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

	org, ok := data["organization"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected organization to be a map")
	}

	// Verify fields
	if org["login"] != "github" {
		t.Errorf("Expected login 'github', got %v", org["login"])
	}
	if org["name"] != "GitHub" {
		t.Errorf("Expected name 'GitHub', got %v", org["name"])
	}
	if org["description"] != "How people build software" {
		t.Errorf("Expected description, got %v", org["description"])
	}
	if org["email"] != "support@github.com" {
		t.Errorf("Expected email 'support@github.com', got %v", org["email"])
	}
}

func TestOrganizationQuery_NotFound(t *testing.T) {
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
			organization(login: "nonexistent") {
				login
			}
		}
	`

	ctx := context.WithValue(context.Background(), githubClientKey, client)
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
		Context:       ctx,
	})

	// Should have errors
	if len(result.Errors) == 0 {
		t.Fatal("Expected GraphQL errors for non-existent organization")
	}

	errorMsg := result.Errors[0].Message
	if errorMsg != "organization not found: nonexistent" {
		t.Errorf("Expected error 'organization not found: nonexistent', got '%s'", errorMsg)
	}
}

func TestViewerRepositoriesQuery_Success(t *testing.T) {
	// Track which endpoints were called
	var userCalled, reposCalled bool

	// Create mock GitHub API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.URL.Path == "/user" {
			userCalled = true
			mockUser := github.User{
				ID:    12345,
				Login: "testuser",
			}
			json.NewEncoder(w).Encode(mockUser)
		} else if r.URL.Path == "/users/testuser/repos" {
			reposCalled = true
			// Verify pagination parameters
			if r.URL.Query().Get("per_page") != "2" {
				t.Errorf("Expected per_page=2, got %s", r.URL.Query().Get("per_page"))
			}
			if r.URL.Query().Get("page") != "1" {
				t.Errorf("Expected page=1, got %s", r.URL.Query().Get("page"))
			}

			mockRepos := []*github.Repository{
				{
					ID:              1,
					Name:            "repo1",
					FullName:        "testuser/repo1",
					Description:     stringPtr("First repo"),
					Private:         false,
					HTMLURL:         "https://github.com/testuser/repo1",
					StargazersCount: 10,
					ForksCount:      5,
					DefaultBranch:   "main",
					Owner: github.User{
						ID:    12345,
						Login: "testuser",
					},
				},
				{
					ID:              2,
					Name:            "repo2",
					FullName:        "testuser/repo2",
					Description:     stringPtr("Second repo"),
					Private:         true,
					HTMLURL:         "https://github.com/testuser/repo2",
					StargazersCount: 20,
					ForksCount:      3,
					DefaultBranch:   "main",
					Owner: github.User{
						ID:    12345,
						Login: "testuser",
					},
				},
			}
			w.Header().Set("Link", "<http://test/users/testuser/repos?page=2>; rel=\"next\"")
			json.NewEncoder(w).Encode(mockRepos)
		}
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
				repositories(first: 2) {
					totalCount
					edges {
						cursor
						node {
							name
							description
							stargazerCount
						}
					}
					nodes {
						name
						fullName: nameWithOwner
					}
					pageInfo {
						hasNextPage
						hasPreviousPage
						startCursor
						endCursor
					}
				}
			}
		}
	`

	// Create context with GitHub client
	ctx := context.Background()
	ctx = context.WithValue(ctx, githubClientKey, client)

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
		Context:       ctx,
	})

	// Check for errors
	if len(result.Errors) > 0 {
		t.Fatalf("GraphQL errors: %v", result.Errors)
	}

	// Verify both endpoints were called
	if !userCalled {
		t.Error("Expected /user endpoint to be called")
	}
	if !reposCalled {
		t.Error("Expected /users/testuser/repos endpoint to be called")
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

	repos, ok := viewer["repositories"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected repositories to be a map")
	}

	// Verify totalCount
	if repos["totalCount"] != 2 {
		t.Errorf("Expected totalCount 2, got %v", repos["totalCount"])
	}

	// Verify edges
	edges, ok := repos["edges"].([]interface{})
	if !ok {
		t.Fatal("Expected edges to be an array")
	}
	if len(edges) != 2 {
		t.Fatalf("Expected 2 edges, got %d", len(edges))
	}

	firstEdge, ok := edges[0].(map[string]interface{})
	if !ok {
		t.Fatal("Expected first edge to be a map")
	}
	if firstEdge["cursor"] == nil {
		t.Error("Expected cursor to be present")
	}

	firstNode, ok := firstEdge["node"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected node to be a map")
	}
	if firstNode["name"] != "repo1" {
		t.Errorf("Expected first repo name 'repo1', got %v", firstNode["name"])
	}
	if firstNode["stargazerCount"] != 10 {
		t.Errorf("Expected stargazerCount 10, got %v", firstNode["stargazerCount"])
	}

	// Verify nodes
	nodes, ok := repos["nodes"].([]interface{})
	if !ok {
		t.Fatal("Expected nodes to be an array")
	}
	if len(nodes) != 2 {
		t.Fatalf("Expected 2 nodes, got %d", len(nodes))
	}

	firstNodeDirect, ok := nodes[0].(map[string]interface{})
	if !ok {
		t.Fatal("Expected first node to be a map")
	}
	if firstNodeDirect["name"] != "repo1" {
		t.Errorf("Expected first node name 'repo1', got %v", firstNodeDirect["name"])
	}
	if firstNodeDirect["fullName"] != "testuser/repo1" {
		t.Errorf("Expected fullName 'testuser/repo1', got %v", firstNodeDirect["fullName"])
	}

	// Verify pageInfo
	pageInfo, ok := repos["pageInfo"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected pageInfo to be a map")
	}
	if pageInfo["hasNextPage"] != true {
		t.Error("Expected hasNextPage to be true")
	}
	if pageInfo["hasPreviousPage"] != false {
		t.Error("Expected hasPreviousPage to be false")
	}
	if pageInfo["startCursor"] == nil {
		t.Error("Expected startCursor to be present")
	}
	if pageInfo["endCursor"] == nil {
		t.Error("Expected endCursor to be present")
	}
}

func TestUserRepositoriesQuery_Success(t *testing.T) {
	// Create mock GitHub API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.URL.Path == "/users/octocat" {
			mockUser := github.User{
				ID:    583231,
				Login: "octocat",
				Name:  stringPtr("The Octocat"),
			}
			json.NewEncoder(w).Encode(mockUser)
		} else if r.URL.Path == "/users/octocat/repos" {
			mockRepos := []*github.Repository{
				{
					ID:       1,
					Name:     "Hello-World",
					FullName: "octocat/Hello-World",
					Owner: github.User{
						ID:    583231,
						Login: "octocat",
					},
				},
			}
			json.NewEncoder(w).Encode(mockRepos)
		}
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
				repositories(first: 10) {
					totalCount
					nodes {
						name
					}
				}
			}
		}
	`

	// Create context with GitHub client
	ctx := context.Background()
	ctx = context.WithValue(ctx, githubClientKey, client)

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
		Context:       ctx,
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

	if user["login"] != "octocat" {
		t.Errorf("Expected login 'octocat', got %v", user["login"])
	}

	repos, ok := user["repositories"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected repositories to be a map")
	}

	nodes, ok := repos["nodes"].([]interface{})
	if !ok {
		t.Fatal("Expected nodes to be an array")
	}
	if len(nodes) != 1 {
		t.Fatalf("Expected 1 node, got %d", len(nodes))
	}

	firstNode, ok := nodes[0].(map[string]interface{})
	if !ok {
		t.Fatal("Expected first node to be a map")
	}
	if firstNode["name"] != "Hello-World" {
		t.Errorf("Expected repo name 'Hello-World', got %v", firstNode["name"])
	}
}

func TestRepositoriesQuery_WithPagination(t *testing.T) {
	// Create mock GitHub API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.URL.Path == "/user" {
			mockUser := github.User{
				ID:    12345,
				Login: "testuser",
			}
			json.NewEncoder(w).Encode(mockUser)
		} else if r.URL.Path == "/users/testuser/repos" {
			// Check if this is page 2
			page := r.URL.Query().Get("page")
			if page == "2" {
				// Return second page
				mockRepos := []*github.Repository{
					{
						ID:       3,
						Name:     "repo3",
						FullName: "testuser/repo3",
						Owner: github.User{
							ID:    12345,
							Login: "testuser",
						},
					},
				}
				// No next link - this is the last page
				json.NewEncoder(w).Encode(mockRepos)
			} else {
				// Return first page
				mockRepos := []*github.Repository{
					{
						ID:       1,
						Name:     "repo1",
						FullName: "testuser/repo1",
						Owner: github.User{
							ID:    12345,
							Login: "testuser",
						},
					},
				}
				w.Header().Set("Link", "<http://test/users/testuser/repos?page=2>; rel=\"next\"")
				json.NewEncoder(w).Encode(mockRepos)
			}
		}
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

	// Execute query for first page
	query := `
		query {
			viewer {
				repositories(first: 1) {
					pageInfo {
						hasNextPage
						endCursor
					}
					nodes {
						name
					}
				}
			}
		}
	`

	// Create context with GitHub client
	ctx := context.Background()
	ctx = context.WithValue(ctx, githubClientKey, client)

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
		Context:       ctx,
	})

	// Check for errors
	if len(result.Errors) > 0 {
		t.Fatalf("GraphQL errors: %v", result.Errors)
	}

	// Verify first page
	data := result.Data.(map[string]interface{})
	viewer := data["viewer"].(map[string]interface{})
	repos := viewer["repositories"].(map[string]interface{})
	pageInfo := repos["pageInfo"].(map[string]interface{})

	if pageInfo["hasNextPage"] != true {
		t.Error("Expected hasNextPage to be true on first page")
	}

	endCursor, ok := pageInfo["endCursor"].(string)
	if !ok || endCursor == "" {
		t.Fatal("Expected endCursor to be present")
	}

	// Now query second page with the cursor
	query2 := `
		query($cursor: String!) {
			viewer {
				repositories(first: 1, after: $cursor) {
					pageInfo {
						hasNextPage
						hasPreviousPage
					}
					nodes {
						name
					}
				}
			}
		}
	`

	result2 := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query2,
		VariableValues: map[string]interface{}{
			"cursor": endCursor,
		},
		Context: ctx,
	})

	// Check for errors
	if len(result2.Errors) > 0 {
		t.Fatalf("GraphQL errors on second page: %v", result2.Errors)
	}

	// Verify second page
	data2 := result2.Data.(map[string]interface{})
	viewer2 := data2["viewer"].(map[string]interface{})
	repos2 := viewer2["repositories"].(map[string]interface{})
	pageInfo2 := repos2["pageInfo"].(map[string]interface{})

	if pageInfo2["hasNextPage"] != false {
		t.Error("Expected hasNextPage to be false on last page")
	}
	if pageInfo2["hasPreviousPage"] != true {
		t.Error("Expected hasPreviousPage to be true on second page")
	}

	nodes2 := repos2["nodes"].([]interface{})
	if len(nodes2) != 1 {
		t.Fatalf("Expected 1 node on second page, got %d", len(nodes2))
	}
	node2 := nodes2[0].(map[string]interface{})
	if node2["name"] != "repo3" {
		t.Errorf("Expected second page to have 'repo3', got %v", node2["name"])
	}
}

func stringPtr(s string) *string {
	return &s
}
