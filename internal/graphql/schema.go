package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/srt32/ghgraph/internal/github"
)

// userType represents the GraphQL User type matching GitHub's schema
var userType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "User",
	Description: "A user is an individual's account on GitHub that owns repositories and can make new content.",
	Fields: graphql.Fields{
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
	},
})

// repositoryType represents the GraphQL Repository type matching GitHub's schema
var repositoryType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Repository",
	Description: "A repository contains the content for a project.",
	Fields: graphql.Fields{
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
	},
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

// NewSchema creates a new GraphQL schema
func NewSchema(githubClient *github.Client) (graphql.Schema, error) {
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"viewer": &graphql.Field{
				Type:        userType,
				Description: "The currently authenticated user.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
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
					login, ok := p.Args["login"].(string)
					if !ok {
						return nil, nil
					}
					return githubClient.GetOrganization(login)
				},
			},
		},
	})

	return graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
	})
}
