package github

import (
	"signalboard/internal/sources"
	"sync"
)

type GithubSource struct {
	mu sync.Mutex

	client *Client

	calendar ContributionCalendar
}

func NewGithubSource(client *Client) *GithubSource {
	return &GithubSource{
		client: client,
	}
}

func (s *GithubSource) Name() string {
	return "github"
}

func (s *GithubSource) Refresh() error {
	// Run method to retrieve github contributions!
	return nil
}

func (s *GithubSource) Endpoints() []sources.Endpoint {
	// Return slice with sources.Endpoints that will be created onto the subpath
	return nil
}
