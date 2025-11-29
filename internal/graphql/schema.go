package graphql

import (
	"context"
	"encoding/base64"
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/srt32/ghgraph/internal/github"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const githubClientKey contextKey = "githubClient"

// getClientFromContext retrieves the GitHub client from the request context
func getClientFromContext(ctx context.Context) (*github.Client, bool) {
	client, ok := ctx.Value(githubClientKey).(*github.Client)
	return client, ok
}

// pageInfoType represents pagination information
var pageInfoType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "PageInfo",
	Description: "Information about pagination in a connection.",
	Fields: graphql.Fields{
		"hasNextPage": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "When paginating forwards, are there more items?",
		},
		"hasPreviousPage": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "When paginating backwards, are there more items?",
		},
		"startCursor": &graphql.Field{
			Type:        graphql.String,
			Description: "When paginating backwards, the cursor to continue.",
		},
		"endCursor": &graphql.Field{
			Type:        graphql.String,
			Description: "When paginating forwards, the cursor to continue.",
		},
	},
})

// repositoryEdgeType represents an edge in a repository connection
var repositoryEdgeType *graphql.Object

// repositoryConnectionType represents a paginated list of repositories
var repositoryConnectionType *graphql.Object

// issueEdgeType represents an edge in an issue connection
var issueEdgeType *graphql.Object

// issueConnectionType represents a paginated list of issues
var issueConnectionType *graphql.Object

// pullRequestEdgeType represents an edge in a pull request connection
var pullRequestEdgeType *graphql.Object

// pullRequestConnectionType represents a paginated list of pull requests
var pullRequestConnectionType *graphql.Object

// issueStateEnum represents the state of an issue
var issueStateEnum = graphql.NewEnum(graphql.EnumConfig{
	Name:        "IssueState",
	Description: "The possible states of an issue.",
	Values: graphql.EnumValueConfigMap{
		"OPEN": &graphql.EnumValueConfig{
			Value:       "open",
			Description: "An issue that is still open",
		},
		"CLOSED": &graphql.EnumValueConfig{
			Value:       "closed",
			Description: "An issue that has been closed",
		},
	},
})

// pullRequestStateEnum represents the state of a pull request
var pullRequestStateEnum = graphql.NewEnum(graphql.EnumConfig{
	Name:        "PullRequestState",
	Description: "The possible states of a pull request.",
	Values: graphql.EnumValueConfigMap{
		"OPEN": &graphql.EnumValueConfig{
			Value:       "open",
			Description: "A pull request that is still open",
		},
		"CLOSED": &graphql.EnumValueConfig{
			Value:       "closed",
			Description: "A pull request that has been closed",
		},
		"MERGED": &graphql.EnumValueConfig{
			Value:       "merged",
			Description: "A pull request that has been merged",
		},
	},
})

