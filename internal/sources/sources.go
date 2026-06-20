package sources

import (
	"context"
	// "signalboard/internal/content"
)

type Source interface {
	Name() string
	Refresh(ctx context.Context) error
	// Content() []content.Content
}
