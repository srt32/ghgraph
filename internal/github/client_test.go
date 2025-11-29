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

func TestGetRepository_Success(t *testing.T) {
	// Create mock repository data
	mockRepo := Repository{
		ID:              12345,
		Name:            "hello-world",
		FullName:        "octocat/hello-world",
		Description:     stringPtr("My first repository"),
		Private:         false,
		HTMLURL:         "https://github.com/octocat/hello-world",
		StargazersCount: 100,
		ForksCount:      25,
		DefaultBranch:   "main",
		Owner: User{
			ID:    583231,
			Login: "octocat",
		},
	}

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request headers
		if auth := r.Header.Get("Authorization"); auth != "Bearer test-token" {
			t.Errorf("Expected Authorization header 'Bearer test-token', got '%s'", auth)
		}

		// Verify endpoint
		expectedPath := "/repos/octocat/hello-world"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path '%s', got '%s'", expectedPath, r.URL.Path)
		}

		// Return mock repository
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockRepo)
	}))
	defer server.Close()

	// Create client pointing to test server
	client := NewClient("test-token")
	client.baseURL = server.URL

	// Call GetRepository
	repo, err := client.GetRepository("octocat", "hello-world")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify response
	if repo.Name != mockRepo.Name {
		t.Errorf("Expected Name %s, got %s", mockRepo.Name, repo.Name)
	}
	if repo.FullName != mockRepo.FullName {
		t.Errorf("Expected FullName %s, got %s", mockRepo.FullName, repo.FullName)
	}
	if repo.StargazersCount != mockRepo.StargazersCount {
		t.Errorf("Expected StargazersCount %d, got %d", mockRepo.StargazersCount, repo.StargazersCount)
	}
	if repo.Owner.Login != mockRepo.Owner.Login {
		t.Errorf("Expected Owner.Login %s, got %s", mockRepo.Owner.Login, repo.Owner.Login)
	}
}

func TestGetRepository_NotFound(t *testing.T) {
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

	// Call GetRepository
	_, err := client.GetRepository("nonexistent", "repo")
	if err == nil {
		t.Fatal("Expected error for non-existent repository, got nil")
	}

	// Verify error message
	expectedMsg := "repository not found: nonexistent/repo"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestGetOrganization_Success(t *testing.T) {
	// Create mock organization data
	mockOrg := Organization{
		ID:          1,
		Login:       "github",
		Name:        stringPtr("GitHub"),
		Description: stringPtr("How people build software"),
		AvatarURL:   "https://avatars.githubusercontent.com/u/9919",
		HTMLURL:     "https://github.com/github",
		Email:       stringPtr("support@github.com"),
		Location:    stringPtr("San Francisco, CA"),
	}

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request headers
		if auth := r.Header.Get("Authorization"); auth != "Bearer test-token" {
			t.Errorf("Expected Authorization header 'Bearer test-token', got '%s'", auth)
		}

		// Verify endpoint
		expectedPath := "/orgs/github"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path '%s', got '%s'", expectedPath, r.URL.Path)
		}

		// Return mock organization
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockOrg)
	}))
	defer server.Close()

	// Create client pointing to test server
	client := NewClient("test-token")
	client.baseURL = server.URL

	// Call GetOrganization
	org, err := client.GetOrganization("github")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify response
	if org.Login != mockOrg.Login {
		t.Errorf("Expected Login %s, got %s", mockOrg.Login, org.Login)
	}
	if *org.Name != *mockOrg.Name {
		t.Errorf("Expected Name %s, got %s", *mockOrg.Name, *org.Name)
	}
	if *org.Description != *mockOrg.Description {
		t.Errorf("Expected Description %s, got %s", *mockOrg.Description, *org.Description)
	}
}