// userType represents the GraphQL User type matching GitHub's schema
var userType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "User",
	Description: "A user is an individual's account on GitHub that owns repositories and can make new content.",
	Fields: graphql.FieldsThunk(func() graphql.Fields {
		return graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The Node ID of the User object",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if user, ok := p.Source.(*github.User); ok {
					// GitHub GraphQL uses global node IDs, but for now we'll use the numeric ID
					// In a full implementation, we'd convert this to a proper base64-encoded node ID
					return user.Login, nil
				}
				return nil, nil
			},
		},
		"login": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The username used to login.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if user, ok := p.Source.(*github.User); ok {
					return user.Login, nil
				}
				return nil, nil
			},
		},
		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "The user's public profile name.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if user, ok := p.Source.(*github.User); ok {
					return user.Name, nil
				}
				return nil, nil
			},
		},
		"email": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The user's publicly visible profile email.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if user, ok := p.Source.(*github.User); ok {
					if user.Email != nil {
						return user.Email, nil
					}
					// Return empty string if no email (to match NonNull requirement)
					return "", nil
				}
				return "", nil
			},
		},
		"avatarUrl": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "A URL pointing to the user's public avatar.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if user, ok := p.Source.(*github.User); ok {
					return user.AvatarURL, nil
				}
				return nil, nil
			},
		},
		"bio": &graphql.Field{
			Type:        graphql.String,
			Description: "The user's public profile bio.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if user, ok := p.Source.(*github.User); ok {
					return user.Bio, nil
				}
				return nil, nil
			},
		},
		"company": &graphql.Field{
			Type:        graphql.String,
			Description: "The user's public profile company.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if user, ok := p.Source.(*github.User); ok {
					return user.Company, nil
				}
				return nil, nil
			},
		},
		"location": &graphql.Field{
			Type:        graphql.String,
			Description: "The user's public profile location.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if user, ok := p.Source.(*github.User); ok {
					return user.Location, nil
				}
				return nil, nil
			},
		},
		"repositories": &graphql.Field{
			Type:        repositoryConnectionType,
			Description: "A list of repositories that the user owns.",
			Args: graphql.FieldConfigArgument{
				"first": &graphql.ArgumentConfig{
					Type:        graphql.Int,
					Description: "Returns the first n repositories from the list.",
				},
				"after": &graphql.ArgumentConfig{
					Type:        graphql.String,
					Description: "Returns the repositories that come after the specified cursor.",
				},
				"last": &graphql.ArgumentConfig{
					Type:        graphql.Int,
					Description: "Returns the last n repositories from the list.",
				},
				"before": &graphql.ArgumentConfig{
					Type:        graphql.String,
					Description: "Returns the repositories that come before the specified cursor.",
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				user, ok := p.Source.(*github.User)
				if !ok {
					return nil, nil
				}

				// Get pagination arguments
				var first, last *int
				var after, before *string

				if firstVal, ok := p.Args["first"].(int); ok {
					first = &firstVal
				}
				if afterVal, ok := p.Args["after"].(string); ok && afterVal != "" {
					after = &afterVal
				}
				if lastVal, ok := p.Args["last"].(int); ok {
					last = &lastVal
				}
				if beforeVal, ok := p.Args["before"].(string); ok && beforeVal != "" {
					before = &beforeVal
				}

				// Get GitHub client from context
				githubClient, ok := getClientFromContext(p.Context)
				if !ok {
					return nil, nil
				}

				// Fetch repositories
				result, err := githubClient.ListUserRepositories(user.Login, first, after, last, before)
				if err != nil {
					return nil, err
				}

				// Build edges
				edges := make([]map[string]interface{}, 0, len(result.Repositories))
				for i, repo := range result.Repositories {
					// Create cursor for this item
					cursor := base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(i)))
					edges = append(edges, map[string]interface{}{
						"cursor": cursor,
						"node":   repo,
					})
				}

				// Build connection
				return map[string]interface{}{
					"edges": edges,
					"nodes": result.Repositories,
					"pageInfo": map[string]interface{}{
						"hasNextPage":     result.HasNextPage,
						"hasPreviousPage": result.HasPrevPage,
						"startCursor":     result.StartCursor,
						"endCursor":       result.EndCursor,
					},
					"totalCount": result.TotalCount,
				}, nil
			},
		},
		}
	}),
})

