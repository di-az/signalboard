package plugin

import "encoding/json"

// Request sent by Signalboard for executing a widget
type Request struct {
	Plugin string          `json:"plugin"`
	Widget string          `json:"widget"`
	Config json.RawMessage `json:"config,omitempty"`
}

// Response sent by a plugin to Signalboard
type Response struct {
	Type    string  `json:"type"`
	Content string  `json:"content"`
	Assets  []Asset `json:"assets,omitempty"`
}

// Asset describes a resource that can be referenced by the widget.
type Asset struct {
	ID          string    `json:"id"`
	Type        AssetType `json:"type"`
	ContentType string    `json:"content_type"`
	Source      string    `json:"source"`
}

// AssetType identifies where an asset is provided from.
type AssetType string

const (
	AssetStatic   AssetType = "static"
	AssetDynamic  AssetType = "dynamic"
	AssetExternal AssetType = "external"
)
