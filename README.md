# ghgraph

A stateless GraphQL proxy that translates GitHub GraphQL API queries into GitHub REST API calls.

## Overview

This project clones the GitHub GraphQL API interface but translates all queries to REST API calls behind the scenes. It's designed to work as if the GitHub GraphQL API doesn't exist and only the REST API is available.

## Features

- **GraphQL Interface**: Exposes a GraphQL API that matches GitHub's schema
- **REST Backend**: Translates all queries to GitHub REST API v3 calls
- **Stateless**: No database or persistent state required
- **Viewer Query**: Currently supports the `viewer` query with scalar fields
- **Type-Safe**: Written in Go with full type safety

## Currently Supported

### Queries

- `viewer` - Returns the authenticated user

### User Fields

The following scalar fields are supported on the `User` type:

- `id` (String!) - The user's ID
- `login` (String!) - The username
- `name` (String) - The user's full name
- `email` (String!) - The user's email (may be empty string)
- `avatarUrl` (String!) - URL to the user's avatar
- `bio` (String) - The user's bio
- `company` (String) - The user's company
- `location` (String) - The user's location

## Installation

```bash
# Clone the repository
git clone https://github.com/srt32/ghgraph.git
cd ghgraph

# Install dependencies
go mod download

# Build the server
go build -o ghgraph ./cmd/server
```

## Usage

### Starting the Server

```bash
# Default port (8080)
./ghgraph

# Custom port
PORT=3000 ./ghgraph
```

### Making Queries

The server accepts GraphQL queries at the `/graphql` endpoint. You must provide a GitHub personal access token in the `Authorization` header.

#### Example: Query All Fields

```bash
curl -X POST http://localhost:8080/graphql \
  -H "Authorization: Bearer YOUR_GITHUB_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "{ viewer { login name email avatarUrl bio company location } }"
  }'
```

#### Example: Query Specific Fields

```bash
curl -X POST http://localhost:8080/graphql \
  -H "Authorization: Bearer YOUR_GITHUB_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "{ viewer { login name } }"
  }'
```

#### Example: Using Variables

```bash
curl -X POST http://localhost:8080/graphql \
  -H "Authorization: Bearer YOUR_GITHUB_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "query GetUser { viewer { login email } }",
    "operationName": "GetUser"
  }'
```

### Health Check

```bash
curl http://localhost:8080/health
```

## Development

### Running Tests

```bash
# Run all unit tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for a specific package
go test ./internal/github/...
go test ./internal/graphql/...
```

### Running Integration Tests

Integration tests verify that our proxy matches the behavior of the real GitHub GraphQL API.

```bash
# Set your GitHub token
export GITHUB_TOKEN=your_token_here

# Run integration tests
go test -v -run TestIntegration

# Run all tests including integration
go test -v ./...
```

## Architecture

```
/cmd/server              - Main application entry point
/internal/github         - GitHub REST API client
/internal/graphql        - GraphQL schema and resolvers
/internal/server         - HTTP server and request handling
/integration_test.go     - Integration tests against real GitHub API
```

### How It Works

1. Client sends GraphQL query to `/graphql` endpoint
2. Server extracts GitHub token from `Authorization` header
3. GraphQL query is parsed and validated against schema
4. Resolvers translate the query into GitHub REST API calls
5. REST API responses are mapped back to GraphQL response format
6. Response is returned to client

## Project Structure

```
ghgraph/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── github/
│   │   ├── client.go            # REST API client
│   │   └── client_test.go       # Client tests
│   ├── graphql/
│   │   ├── schema.go            # GraphQL schema definition
│   │   └── schema_test.go       # Schema tests
│   └── server/
│       └── server.go            # HTTP server
├── integration_test.go          # Integration tests
├── go.mod                       # Go module file
└── README.md                    # This file
```

## Roadmap

### Next Steps

- [ ] Add support for nested queries (e.g., `viewer.repositories`)
- [ ] Implement pagination for list fields
- [ ] Add more root queries (e.g., `user`, `repository`, `organization`)
- [ ] Support mutations
- [ ] Add request caching
- [ ] Implement rate limiting

### Future Enhancements

- [ ] Support for GitHub Enterprise
- [ ] Metrics and monitoring
- [ ] Request/response logging
- [ ] GraphQL subscriptions (via polling REST endpoints)
- [ ] Admin UI for testing queries

## Testing Philosophy

This project includes two types of tests:

1. **Unit Tests**: Mock the GitHub REST API to test our translation logic in isolation
2. **Integration Tests**: Query the real GitHub GraphQL API to ensure our responses match the actual API behavior

Integration tests require a `GITHUB_TOKEN` environment variable and will be skipped if not provided.

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## License

MIT License - See LICENSE file for details

## References

- [GitHub GraphQL API Documentation](https://docs.github.com/en/graphql/reference)
- [GitHub REST API Documentation](https://docs.github.com/en/rest)
- [GraphQL Go Library](https://github.com/graphql-go/graphql)
