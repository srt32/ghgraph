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
		},
	})

	return graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
	})
}
