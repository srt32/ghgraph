package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/graphql-go/graphql"
	githubclient "github.com/srt32/ghgraph/internal/github"
	ghgraphql "github.com/srt32/ghgraph/internal/graphql"
)

// Server represents the GraphQL HTTP server
type Server struct {
	port int
}

// NewServer creates a new server instance
func NewServer(port int) *Server {
	return &Server{
		port: port,
	}
}

// graphQLRequest represents an incoming GraphQL request
type graphQLRequest struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

// Start starts the HTTP server
func (s *Server) Start() error {
	http.HandleFunc("/graphql", s.handleGraphQL)
	http.HandleFunc("/health", s.handleHealth)

	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("Starting GraphQL server on %s", addr)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) handleGraphQL(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract GitHub token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		http.Error(w, "Invalid Authorization header format. Expected: Bearer <token>", http.StatusUnauthorized)
		return
	}

	// Parse GraphQL request
	var req graphQLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Create GitHub client with the provided token
	githubClient := githubclient.NewClient(token)

	// Create GraphQL schema
	schema, err := ghgraphql.NewSchema(githubClient)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create schema: %v", err), http.StatusInternalServerError)
		return
	}

	// Create context with GitHub client
	ctx := context.WithValue(context.Background(), "githubClient", githubClient)

	// Execute GraphQL query
	result := graphql.Do(graphql.Params{
		Schema:         schema,
		RequestString:  req.Query,
		VariableValues: req.Variables,
		OperationName:  req.OperationName,
		Context:        ctx,
	})

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Write response
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
