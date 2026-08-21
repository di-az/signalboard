package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"signalboard/internal/plugin"
)

type Executor struct {
}

func NewExecutor() *Executor {
	return &Executor{}
}

func (e *Executor) Execute(
	ctx context.Context,
	pluginPath string,
	req plugin.Request,
) (*plugin.Response, error) {
	input, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal plugin request: %w", err)
	}

	cmd := exec.CommandContext(ctx, pluginPath)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("create plugin stdin: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("create plugin stdout: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start plugin: %w", err)
	}

	if _, err := stdin.Write(input); err != nil {
		return nil, fmt.Errorf("write plugin request: %w", err)
	}

	if err := stdin.Close(); err != nil {
		return nil, fmt.Errorf("close plugin stdin: %w", err)
	}

	var resp plugin.Response
	if err := json.NewDecoder(stdout).Decode(&resp); err != nil {
		return nil, fmt.Errorf("decode plugin response: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("plugin execution: %w", err)
	}

	return &resp, nil
}
