package executor

import (
	"context"
	"signalboard/internal/plugin"
	"testing"
)

func TestExecutor(t *testing.T) {
	ctx := context.Background()

	executor := New()

	req := plugin.Request{
		Plugin: "hello",
		Widget: "greeting",
	}

	resp, err := executor.Execute(
		ctx,
		"/tmp/signalboard-hello",
		req,
	)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Type != "text/html" {
		t.Fatalf("unexpected type: %s", resp.Type)
	}

	if resp.Content != "<h1>Hello from Signalboard!</h1>" {
		t.Fatalf("unexpected content: %s", resp.Content)
	}
}