// repositoryType represents the GraphQL Repository type matching GitHub's schema
var repositoryType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Repository",
	Description: "A repository contains the content for a project.",
	Fields: graphql.FieldsThunk(func() graphql.Fields {
		return graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The Node ID of the Repository object",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if repo, ok := p.Source.(*github.Repository); ok {
					return repo.Name, nil
				}
				return nil, nil
			},
		},
		"name": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The name of the repository.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if repo, ok := p.Source.(*github.Repository); ok {
					return repo.Name, nil
				}
				return nil, nil
			},
		},
		"nameWithOwner": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The repository's name with owner.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if repo, ok := p.Source.(*github.Repository); ok {
					return repo.FullName, nil
				}
				return nil, nil
			},
		},
		"description": &graphql.Field{
			Type:        graphql.String,
			Description: "The description of the repository.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if repo, ok := p.Source.(*github.Repository); ok {
					return repo.Description, nil
				}
				return nil, nil
			},
		},
		"isPrivate": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "Identifies if the repository is private.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if repo, ok := p.Source.(*github.Repository); ok {
					return repo.Private, nil
				}
				return false, nil
			},
		},
		"url": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The HTTP URL for this repository",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if repo, ok := p.Source.(*github.Repository); ok {
					return repo.HTMLURL, nil
				}
				return nil, nil
			},
		},
		"stargazerCount": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.Int),
			Description: "Returns a count of how many stargazers there are on this repository",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if repo, ok := p.Source.(*github.Repository); ok {
					return repo.StargazersCount, nil
				}
				return 0, nil
			},
		},
		"forkCount": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.Int),
			Description: "Returns how many forks there are of this repository in the whole network.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if repo, ok := p.Source.(*github.Repository); ok {
					return repo.ForksCount, nil
				}
				return 0, nil
			},
		},
		"defaultBranchRef": &graphql.Field{
			Type:        graphql.String,
			Description: "The name of the default branch (simplified, returns just the branch name).",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if repo, ok := p.Source.(*github.Repository); ok {
					return repo.DefaultBranch, nil
				}
				return nil, nil
			},
		},
		"owner": &graphql.Field{
			Type:        userType,
			Description: "The User owner of the repository.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if repo, ok := p.Source.(*github.Repository); ok {
					return &repo.Owner, nil
				}
				return nil, nil
			},
		},
		"issues": &graphql.Field{
			Type:        issueConnectionType,
			Description: "A list of issues that have been opened in the repository.",
			Args: graphql.FieldConfigArgument{
				"first": &graphql.ArgumentConfig{
					Type:        graphql.Int,
					Description: "Returns the first n issues.",
				},
				"after": &graphql.ArgumentConfig{
					Type:        graphql.String,
					Description: "Returns the issues after the specified cursor.",
				},
				"last": &graphql.ArgumentConfig{
					Type:        graphql.Int,
					Description: "Returns the last n issues.",
				},
				"before": &graphql.ArgumentConfig{
					Type:        graphql.String,
					Description: "Returns the issues before the specified cursor.",
				},
				"states": &graphql.ArgumentConfig{
					Type:        graphql.NewList(issueStateEnum),
					Description: "Filter issues by state.",
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				repo, ok := p.Source.(*github.Repository)
				if !ok {
					return nil, nil
				}

				// Get pagination arguments
				var first, last *int
				var after, before *string
				if firstVal, ok := p.Args["first"].(int); ok {
					first = &firstVal
				}
				if afterVal, ok := p.Args["after"].(string); ok && afterVal != "" {
					after = &afterVal
				}
				if lastVal, ok := p.Args["last"].(int); ok {
					last = &lastVal
				}
				if beforeVal, ok := p.Args["before"].(string); ok && beforeVal != "" {
					before = &beforeVal
				}

				// Get state filter (default to "open")
				state := "open"
				if states, ok := p.Args["states"].([]interface{}); ok && len(states) > 0 {
					// For simplicity, use the first state provided
					state = states[0].(string)
				}

				githubClient, ok := getClientFromContext(p.Context)
				if !ok {
					return nil, nil
				}

				result, err := githubClient.ListRepositoryIssues(repo.Owner.Login, repo.Name, state, first, after, last, before)
				if err != nil {
					return nil, err
				}

				// Build edges
				edges := make([]map[string]interface{}, 0, len(result.Issues))
				for i, issue := range result.Issues {
					cursor := base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(i)))
					edges = append(edges, map[string]interface{}{
						"cursor": cursor,
						"node":   issue,
					})
				}

				return map[string]interface{}{
					"edges": edges,
					"nodes": result.Issues,
					"pageInfo": map[string]interface{}{
						"hasNextPage":     result.HasNextPage,
						"hasPreviousPage": result.HasPrevPage,
						"startCursor":     result.StartCursor,
						"endCursor":       result.EndCursor,
					},
					"totalCount": result.TotalCount,
				}, nil
			},
		},
		"pullRequests": &graphql.Field{
			Type:        pullRequestConnectionType,
			Description: "A list of pull requests that have been opened in the repository.",
			Args: graphql.FieldConfigArgument{
				"first": &graphql.ArgumentConfig{
					Type:        graphql.Int,
					Description: "Returns the first n pull requests.",
				},
				"after": &graphql.ArgumentConfig{
					Type:        graphql.String,
					Description: "Returns the pull requests after the specified cursor.",
				},
				"last": &graphql.ArgumentConfig{
					Type:        graphql.Int,
					Description: "Returns the last n pull requests.",
				},
				"before": &graphql.ArgumentConfig{
					Type:        graphql.String,
					Description: "Returns the pull requests before the specified cursor.",
				},
				"states": &graphql.ArgumentConfig{
					Type:        graphql.NewList(pullRequestStateEnum),
					Description: "Filter pull requests by state.",
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				repo, ok := p.Source.(*github.Repository)
				if !ok {
					return nil, nil
				}

				// Get pagination arguments
				var first, last *int
				var after, before *string
				if firstVal, ok := p.Args["first"].(int); ok {
					first = &firstVal
				}
				if afterVal, ok := p.Args["after"].(string); ok && afterVal != "" {
					after = &afterVal
				}
				if lastVal, ok := p.Args["last"].(int); ok {
					last = &lastVal
				}
				if beforeVal, ok := p.Args["before"].(string); ok && beforeVal != "" {
					before = &beforeVal
				}

				// Get state filter (default to "open")
				state := "open"
				if states, ok := p.Args["states"].([]interface{}); ok && len(states) > 0 {
					state = states[0].(string)
				}

				githubClient, ok := getClientFromContext(p.Context)
				if !ok {
					return nil, nil
				}

				result, err := githubClient.ListRepositoryPullRequests(repo.Owner.Login, repo.Name, state, first, after, last, before)
				if err != nil {
					return nil, err
				}

				// Build edges
				edges := make([]map[string]interface{}, 0, len(result.PullRequests))
				for i, pr := range result.PullRequests {
					cursor := base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(i)))
					edges = append(edges, map[string]interface{}{
						"cursor": cursor,
						"node":   pr,
					})
				}

				return map[string]interface{}{
					"edges": edges,
					"nodes": result.PullRequests,
					"pageInfo": map[string]interface{}{
						"hasNextPage":     result.HasNextPage,
						"hasPreviousPage": result.HasPrevPage,
						"startCursor":     result.StartCursor,
						"endCursor":       result.EndCursor,
					},
					"totalCount": result.TotalCount,
				}, nil
			},
		},
		}
	}),
})