func TestGetOrganization_NotFound(t *testing.T) {
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

	// Call GetOrganization
	_, err := client.GetOrganization("nonexistent")
	if err == nil {
		t.Fatal("Expected error for non-existent organization, got nil")
	}

	// Verify error message
	expectedMsg := "organization not found: nonexistent"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestListUserRepositories_Success(t *testing.T) {
	// Create mock repositories data
	mockRepos := []*Repository{
		{
			ID:              1,
			Name:            "repo1",
			FullName:        "octocat/repo1",
			Description:     stringPtr("First repo"),
			Private:         false,
			HTMLURL:         "https://github.com/octocat/repo1",
			StargazersCount: 10,
			ForksCount:      5,
			DefaultBranch:   "main",
			Owner: User{
				ID:    583231,
				Login: "octocat",
			},
		},
		{
			ID:              2,
			Name:            "repo2",
			FullName:        "octocat/repo2",
			Description:     stringPtr("Second repo"),
			Private:         true,
			HTMLURL:         "https://github.com/octocat/repo2",
			StargazersCount: 20,
			ForksCount:      3,
			DefaultBranch:   "main",
			Owner: User{
				ID:    583231,
				Login: "octocat",
			},
		},
	}

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request headers
		if auth := r.Header.Get("Authorization"); auth != "Bearer test-token" {
			t.Errorf("Expected Authorization header 'Bearer test-token', got '%s'", auth)
		}

		// Verify endpoint
		expectedPath := "/users/octocat/repos"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path '%s', got '%s'", expectedPath, r.URL.Path)
		}

		// Verify query parameters
		if r.URL.Query().Get("per_page") != "2" {
			t.Errorf("Expected per_page=2, got %s", r.URL.Query().Get("per_page"))
		}
		if r.URL.Query().Get("page") != "1" {
			t.Errorf("Expected page=1, got %s", r.URL.Query().Get("page"))
		}

		// Return mock repositories with pagination link
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Link", "<http://test/users/octocat/repos?page=2>; rel=\"next\", <http://test/users/octocat/repos?page=5>; rel=\"last\"")
		json.NewEncoder(w).Encode(mockRepos)
	}))
	defer server.Close()

	// Create client pointing to test server
	client := NewClient("test-token")
	client.baseURL = server.URL

	// Call ListUserRepositories
	first := 2
	result, err := client.ListUserRepositories("octocat", &first, nil, nil, nil)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify response
	if len(result.Repositories) != 2 {
		t.Errorf("Expected 2 repositories, got %d", len(result.Repositories))
	}
	if result.HasNextPage != true {
		t.Error("Expected HasNextPage to be true")
	}
	if result.HasPrevPage != false {
		t.Error("Expected HasPrevPage to be false")
	}
	if result.EndCursor == nil {
		t.Error("Expected EndCursor to be non-nil")
	}
	if result.Repositories[0].Name != "repo1" {
		t.Errorf("Expected first repo name 'repo1', got '%s'", result.Repositories[0].Name)
	}
}

func TestListUserRepositories_WithCursor(t *testing.T) {
	// Create mock repositories data
	mockRepos := []*Repository{
		{
			ID:       3,
			Name:     "repo3",
			FullName: "octocat/repo3",
			Owner: User{
				ID:    583231,
				Login: "octocat",
			},
		},
	}

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify page parameter is 2 (decoded from cursor)
		if r.URL.Query().Get("page") != "2" {
			t.Errorf("Expected page=2, got %s", r.URL.Query().Get("page"))
		}

		// Return mock repositories without next link
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockRepos)
	}))
	defer server.Close()

	// Create client pointing to test server
	client := NewClient("test-token")
	client.baseURL = server.URL

	// Create cursor for page 2
	cursor := "Mg==" // base64 encoded "2"

	// Call ListUserRepositories with cursor
	first := 1
	result, err := client.ListUserRepositories("octocat", &first, &cursor, nil, nil)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify response
	if len(result.Repositories) != 1 {
		t.Errorf("Expected 1 repository, got %d", len(result.Repositories))
	}
	if result.HasNextPage != false {
		t.Error("Expected HasNextPage to be false")
	}
	if result.HasPrevPage != true {
		t.Error("Expected HasPrevPage to be true")
	}
}

