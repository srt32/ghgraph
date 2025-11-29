package github

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/srt32/ghgraph/internal/cache"
)

const (
	defaultBaseURL = "https://api.github.com"
	apiVersion     = "2022-11-28"
	cacheTTL       = 5 * time.Minute
)

var (
	// sharedCache is a global cache instance shared across all clients
	sharedCache = cache.NewCache(cacheTTL)
)

// Client is a GitHub REST API client with caching support
type Client struct {
	baseURL    string
	httpClient *http.Client
	token      string
	cache      *cache.Cache
}

// User represents a GitHub user from the REST API
type User struct {
	ID        int64   `json:"id"`
	Login     string  `json:"login"`
	Name      *string `json:"name"`
	Email     *string `json:"email"`
	AvatarURL string  `json:"avatar_url"`
	Bio       *string `json:"bio"`
	Company   *string `json:"company"`
	Location  *string `json:"location"`
}

// Repository represents a GitHub repository from the REST API
type Repository struct {
	ID              int64   `json:"id"`
	Name            string  `json:"name"`
	FullName        string  `json:"full_name"`
	Description     *string `json:"description"`
	Private         bool    `json:"private"`
	HTMLURL         string  `json:"html_url"`
	StargazersCount int     `json:"stargazers_count"`
	ForksCount      int     `json:"forks_count"`
	DefaultBranch   string  `json:"default_branch"`
	Owner           User    `json:"owner"`
}

// Organization represents a GitHub organization from the REST API
type Organization struct {
	ID          int64   `json:"id"`
	Login       string  `json:"login"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	AvatarURL   string  `json:"avatar_url"`
	HTMLURL     string  `json:"html_url"`
	Email       *string `json:"email"`
	Location    *string `json:"location"`
}

// RepositoriesResult represents a paginated list of repositories
type RepositoriesResult struct {
	Repositories []*Repository
	TotalCount   int
	HasNextPage  bool
	HasPrevPage  bool
	EndCursor    *string
	StartCursor  *string
}

// Issue represents a GitHub issue from the REST API
type Issue struct {
	ID        int64   `json:"id"`
	Number    int     `json:"number"`
	Title     string  `json:"title"`
	Body      *string `json:"body"`
	State     string  `json:"state"`
	HTMLURL   string  `json:"html_url"`
	User      User    `json:"user"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	ClosedAt  *string `json:"closed_at"`
}

