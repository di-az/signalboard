package executor

import (
	"context"
	"fmt"
	"signalboard/internal/plugin"
)

type Executor struct {
}

func New() *Executor {
	return &Executor{}
}

func (e *Executor) Execute(
	ctx context.Context,
	pluginPath string,
	req plugin.Request,
) (*plugin.Response, error) {
	// TODO: serialize request and execute plugin
	return nil, fmt.Errorf("plugin execution not implemented")
}
