package resolver_test

import (
	"context"
	"testing"

	"github.com/google/go-github/v28/github"
	"github.com/tj/assert"
	"github.com/tj/go/env"
	"golang.org/x/oauth2"

	"github.com/tj/gobinaries"
	"github.com/tj/gobinaries/resolver"
)

// newGitHubResolver returns a new GitHub resolver.
func newGitHubResolver() gobinaries.Resolver {
	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: env.Get("GITHUB_TOKEN"),
		},
	)

	return &resolver.GitHub{
		Client: github.NewClient(oauth2.NewClient(ctx, ts)),
	}
}

// Test resolver.
func TestGitHub_Resolve(t *testing.T) {
	r := newGitHubResolver()

	t.Run("exact match", func(t *testing.T) {
		repo := resolver.Repository{
			Owner:   "tj",
			Project: "d3-bar",
			Version: "v1.8.0",
		}
		v, err := r.Resolve(repo)
		assert.NoError(t, err)
		assert.Equal(t, "v1.8.0", v)
	})

	t.Run("exact match without leading v", func(t *testing.T) {
		repo := resolver.Repository{
			Owner:   "tj",
			Project: "d3-bar",
			Version: "1.8.0",
		}
		v, err := r.Resolve(repo)
		assert.NoError(t, err)
		assert.Equal(t, "v1.8.0", v)
	})

	t.Run("major wildcard match", func(t *testing.T) {
		repo := resolver.Repository{
			Owner:   "tj",
			Project: "d3-bar",
			Version: "1.x",
		}
		v, err := r.Resolve(repo)
		assert.NoError(t, err)
		assert.Equal(t, "v1.8.0", v)
	})

	t.Run("minor wildcard match", func(t *testing.T) {
		repo := resolver.Repository{
			Owner:   "tj",
			Project: "d3-bar",
			Version: "1.6.x",
		}
		v, err := r.Resolve(repo)
		assert.NoError(t, err)
		assert.Equal(t, "v1.6.0", v)
	})

	t.Run("minor match", func(t *testing.T) {
		repo := resolver.Repository{
			Owner:   "tj",
			Project: "d3-bar",
			Version: "1.6",
		}
		v, err := r.Resolve(repo)
		assert.NoError(t, err)
		assert.Equal(t, "v1.6.0", v)
	})

	t.Run("master", func(t *testing.T) {
		repo := resolver.Repository{
			Owner:   "tj",
			Project: "d3-bar",
			Version: "master",
		}
		v, err := r.Resolve(repo)
		assert.NoError(t, err)
		assert.Equal(t, "v1.8.0", v)
	})
}
