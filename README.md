# ghgraph

A stateless GraphQL proxy that translates GitHub GraphQL API queries into GitHub REST API calls.

## Overview

This project clones the GitHub GraphQL API interface but translates all queries to REST API calls behind the scenes. It's designed to work as if the GitHub GraphQL API doesn't exist and only the REST API is available.

## Features

- **GraphQL Interface**: Exposes a GraphQL API that matches GitHub's schema
- **REST Backend**: Translates all queries to GitHub REST API v3 calls
- **Intelligent Caching**: Token-scoped in-memory caching with TTL for improved performance
- **Pagination Support**: Cursor-based pagination for repositories, issues, and pull requests
- **State Filtering**: Filter issues and pull requests by state (OPEN, CLOSED, MERGED)
- **Type-Safe**: Written in Go with full type safety
- **Multi-Tenant**: Isolated cache per user token for security

## Currently Supported

### Queries

- `viewer` - Returns the authenticated user
- `user(login: String!)` - Lookup a user by login
- `repository(owner: String!, name: String!)` - Lookup a repository by owner and name
- `organization(login: String!)` - Lookup an organization by login

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
- `repositories` (RepositoryConnection!) - A paginated list of repositories owned by the user
  - Arguments:
    - `first` (Int) - Number of repositories to return (forward pagination)
    - `after` (String) - Cursor to start after for forward pagination
    - `last` (Int) - Number of repositories to return (backward pagination)
    - `before` (String) - Cursor to start before for backward pagination

### Repository Fields

The following fields are supported on the `Repository` type:

- `id` (String!) - The repository's ID
- `name` (String!) - The name of the repository
- `nameWithOwner` (String!) - The repository's name with owner (e.g., "owner/repo")
- `description` (String) - The description of the repository
- `isPrivate` (Boolean!) - Whether the repository is private
- `url` (String!) - The HTTP URL for this repository
- `stargazerCount` (Int!) - Number of stars
- `forkCount` (Int!) - Number of forks
- `defaultBranchRef` (String) - Name of the default branch
- `owner` (User!) - The User who owns the repository
- `issues` (IssueConnection!) - A paginated list of issues for the repository
  - Arguments:
    - `first` (Int) - Number of issues to return (forward pagination)
    - `after` (String) - Cursor to start after for forward pagination
    - `last` (Int) - Number of issues to return (backward pagination)
    - `before` (String) - Cursor to start before for backward pagination
    - `states` ([IssueState]) - Filter by issue states (OPEN, CLOSED)
- `pullRequests` (PullRequestConnection!) - A paginated list of pull requests for the repository
  - Arguments:
    - `first` (Int) - Number of PRs to return (forward pagination)
    - `after` (String) - Cursor to start after for forward pagination
    - `last` (Int) - Number of PRs to return (backward pagination)
    - `before` (String) - Cursor to start before for backward pagination
    - `states` ([PullRequestState]) - Filter by PR states (OPEN, CLOSED, MERGED)

### Organization Fields

The following fields are supported on the `Organization` type:

- `id` (String!) - The organization's ID
- `login` (String!) - The organization's login name
- `name` (String) - The organization's public profile name
- `description` (String) - The organization's description
- `avatarUrl` (String!) - URL to the organization's avatar
- `url` (String!) - The HTTP URL for this organization
- `email` (String) - The organization's public email
- `location` (String) - The organization's location

### Issue Fields

The following fields are supported on the `Issue` type:

- `id` (String!) - The issue's ID
- `number` (Int!) - The issue number in the repository
- `title` (String!) - The title of the issue
- `body` (String) - The body content of the issue
- `state` (IssueState!) - The state of the issue (OPEN, CLOSED)
- `url` (String!) - The HTTP URL for this issue
- `author` (User!) - The user who created the issue
- `createdAt` (String!) - When the issue was created

### PullRequest Fields

The following fields are supported on the `PullRequest` type:

- `id` (String!) - The pull request's ID
- `number` (Int!) - The pull request number in the repository
- `title` (String!) - The title of the pull request
- `body` (String) - The body content of the pull request
- `state` (PullRequestState!) - The state of the pull request (OPEN, CLOSED, MERGED)
- `url` (String!) - The HTTP URL for this pull request
- `author` (User!) - The user who created the pull request
- `createdAt` (String!) - When the pull request was created
- `merged` (Boolean!) - Whether the pull request was merged

### Pagination Types

The following types are used for cursor-based pagination:

#### PageInfo

Contains pagination metadata:

- `hasNextPage` (Boolean!) - Are there more items when paginating forward?
- `hasPreviousPage` (Boolean!) - Are there more items when paginating backward?
- `startCursor` (String) - Cursor to continue from the start
- `endCursor` (String) - Cursor to continue from the end

#### RepositoryConnection

Represents a paginated list of repositories:

- `edges` ([RepositoryEdge]) - List of edges containing repositories and cursors
- `nodes` ([Repository]) - Direct list of repositories
- `pageInfo` (PageInfo!) - Pagination information
- `totalCount` (Int!) - Count of repositories in the current page

**Note:** `totalCount` returns the number of items in the current page, not the total count across all pages. This is a limitation of the GitHub REST API, which doesn't provide total counts in paginated responses.

#### RepositoryEdge

Represents an edge in a repository connection:

- `cursor` (String!) - Cursor for this edge
- `node` (Repository) - The repository at this edge

#### IssueConnection

Represents a paginated list of issues:

- `edges` ([IssueEdge]) - List of edges containing issues and cursors
- `nodes` ([Issue]) - Direct list of issues
- `pageInfo` (PageInfo!) - Pagination information
- `totalCount` (Int!) - Count of issues in the current page

**Note:** `totalCount` returns the number of items in the current page, not the total count across all pages.

#### IssueEdge

Represents an edge in an issue connection:

- `cursor` (String!) - Cursor for this edge
- `node` (Issue) - The issue at this edge

#### PullRequestConnection

Represents a paginated list of pull requests:

- `edges` ([PullRequestEdge]) - List of edges containing pull requests and cursors
- `nodes` ([PullRequest]) - Direct list of pull requests
- `pageInfo` (PageInfo!) - Pagination information
- `totalCount` (Int!) - Count of pull requests in the current page

**Note:** `totalCount` returns the number of items in the current page, not the total count across all pages.

#### PullRequestEdge

Represents an edge in a pull request connection:

- `cursor` (String!) - Cursor for this edge
- `node` (PullRequest) - The pull request at this edge

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

#### Example: Query Specific User

```bash
curl -X POST http://localhost:8080/graphql \
  -H "Authorization: Bearer YOUR_GITHUB_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "{ user(login: \"octocat\") { login name bio company } }"
  }'
```

#### Example: Query Repository with Nested Owner

```bash
curl -X POST http://localhost:8080/graphql \
  -H "Authorization: Bearer YOUR_GITHUB_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "{ repository(owner: \"facebook\", name: \"react\") { name nameWithOwner description stargazerCount owner { login name } } }"
  }'
```

#### Example: Query Organization

```bash
curl -X POST http://localhost:8080/graphql \
  -H "Authorization: Bearer YOUR_GITHUB_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "{ organization(login: \"github\") { login name description location } }"
  }'
```

#### Example: Query User Repositories with Pagination

```bash
curl -X POST http://localhost:8080/graphql \
  -H "Authorization: Bearer YOUR_GITHUB_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "{ viewer { login repositories(first: 5) { totalCount pageInfo { hasNextPage endCursor } nodes { name description stargazerCount } } } }"
  }'
```

#### Example: Paginate Through Repositories

```bash
# First page
curl -X POST http://localhost:8080/graphql \
  -H "Authorization: Bearer YOUR_GITHUB_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "{ viewer { repositories(first: 2) { pageInfo { hasNextPage endCursor } edges { cursor node { name } } } } }"
  }'

# Next page (use endCursor from previous response)
curl -X POST http://localhost:8080/graphql \
  -H "Authorization: Bearer YOUR_GITHUB_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "query($cursor: String!) { viewer { repositories(first: 2, after: $cursor) { pageInfo { hasNextPage endCursor } nodes { name } } } }",
    "variables": {
      "cursor": "Mg=="
    }
  }'
```