// PullRequest represents a GitHub pull request from the REST API
type PullRequest struct {
	ID        int64   `json:"id"`
	Number    int     `json:"number"`
	Title     string  `json:"title"`
	Body      *string `json:"body"`
	State     string  `json:"state"`
	HTMLURL   string  `json:"html_url"`
	User      User    `json:"user"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	ClosedAt  *string `json:"closed_at"`
	MergedAt  *string `json:"merged_at"`
	Merged    bool    `json:"merged"`
}

// IssuesResult represents a paginated list of issues
type IssuesResult struct {
	Issues      []*Issue
	TotalCount  int
	HasNextPage bool
	HasPrevPage bool
	EndCursor   *string
	StartCursor *string
}

// PullRequestsResult represents a paginated list of pull requests
type PullRequestsResult struct {
	PullRequests []*PullRequest
	TotalCount   int
	HasNextPage  bool
	HasPrevPage  bool
	EndCursor    *string
	StartCursor  *string
}

// NewClient creates a new GitHub REST API client
func NewClient(token string) *Client {
	return &Client{
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		token: token,
		cache: sharedCache,
	}
}

// SetBaseURL sets the base URL for the GitHub API (useful for testing)
func (c *Client) SetBaseURL(url string) {
	c.baseURL = url
}

// GetAuthenticatedUser fetches the authenticated user's information
func (c *Client) GetAuthenticatedUser() (*User, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/user", c.baseURL), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", apiVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &user, nil
}

// GetUser fetches a specific user by their login (username)
func (c *Client) GetUser(login string) (*User, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/users/%s", c.baseURL, login), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", apiVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("user not found: %s", login)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &user, nil
}

// GetRepository fetches a specific repository by owner and name
func (c *Client) GetRepository(owner, name string) (*Repository, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/repos/%s/%s", c.baseURL, owner, name), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", apiVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("repository not found: %s/%s", owner, name)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var repo Repository
	if err := json.NewDecoder(resp.Body).Decode(&repo); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &repo, nil
}

// GetOrganization fetches a specific organization by login
func (c *Client) GetOrganization(login string) (*Organization, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/orgs/%s", c.baseURL, login), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", apiVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("organization not found: %s", login)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var org Organization
	if err := json.NewDecoder(resp.Body).Decode(&org); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &org, nil
}

// ListUserRepositories fetches repositories for a user with pagination support
// If login is empty, fetches repositories for the authenticated user
// Supports both forward (first/after) and backward (last/before) pagination
func (c *Client) ListUserRepositories(login string, first *int, after *string, last *int, before *string) (*RepositoriesResult, error) {
	// Determine page size and number
	perPage := 30 // default
	page := 1

	// Forward pagination: first/after
	if first != nil {
		perPage = *first
		if after != nil && *after != "" {
			decoded, err := base64.StdEncoding.DecodeString(*after)
			if err == nil {
				if p, err := strconv.Atoi(string(decoded)); err == nil && p > 0 {
					page = p
				}
			}
		}
	}

	// Backward pagination: last/before
	if last != nil {
		perPage = *last
		if before != nil && *before != "" {
			decoded, err := base64.StdEncoding.DecodeString(*before)
			if err == nil {
				if p, err := strconv.Atoi(string(decoded)); err == nil && p > 1 {
					page = p - 1 // Go to previous page
				}
			}
		}
	}

	// Determine endpoint
	var endpoint string
	if login == "" {
		endpoint = fmt.Sprintf("%s/user/repos", c.baseURL)
	} else {
		endpoint = fmt.Sprintf("%s/users/%s/repos", c.baseURL, login)
	}

	// Build query parameters
	params := url.Values{}
	params.Set("per_page", strconv.Itoa(perPage))
	params.Set("page", strconv.Itoa(page))
	params.Set("sort", "updated")
	params.Set("direction", "desc")

	fullURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", apiVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		if login == "" {
			return nil, fmt.Errorf("authenticated user repositories not found")
		}
		return nil, fmt.Errorf("user not found: %s", login)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var repos []*Repository
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	// Parse Link header to determine pagination info
	hasNextPage := false
	hasPrevPage := page > 1
	linkHeader := resp.Header.Get("Link")
	if linkHeader != "" {
		links := parseLinkHeader(linkHeader)
		if _, ok := links["next"]; ok {
			hasNextPage = true
		}
	}

	// Create cursors
	var startCursor, endCursor *string
	if len(repos) > 0 {
		start := base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(page)))
		startCursor = &start
		if hasNextPage {
			end := base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(page + 1)))
			endCursor = &end
		}
	}

	return &RepositoriesResult{
		Repositories: repos,
		// NOTE: totalCount returns the count of items in the current page, not the total across all pages.
		// GitHub's REST API doesn't provide a total count in the response, and we would need to
		// parse the "last" rel from the Link header and multiply by per_page to estimate it,
		// which would still be inaccurate. This matches the behavior of returning page-level counts.
		TotalCount:   len(repos),
		HasNextPage:  hasNextPage,
		HasPrevPage:  hasPrevPage,
		EndCursor:    endCursor,
		StartCursor:  startCursor,
	}, nil
}

// ListRepositoryIssues fetches issues for a repository with pagination support
// Supports both forward (first/after) and backward (last/before) pagination
// state can be "open", "closed", or "all"
func (c *Client) ListRepositoryIssues(owner, name string, state string, first *int, after *string, last *int, before *string) (*IssuesResult, error) {
	// Determine page size and number
	perPage := 30
	page := 1

	if first != nil {
		perPage = *first
		if after != nil && *after != "" {
			if decoded, err := base64.StdEncoding.DecodeString(*after); err == nil {
				if p, err := strconv.Atoi(string(decoded)); err == nil && p > 0 {
					page = p
				}
			}
		}
	}

	if last != nil {
		perPage = *last
		if before != nil && *before != "" {
			if decoded, err := base64.StdEncoding.DecodeString(*before); err == nil {
				if p, err := strconv.Atoi(string(decoded)); err == nil && p > 1 {
					page = p - 1
				}
			}
		}
	}

	// Check cache first
	cacheKey := fmt.Sprintf("issues:%s/%s:%s:page%d:per%d", owner, name, state, page, perPage)
	if cached := c.cache.Get(c.token, cacheKey); cached != nil {
		return cached.(*IssuesResult), nil
	}

	endpoint := fmt.Sprintf("%s/repos/%s/%s/issues", c.baseURL, owner, name)
	params := url.Values{}
	params.Set("state", state)
	params.Set("per_page", strconv.Itoa(perPage))
	params.Set("page", strconv.Itoa(page))
	params.Set("sort", "updated")
	params.Set("direction", "desc")

	fullURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", apiVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("repository not found: %s/%s", owner, name)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var issues []*Issue
	if err := json.NewDecoder(resp.Body).Decode(&issues); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	hasNextPage := false
	hasPrevPage := page > 1
	linkHeader := resp.Header.Get("Link")
	if linkHeader != "" {
		links := parseLinkHeader(linkHeader)
		if _, ok := links["next"]; ok {
			hasNextPage = true
		}
	}

	var startCursor, endCursor *string
	if len(issues) > 0 {
		start := base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(page)))
		startCursor = &start
		if hasNextPage {
			end := base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(page + 1)))
			endCursor = &end
		}
	}

	result := &IssuesResult{
		Issues:      issues,
		TotalCount:  len(issues),
		HasNextPage: hasNextPage,
		HasPrevPage: hasPrevPage,
		EndCursor:   endCursor,
		StartCursor: startCursor,
	}

	// Cache the result
	c.cache.Set(c.token, cacheKey, result)

	return result, nil
}

// ListRepositoryPullRequests fetches pull requests for a repository with pagination support
// Supports both forward (first/after) and backward (last/before) pagination
// state can be "open", "closed", or "all"
func (c *Client) ListRepositoryPullRequests(owner, name string, state string, first *int, after *string, last *int, before *string) (*PullRequestsResult, error) {
	// Determine page size and number
	perPage := 30
	page := 1

	if first != nil {
		perPage = *first
		if after != nil && *after != "" {
			if decoded, err := base64.StdEncoding.DecodeString(*after); err == nil {
				if p, err := strconv.Atoi(string(decoded)); err == nil && p > 0 {
					page = p
				}
			}
		}
	}

	if last != nil {
		perPage = *last
		if before != nil && *before != "" {
			if decoded, err := base64.StdEncoding.DecodeString(*before); err == nil {
				if p, err := strconv.Atoi(string(decoded)); err == nil && p > 1 {
					page = p - 1
				}
			}
		}
	}

	// Check cache first
	cacheKey := fmt.Sprintf("pullRequests:%s/%s:%s:page%d:per%d", owner, name, state, page, perPage)
	if cached := c.cache.Get(c.token, cacheKey); cached != nil {
		return cached.(*PullRequestsResult), nil
	}

	endpoint := fmt.Sprintf("%s/repos/%s/%s/pulls", c.baseURL, owner, name)
	params := url.Values{}
	params.Set("state", state)
	params.Set("per_page", strconv.Itoa(perPage))
	params.Set("page", strconv.Itoa(page))
	params.Set("sort", "updated")
	params.Set("direction", "desc")

	fullURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", apiVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("repository not found: %s/%s", owner, name)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var prs []*PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&prs); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	hasNextPage := false
	hasPrevPage := page > 1
	linkHeader := resp.Header.Get("Link")
	if linkHeader != "" {
		links := parseLinkHeader(linkHeader)
		if _, ok := links["next"]; ok {
			hasNextPage = true
		}
	}

	var startCursor, endCursor *string
	if len(prs) > 0 {
		start := base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(page)))
		startCursor = &start
		if hasNextPage {
			end := base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(page + 1)))
			endCursor = &end
		}
	}

	result := &PullRequestsResult{
		PullRequests: prs,
		TotalCount:   len(prs),
		HasNextPage:  hasNextPage,
		HasPrevPage:  hasPrevPage,
		EndCursor:    endCursor,
		StartCursor:  startCursor,
	}

	// Cache the result
	c.cache.Set(c.token, cacheKey, result)

	return result, nil
}

// parseLinkHeader parses the GitHub Link header format
// Example: <https://api.github.com/user/repos?page=2>; rel="next", <https://api.github.com/user/repos?page=5>; rel="last"
func parseLinkHeader(header string) map[string]string {
	links := make(map[string]string)
	parts := strings.Split(header, ",")
	for _, part := range parts {
		section := strings.Split(strings.TrimSpace(part), ";")
		if len(section) != 2 {
			continue
		}
		url := strings.Trim(strings.TrimSpace(section[0]), "<>")
		rel := strings.Trim(strings.TrimSpace(section[1]), " ")
		if strings.HasPrefix(rel, "rel=\"") {
			rel = strings.Trim(rel[5:], "\"")
			links[rel] = url
		}
	}
	return links
}
