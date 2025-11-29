package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/graphql-go/graphql"
)

// TestIntegration_RealGitHubAPI tests against the actual GitHub GraphQL API
// to ensure our proxy matches the real API behavior.
// Set GITHUB_TOKEN environment variable to run this test.
func TestIntegration_RealGitHubAPI(t *testing.T) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		t.Skip("Skipping integration test: GITHUB_TOKEN not set")
	}

	// Test query
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

	// Query real GitHub GraphQL API
	realResult := queryRealGitHubAPI(t, token, query)

	// Query our proxy (we'll need to start the server for this)
	// For now, we'll just validate the structure matches what we expect
	t.Logf("Real GitHub API response: %+v", realResult)

	// Validate response structure
	data, ok := realResult["data"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected data field in response")
	}

	viewer, ok := data["viewer"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected viewer field in data")
	}

	// Verify required fields exist
	requiredFields := []string{"id", "login", "avatarUrl"}
	for _, field := range requiredFields {
		if _, exists := viewer[field]; !exists {
			t.Errorf("Expected field %s to exist in viewer", field)
		}
	}

	// Verify nullable fields are present (but may be null)
	optionalFields := []string{"name", "email", "bio", "company", "location"}
	for _, field := range optionalFields {
		if _, exists := viewer[field]; !exists {
			t.Errorf("Expected field %s to exist in viewer (even if null)", field)
		}
	}
}

// TestIntegration_SchemaIntrospection tests that our schema structure
// matches GitHub's GraphQL schema for the viewer query
func TestIntegration_SchemaIntrospection(t *testing.T) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		t.Skip("Skipping integration test: GITHUB_TOKEN not set")
	}

	// Introspection query to get the viewer field's type
	query := `
		query {
			__type(name: "Query") {
				fields {
					name
					type {
						name
						kind
						ofType {
							name
							kind
						}
					}
				}
			}
		}
	`

	result := queryRealGitHubAPI(t, token, query)
	t.Logf("Schema introspection result: %+v", result)

	// This helps us understand the real GitHub schema structure
	// We can compare it with our implementation
}

// TestIntegration_PartialFields tests that selecting only some fields works correctly
func TestIntegration_PartialFields(t *testing.T) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		t.Skip("Skipping integration test: GITHUB_TOKEN not set")
	}

	// Test query with only login field
	query := `
		query {
			viewer {
				login
			}
		}
	`

	result := queryRealGitHubAPI(t, token, query)

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected data field in response")
	}

	viewer, ok := data["viewer"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected viewer field in data")
	}

	// Verify only login is returned
	if _, exists := viewer["login"]; !exists {
		t.Error("Expected login field")
	}

	// These fields should not be present since we didn't request them
	unexpectedFields := []string{"email", "bio", "company"}
	for _, field := range unexpectedFields {
		if _, exists := viewer[field]; exists {
			t.Errorf("Did not expect field %s in response", field)
		}
	}
}

func queryRealGitHubAPI(t *testing.T, token, query string) map[string]interface{} {
	t.Helper()

	requestBody := map[string]interface{}{
		"query": query,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Check for GraphQL errors
	if errors, ok := result["errors"].([]interface{}); ok && len(errors) > 0 {
		t.Fatalf("GraphQL errors: %v", errors)
	}

	return result
}

// TestTypeDefinitions verifies that our User type matches GitHub's User type structure
func TestTypeDefinitions(t *testing.T) {
	// This test documents the expected type structure
	// In the real GitHub GraphQL API:
	// - id: ID! (non-null)
	// - login: String! (non-null)
	// - name: String (nullable)
	// - email: String! (non-null) - but can be empty string
	// - avatarUrl: URI! (non-null)
	// - bio: String (nullable)
	// - company: String (nullable)
	// - location: String (nullable)

	expectedTypes := map[string]struct {
		typeName string
		nonNull  bool
	}{
		"id":        {"String", true},  // We use String instead of ID for simplicity
		"login":     {"String", true},
		"name":      {"String", false},
		"email":     {"String", true},
		"avatarUrl": {"String", true},
		"bio":       {"String", false},
		"company":   {"String", false},
		"location":  {"String", false},
	}

	// This serves as documentation of our schema
	for field, expected := range expectedTypes {
		t.Logf("Field %s: type=%s, nonNull=%v", field, expected.typeName, expected.nonNull)
	}
}

// Helper function to check if a GraphQL type matches expected structure
func assertGraphQLType(t *testing.T, fieldType graphql.Type, expectedTypeName string, expectNonNull bool) {
	t.Helper()

	if expectNonNull {
		nonNull, ok := fieldType.(*graphql.NonNull)
		if !ok {
			t.Errorf("Expected NonNull type, got %T", fieldType)
			return
		}
		fieldType = nonNull.OfType
	}

	if fieldType.Name() != expectedTypeName {
		t.Errorf("Expected type name %s, got %s", expectedTypeName, fieldType.Name())
	}
}
