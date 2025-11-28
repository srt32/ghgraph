package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultBaseURL = "https://api.github.com"
	apiVersion     = "2022-11-28"
)

// Client is a GitHub REST API client
type Client struct {
	baseURL    string
	httpClient *http.Client
	token      string
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

// NewClient creates a new GitHub REST API client
func NewClient(token string) *Client {
	return &Client{
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		token: token,
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
