package sources

import (
	"context"
	"net/http"
)

type Endpoint struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

type Source interface {
	Name() string
	Refresh(ctx context.Context) error
	Endpoints() []Endpoint
	// Content() []content.Content
}