func TestListUserRepositories_AuthenticatedUser(t *testing.T) {
	// Create mock repositories data
	mockRepos := []*Repository{
		{
			ID:       1,
			Name:     "my-repo",
			FullName: "myuser/my-repo",
			Owner: User{
				ID:    12345,
				Login: "myuser",
			},
		},
	}

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify endpoint is /user/repos (authenticated user)
		expectedPath := "/user/repos"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path '%s', got '%s'", expectedPath, r.URL.Path)
		}

		// Return mock repositories
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockRepos)
	}))
	defer server.Close()

	// Create client pointing to test server
	client := NewClient("test-token")
	client.baseURL = server.URL

	// Call ListUserRepositories with empty login (authenticated user)
	first := 30
	result, err := client.ListUserRepositories("", &first, nil, nil, nil)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify response
	if len(result.Repositories) != 1 {
		t.Errorf("Expected 1 repository, got %d", len(result.Repositories))
	}
	if result.Repositories[0].Name != "my-repo" {
		t.Errorf("Expected repo name 'my-repo', got '%s'", result.Repositories[0].Name)
	}
}

func TestListUserRepositories_NotFound(t *testing.T) {
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

	// Call ListUserRepositories
	first := 30
	_, err := client.ListUserRepositories("nonexistent", &first, nil, nil, nil)
	if err == nil {
		t.Fatal("Expected error for non-existent user, got nil")
	}

	// Verify error message
	expectedMsg := "user not found: nonexistent"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestListRepositoryIssues_Success(t *testing.T) {
	// Create mock issues data
	mockIssues := []*Issue{
		{
			ID:        1,
			Number:    10,
			Title:     "First issue",
			Body:      stringPtr("Issue body"),
			State:     "open",
			HTMLURL:   "https://github.com/octocat/repo/issues/10",
			User:      User{ID: 1, Login: "user1"},
			CreatedAt: "2023-01-01T00:00:00Z",
			UpdatedAt: "2023-01-02T00:00:00Z",
			ClosedAt:  nil,
		},
		{
			ID:        2,
			Number:    11,
			Title:     "Second issue",
			Body:      stringPtr("Another issue"),
			State:     "closed",
			HTMLURL:   "https://github.com/octocat/repo/issues/11",
			User:      User{ID: 2, Login: "user2"},
			CreatedAt: "2023-01-03T00:00:00Z",
			UpdatedAt: "2023-01-04T00:00:00Z",
			ClosedAt:  stringPtr("2023-01-05T00:00:00Z"),
		},
	}

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request headers
		if auth := r.Header.Get("Authorization"); auth != "Bearer test-token" {
			t.Errorf("Expected Authorization header 'Bearer test-token', got '%s'", auth)
		}

		// Verify endpoint
		expectedPath := "/repos/octocat/repo/issues"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path '%s', got '%s'", expectedPath, r.URL.Path)
		}

		// Verify query parameters
		if r.URL.Query().Get("state") != "all" {
			t.Errorf("Expected state=all, got %s", r.URL.Query().Get("state"))
		}
		if r.URL.Query().Get("per_page") != "2" {
			t.Errorf("Expected per_page=2, got %s", r.URL.Query().Get("per_page"))
		}

		// Return mock issues with pagination link
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Link", "<http://test/repos/octocat/repo/issues?page=2>; rel=\"next\", <http://test/repos/octocat/repo/issues?page=5>; rel=\"last\"")
		json.NewEncoder(w).Encode(mockIssues)
	}))
	defer server.Close()

	// Create client pointing to test server
	client := NewClient("test-token")
	client.baseURL = server.URL

	// Call ListRepositoryIssues
	first := 2
	result, err := client.ListRepositoryIssues("octocat", "repo", "all", &first, nil, nil, nil)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify response
	if len(result.Issues) != 2 {
		t.Errorf("Expected 2 issues, got %d", len(result.Issues))
	}
	if result.HasNextPage != true {
		t.Error("Expected HasNextPage to be true")
	}
	if result.HasPrevPage != false {
		t.Error("Expected HasPrevPage to be false")
	}
	if result.Issues[0].Title != "First issue" {
		t.Errorf("Expected first issue title 'First issue', got '%s'", result.Issues[0].Title)
	}
	if result.Issues[0].State != "open" {
		t.Errorf("Expected first issue state 'open', got '%s'", result.Issues[0].State)
	}
	if result.Issues[1].ClosedAt == nil {
		t.Error("Expected second issue to have ClosedAt")
	}
}

