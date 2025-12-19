package tools

import "modelcontextprotocol/go-sdk/mcp"

func {{ .FunctionName }}() *mcp.Tool {
	return &mcp.Tool{
		Name: "{{ .ToolName }}",
		Description: "TODO",
		Handler: func(ctx *mcp.Context, input any) (any, error) {
			return nil, nil
		},
	}
}