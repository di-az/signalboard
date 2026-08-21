package executor

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"signalboard/internal/plugin"
	"testing"
)

const HelloPluginName = "hello"

func buildPlugin(t *testing.T) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), HelloPluginName)

	cmd := exec.Command(
		"go",
		"build",
		"-o",
		path,
		fmt.Sprintf("./../../plugins/%s", HelloPluginName),
	)

	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build plugin: %v\n%s", err, output)
	}

	return path
}

func TestExecute(t *testing.T) {
	pluginPath := buildPlugin(t)

	executor := NewExecutor()

	req := plugin.Request{
		Plugin: "hello",
		Widget: "greeting",
	}

	resp, err := executor.Execute(
		t.Context(),
		pluginPath,
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

func TestExecutePluginNotFound(t *testing.T) {
	executor := NewExecutor()

	req := plugin.Request{
		Plugin: "hello",
		Widget: "greeting",
	}

	_, err := executor.Execute(
		t.Context(),
		"/does/not/exist",
		req,
	)

	if err == nil {
		t.Fatalf("expected error, got=%s", err)
	}
}