func TestListRepositoryIssues_WithCache(t *testing.T) {
	callCount := 0

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]*Issue{
			{
				ID:        1,
				Number:    10,
				Title:     "Test issue",
				State:     "open",
				HTMLURL:   "https://github.com/octocat/repo/issues/10",
				User:      User{ID: 1, Login: "user1"},
				CreatedAt: "2023-01-01T00:00:00Z",
				UpdatedAt: "2023-01-02T00:00:00Z",
			},
		})
	}))
	defer server.Close()

	// Create client pointing to test server
	client := NewClient("test-token")
	client.baseURL = server.URL

	// First call - should hit the API
	first := 1
	result1, err := client.ListRepositoryIssues("octocat", "repo", "all", &first, nil, nil, nil)
	if err != nil {
		t.Fatalf("Expected no error on first call, got: %v", err)
	}

	// Second call with same parameters - should hit cache
	result2, err := client.ListRepositoryIssues("octocat", "repo", "all", &first, nil, nil, nil)
	if err != nil {
		t.Fatalf("Expected no error on second call, got: %v", err)
	}

	// Verify cache was used (API only called once)
	if callCount != 1 {
		t.Errorf("Expected API to be called once, was called %d times", callCount)
	}

	// Verify both results are identical
	if len(result1.Issues) != len(result2.Issues) {
		t.Error("Expected cached result to match original")
	}
	if result1.Issues[0].Title != result2.Issues[0].Title {
		t.Error("Expected cached result to match original")
	}
}

