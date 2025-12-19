package tools

import (
	"context"
)

// ExampleTool implementa un tool MCP minimale
type ExampleTool struct{}

func (t *ExampleTool) Name() string {
	return "example_tool"
}

func (t *ExampleTool) Description() string {
	return "Un tool di esempio che dimostra la struttura base MCP"
}

func (t *ExampleTool) Execute(
	ctx context.Context,
	input map[string]any,
) (any, error) {
	msg, _ := input["message"].(string)

	return map[string]any{
		"echo": msg,
	}, nil
}
