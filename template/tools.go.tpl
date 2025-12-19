package tools

import "modelcontextprotocol/go-sdk/mcp"

func MyTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "myTool",
		Description: "TODO",
		Handler: func(ctx *mcp.Context, input any) (any, error) {
			return nil, nil
		},
	}
}