#### Example: Query Another User's Repositories

```bash
curl -X POST http://localhost:8080/graphql \
  -H "Authorization: Bearer YOUR_GITHUB_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "{ user(login: \"octocat\") { login repositories(first: 10) { totalCount nodes { name stargazerCount } } } }"
  }'
```

#### Example: Query Repository Issues

```bash
curl -X POST http://localhost:8080/graphql \
  -H "Authorization: Bearer YOUR_GITHUB_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "{ repository(owner: \"facebook\", name: \"react\") { name issues(first: 5, states: [OPEN]) { totalCount pageInfo { hasNextPage endCursor } nodes { number title state author { login } createdAt } } } }"
  }'
```

#### Example: Query Repository Pull Requests

```bash
curl -X POST http://localhost:8080/graphql \
  -H "Authorization: Bearer YOUR_GITHUB_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "{ repository(owner: \"facebook\", name: \"react\") { name pullRequests(first: 10, states: [MERGED]) { totalCount nodes { number title state merged author { login } } } } }"
  }'
```

#### Example: Paginate Through Issues

```bash
# First page
curl -X POST http://localhost:8080/graphql \
  -H "Authorization: Bearer YOUR_GITHUB_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "{ repository(owner: \"octocat\", name: \"hello-world\") { issues(first: 2) { pageInfo { hasNextPage endCursor } edges { cursor node { number title state } } } } }"
  }'

# Next page (use endCursor from previous response)
curl -X POST http://localhost:8080/graphql \
  -H "Authorization: Bearer YOUR_GITHUB_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "query($cursor: String!) { repository(owner: \"octocat\", name: \"hello-world\") { issues(first: 2, after: $cursor) { pageInfo { hasNextPage hasPreviousPage endCursor } nodes { number title } } } }",
    "variables": {
      "cursor": "Mg=="
    }
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
/internal/cache          - Token-scoped TTL caching layer
/internal/github         - GitHub REST API client with caching
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

### Caching

The proxy implements an intelligent in-memory caching layer to reduce GitHub API calls and improve performance:

#### Features

- **Token-Scoped Caching**: Each user's token has isolated cache storage for security and multi-tenancy
- **TTL-Based Expiration**: Cache entries automatically expire after 5 minutes
- **Automatic Cleanup**: Background goroutine removes expired entries
- **Thread-Safe**: Uses sync.Map for concurrent access
- **Privacy**: Tokens are SHA256-hashed (first 8 bytes) for cache key generation

#### What Gets Cached

The following GitHub REST API responses are cached:
- Repository issues
- Pull requests
- User repositories (future)

#### Cache Keys

Cache keys are structured as: `<token-hash>:<resource-type>:<owner>/<repo>:<state>:page<N>:per<M>`

Example: `a1b2c3d4:issues:facebook/react:open:page1:per30`

#### Benefits

- Reduces redundant API calls
- Stays within GitHub's rate limits
- Improves response times for repeated queries
- Supports multi-tenant usage with token isolation

## Project Structure

```
ghgraph/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── cache/
│   │   ├── cache.go             # Token-scoped TTL cache
│   │   └── cache_test.go        # Cache tests
│   ├── github/
│   │   ├── client.go            # REST API client with caching
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

### Completed

- [x] Add more root queries (`user`, `repository`, `organization`)
- [x] Add support for nested queries (e.g., `viewer.repositories`, `user.repositories`)
- [x] Implement cursor-based pagination for list fields
- [x] Support for repository connections with edges, nodes, and pageInfo
- [x] Add paginated fields for issues and pull requests (`repository.issues`, `repository.pullRequests`)
- [x] Implement token-scoped request caching with TTL
- [x] Add state filtering for issues and pull requests

### Next Steps

- [ ] Add more paginated fields (e.g., `user.followers`, `organization.repositories`)
- [ ] Support mutations (e.g., creating issues, commenting on PRs)
- [ ] Implement rate limiting
- [ ] Add support for comments on issues and pull requests

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