func TestListRepositoryIssues_NotFound(t *testing.T) {
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

	// Call ListRepositoryIssues
	first := 30
	_, err := client.ListRepositoryIssues("nonexistent", "repo", "all", &first, nil, nil, nil)
	if err == nil {
		t.Fatal("Expected error for non-existent repository, got nil")
	}

	// Verify error message
	expectedMsg := "repository not found: nonexistent/repo"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestListRepositoryPullRequests_Success(t *testing.T) {
	// Create mock pull requests data
	mockPRs := []*PullRequest{
		{
			ID:        1,
			Number:    20,
			Title:     "First PR",
			Body:      stringPtr("PR body"),
			State:     "open",
			HTMLURL:   "https://github.com/octocat/repo/pull/20",
			User:      User{ID: 1, Login: "user1"},
			CreatedAt: "2023-01-01T00:00:00Z",
			UpdatedAt: "2023-01-02T00:00:00Z",
			ClosedAt:  nil,
			MergedAt:  nil,
			Merged:    false,
		},
		{
			ID:        2,
			Number:    21,
			Title:     "Second PR",
			Body:      stringPtr("Another PR"),
			State:     "closed",
			HTMLURL:   "https://github.com/octocat/repo/pull/21",
			User:      User{ID: 2, Login: "user2"},
			CreatedAt: "2023-01-03T00:00:00Z",
			UpdatedAt: "2023-01-04T00:00:00Z",
			ClosedAt:  stringPtr("2023-01-05T00:00:00Z"),
			MergedAt:  stringPtr("2023-01-05T00:00:00Z"),
			Merged:    true,
		},
	}

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request headers
		if auth := r.Header.Get("Authorization"); auth != "Bearer test-token" {
			t.Errorf("Expected Authorization header 'Bearer test-token', got '%s'", auth)
		}

		// Verify endpoint
		expectedPath := "/repos/octocat/repo/pulls"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path '%s', got '%s'", expectedPath, r.URL.Path)
		}

		// Verify query parameters
		if r.URL.Query().Get("state") != "all" {
			t.Errorf("Expected state=all, got %s", r.URL.Query().Get("state"))
		}
		if r.URL.Query().Get("per_page") != "2" {
			t.Errorf("Expected per_page=2, got %s", r.URL.Query().Get("per_page"))
		}

		// Return mock PRs with pagination link
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Link", "<http://test/repos/octocat/repo/pulls?page=2>; rel=\"next\", <http://test/repos/octocat/repo/pulls?page=5>; rel=\"last\"")
		json.NewEncoder(w).Encode(mockPRs)
	}))
	defer server.Close()

	// Create client pointing to test server
	client := NewClient("test-token")
	client.baseURL = server.URL

	// Call ListRepositoryPullRequests
	first := 2
	result, err := client.ListRepositoryPullRequests("octocat", "repo", "all", &first, nil, nil, nil)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify response
	if len(result.PullRequests) != 2 {
		t.Errorf("Expected 2 pull requests, got %d", len(result.PullRequests))
	}
	if result.HasNextPage != true {
		t.Error("Expected HasNextPage to be true")
	}
	if result.HasPrevPage != false {
		t.Error("Expected HasPrevPage to be false")
	}
	if result.PullRequests[0].Title != "First PR" {
		t.Errorf("Expected first PR title 'First PR', got '%s'", result.PullRequests[0].Title)
	}
	if result.PullRequests[0].Merged != false {
		t.Error("Expected first PR to not be merged")
	}
	if result.PullRequests[1].Merged != true {
		t.Error("Expected second PR to be merged")
	}
	if result.PullRequests[1].MergedAt == nil {
		t.Error("Expected second PR to have MergedAt")
	}
}

func TestListRepositoryPullRequests_WithCache(t *testing.T) {
	callCount := 0

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]*PullRequest{
			{
				ID:        1,
				Number:    20,
				Title:     "Test PR",
				State:     "open",
				HTMLURL:   "https://github.com/octocat/repo/pull/20",
				User:      User{ID: 1, Login: "user1"},
				CreatedAt: "2023-01-01T00:00:00Z",
				UpdatedAt: "2023-01-02T00:00:00Z",
				Merged:    false,
			},
		})
	}))
	defer server.Close()

	// Create client pointing to test server
	client := NewClient("test-token")
	client.baseURL = server.URL

	// First call - should hit the API
	first := 1
	result1, err := client.ListRepositoryPullRequests("octocat", "repo", "all", &first, nil, nil, nil)
	if err != nil {
		t.Fatalf("Expected no error on first call, got: %v", err)
	}

	// Second call with same parameters - should hit cache
	result2, err := client.ListRepositoryPullRequests("octocat", "repo", "all", &first, nil, nil, nil)
	if err != nil {
		t.Fatalf("Expected no error on second call, got: %v", err)
	}

	// Verify cache was used (API only called once)
	if callCount != 1 {
		t.Errorf("Expected API to be called once, was called %d times", callCount)
	}

	// Verify both results are identical
	if len(result1.PullRequests) != len(result2.PullRequests) {
		t.Error("Expected cached result to match original")
	}
	if result1.PullRequests[0].Title != result2.PullRequests[0].Title {
		t.Error("Expected cached result to match original")
	}
}

func TestListRepositoryPullRequests_NotFound(t *testing.T) {
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

	// Call ListRepositoryPullRequests
	first := 30
	_, err := client.ListRepositoryPullRequests("nonexistent", "repo", "all", &first, nil, nil, nil)
	if err == nil {
		t.Fatal("Expected error for non-existent repository, got nil")
	}

	// Verify error message
	expectedMsg := "repository not found: nonexistent/repo"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func stringPtr(s string) *string {
	return &s
}