// organizationType represents the GraphQL Organization type matching GitHub's schema
var organizationType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Organization",
	Description: "An organization is a collection of teams and repositories.",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The Node ID of the Organization object",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if org, ok := p.Source.(*github.Organization); ok {
					return org.Login, nil
				}
				return nil, nil
			},
		},
		"login": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The organization's login name.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if org, ok := p.Source.(*github.Organization); ok {
					return org.Login, nil
				}
				return nil, nil
			},
		},
		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "The organization's public profile name.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if org, ok := p.Source.(*github.Organization); ok {
					return org.Name, nil
				}
				return nil, nil
			},
		},
		"description": &graphql.Field{
			Type:        graphql.String,
			Description: "The organization's public profile description.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if org, ok := p.Source.(*github.Organization); ok {
					return org.Description, nil
				}
				return nil, nil
			},
		},
		"avatarUrl": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "A URL pointing to the organization's public avatar.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if org, ok := p.Source.(*github.Organization); ok {
					return org.AvatarURL, nil
				}
				return nil, nil
			},
		},
		"url": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The HTTP URL for this organization.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if org, ok := p.Source.(*github.Organization); ok {
					return org.HTMLURL, nil
				}
				return nil, nil
			},
		},
		"email": &graphql.Field{
			Type:        graphql.String,
			Description: "The organization's public email.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if org, ok := p.Source.(*github.Organization); ok {
					return org.Email, nil
				}
				return nil, nil
			},
		},
		"location": &graphql.Field{
			Type:        graphql.String,
			Description: "The organization's public profile location.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if org, ok := p.Source.(*github.Organization); ok {
					return org.Location, nil
				}
				return nil, nil
			},
		},
	},
})

// issueType represents the GraphQL Issue type
var issueType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Issue",
	Description: "An issue is a discussion thread about a topic",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The ID of the issue",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if issue, ok := p.Source.(*github.Issue); ok {
					return strconv.FormatInt(issue.ID, 10), nil
				}
				return nil, nil
			},
		},
		"number": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.Int),
			Description: "The issue number",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if issue, ok := p.Source.(*github.Issue); ok {
					return issue.Number, nil
				}
				return nil, nil
			},
		},
		"title": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The title of the issue",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if issue, ok := p.Source.(*github.Issue); ok {
					return issue.Title, nil
				}
				return nil, nil
			},
		},
		"body": &graphql.Field{
			Type:        graphql.String,
			Description: "The body of the issue",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if issue, ok := p.Source.(*github.Issue); ok {
					return issue.Body, nil
				}
				return nil, nil
			},
		},
		"state": &graphql.Field{
			Type:        graphql.NewNonNull(issueStateEnum),
			Description: "The state of the issue",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if issue, ok := p.Source.(*github.Issue); ok {
					return issue.State, nil
				}
				return nil, nil
			},
		},
		"url": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The HTTP URL for this issue",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if issue, ok := p.Source.(*github.Issue); ok {
					return issue.HTMLURL, nil
				}
				return nil, nil
			},
		},
		"author": &graphql.Field{
			Type:        userType,
			Description: "The user who created the issue",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if issue, ok := p.Source.(*github.Issue); ok {
					return &issue.User, nil
				}
				return nil, nil
			},
		},
		"createdAt": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The date and time the issue was created",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if issue, ok := p.Source.(*github.Issue); ok {
					return issue.CreatedAt, nil
				}
				return nil, nil
			},
		},
	},
})

// pullRequestType represents the GraphQL PullRequest type
var pullRequestType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "PullRequest",
	Description: "A pull request is a proposal to merge a set of changes",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The ID of the pull request",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if pr, ok := p.Source.(*github.PullRequest); ok {
					return strconv.FormatInt(pr.ID, 10), nil
				}
				return nil, nil
			},
		},
		"number": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.Int),
			Description: "The pull request number",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if pr, ok := p.Source.(*github.PullRequest); ok {
					return pr.Number, nil
				}
				return nil, nil
			},
		},
		"title": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The title of the pull request",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if pr, ok := p.Source.(*github.PullRequest); ok {
					return pr.Title, nil
				}
				return nil, nil
			},
		},
		"body": &graphql.Field{
			Type:        graphql.String,
			Description: "The body of the pull request",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if pr, ok := p.Source.(*github.PullRequest); ok {
					return pr.Body, nil
				}
				return nil, nil
			},
		},
		"state": &graphql.Field{
			Type:        graphql.NewNonNull(pullRequestStateEnum),
			Description: "The state of the pull request",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if pr, ok := p.Source.(*github.PullRequest); ok {
					return pr.State, nil
				}
				return nil, nil
			},
		},
		"merged": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "Whether the pull request has been merged",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if pr, ok := p.Source.(*github.PullRequest); ok {
					return pr.Merged, nil
				}
				return nil, nil
			},
		},
		"url": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The HTTP URL for this pull request",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if pr, ok := p.Source.(*github.PullRequest); ok {
					return pr.HTMLURL, nil
				}
				return nil, nil
			},
		},
		"author": &graphql.Field{
			Type:        userType,
			Description: "The user who created the pull request",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if pr, ok := p.Source.(*github.PullRequest); ok {
					return &pr.User, nil
				}
				return nil, nil
			},
		},
		"createdAt": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The date and time the pull request was created",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if pr, ok := p.Source.(*github.PullRequest); ok {
					return pr.CreatedAt, nil
				}
				return nil, nil
			},
		},
	},
})

func init() {
	// Initialize repositoryEdgeType with reference to repositoryType
	repositoryEdgeType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "RepositoryEdge",
		Description: "An edge in a repository connection.",
		Fields: graphql.Fields{
			"cursor": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "A cursor for use in pagination.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if edge, ok := p.Source.(map[string]interface{}); ok {
						return edge["cursor"], nil
					}
					return nil, nil
				},
			},
			"node": &graphql.Field{
				Type:        repositoryType,
				Description: "The item at the end of the edge.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if edge, ok := p.Source.(map[string]interface{}); ok {
						return edge["node"], nil
					}
					return nil, nil
				},
			},
		},
	})

	// Initialize repositoryConnectionType with references to edge and pageInfo
	repositoryConnectionType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "RepositoryConnection",
		Description: "A list of repositories.",
		Fields: graphql.Fields{
			"edges": &graphql.Field{
				Type:        graphql.NewList(repositoryEdgeType),
				Description: "A list of edges.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if conn, ok := p.Source.(map[string]interface{}); ok {
						return conn["edges"], nil
					}
					return nil, nil
				},
			},
			"nodes": &graphql.Field{
				Type:        graphql.NewList(repositoryType),
				Description: "A list of nodes.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if conn, ok := p.Source.(map[string]interface{}); ok {
						return conn["nodes"], nil
					}
					return nil, nil
				},
			},
			"pageInfo": &graphql.Field{
				Type:        graphql.NewNonNull(pageInfoType),
				Description: "Information to aid in pagination.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if conn, ok := p.Source.(map[string]interface{}); ok {
						return conn["pageInfo"], nil
					}
					return nil, nil
				},
			},
			"totalCount": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Identifies the total count of items in the connection.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if conn, ok := p.Source.(map[string]interface{}); ok {
						return conn["totalCount"], nil
					}
					return 0, nil
				},
			},
		},
	})

	// Initialize issueEdgeType
	issueEdgeType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "IssueEdge",
		Description: "An edge in an issue connection.",
		Fields: graphql.Fields{
			"cursor": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "A cursor for use in pagination.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if edge, ok := p.Source.(map[string]interface{}); ok {
						return edge["cursor"], nil
					}
					return nil, nil
				},
			},
			"node": &graphql.Field{
				Type:        issueType,
				Description: "The issue at the end of the edge.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if edge, ok := p.Source.(map[string]interface{}); ok {
						return edge["node"], nil
					}
					return nil, nil
				},
			},
		},
	})

	// Initialize issueConnectionType
	issueConnectionType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "IssueConnection",
		Description: "A list of issues.",
		Fields: graphql.Fields{
			"edges": &graphql.Field{
				Type:        graphql.NewList(issueEdgeType),
				Description: "A list of edges.",
			},
			"nodes": &graphql.Field{
				Type:        graphql.NewList(issueType),
				Description: "A list of nodes.",
			},
			"pageInfo": &graphql.Field{
				Type:        graphql.NewNonNull(pageInfoType),
				Description: "Information to aid in pagination.",
			},
			"totalCount": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Identifies the total count of items in the connection.",
			},
		},
	})

	// Initialize pullRequestEdgeType
	pullRequestEdgeType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "PullRequestEdge",
		Description: "An edge in a pull request connection.",
		Fields: graphql.Fields{
			"cursor": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "A cursor for use in pagination.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if edge, ok := p.Source.(map[string]interface{}); ok {
						return edge["cursor"], nil
					}
					return nil, nil
				},
			},
			"node": &graphql.Field{
				Type:        pullRequestType,
				Description: "The pull request at the end of the edge.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if edge, ok := p.Source.(map[string]interface{}); ok {
						return edge["node"], nil
					}
					return nil, nil
				},
			},
		},
	})

	// Initialize pullRequestConnectionType
	pullRequestConnectionType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "PullRequestConnection",
		Description: "A list of pull requests.",
		Fields: graphql.Fields{
			"edges": &graphql.Field{
				Type:        graphql.NewList(pullRequestEdgeType),
				Description: "A list of edges.",
			},
			"nodes": &graphql.Field{
				Type:        graphql.NewList(pullRequestType),
				Description: "A list of nodes.",
			},
			"pageInfo": &graphql.Field{
				Type:        graphql.NewNonNull(pageInfoType),
				Description: "Information to aid in pagination.",
			},
			"totalCount": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Identifies the total count of items in the connection.",
			},
		},
	})
}

var cachedSchema graphql.Schema
var schemaInitialized bool

// GetSchema returns the GraphQL schema, creating it once on first call
func GetSchema() (graphql.Schema, error) {
	if schemaInitialized {
		return cachedSchema, nil
	}

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"viewer": &graphql.Field{
				Type:        userType,
				Description: "The currently authenticated user.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					githubClient, ok := getClientFromContext(p.Context)
					if !ok {
						return nil, nil
					}
					return githubClient.GetAuthenticatedUser()
				},
			},
			"user": &graphql.Field{
				Type:        userType,
				Description: "Lookup a user by login.",
				Args: graphql.FieldConfigArgument{
					"login": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.String),
						Description: "The user's login.",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					githubClient, ok := getClientFromContext(p.Context)
					if !ok {
						return nil, nil
					}
					login, ok := p.Args["login"].(string)
					if !ok {
						return nil, nil
					}
					return githubClient.GetUser(login)
				},
			},
			"repository": &graphql.Field{
				Type:        repositoryType,
				Description: "Lookup a repository by owner and name.",
				Args: graphql.FieldConfigArgument{
					"owner": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.String),
						Description: "The login field of the repository's owner.",
					},
					"name": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.String),
						Description: "The name of the repository.",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					githubClient, ok := getClientFromContext(p.Context)
					if !ok {
						return nil, nil
					}
					owner, ownerOk := p.Args["owner"].(string)
					name, nameOk := p.Args["name"].(string)
					if !ownerOk || !nameOk {
						return nil, nil
					}
					return githubClient.GetRepository(owner, name)
				},
			},
			"organization": &graphql.Field{
				Type:        organizationType,
				Description: "Lookup an organization by login.",
				Args: graphql.FieldConfigArgument{
					"login": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.String),
						Description: "The organization's login.",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					githubClient, ok := getClientFromContext(p.Context)
					if !ok {
						return nil, nil
					}
					login, ok := p.Args["login"].(string)
					if !ok {
						return nil, nil
					}
					return githubClient.GetOrganization(login)
				},
			},
		},
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
	})
	if err != nil {
		return graphql.Schema{}, err
	}

	cachedSchema = schema
	schemaInitialized = true
	return cachedSchema, nil
}

// NewSchema is deprecated. Use GetSchema instead.
// Keeping for backward compatibility during transition.
func NewSchema(githubClient *github.Client) (graphql.Schema, error) {
	return GetSchema()
}
