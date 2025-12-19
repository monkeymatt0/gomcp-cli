package registry

import (
	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

var tools = make(map[string]mcp.Tool)

// Register an mcp tool
func Register(tool mcp.Tool) {
	tools[tool.Name()] = tool
}

// GetAll return all the mcp tools in tools
func GetAll() []mcp.Tool {
	result := make([]mcp.Tool, 0, len(tools))
	for _, t := range tools {
		result = append(result, t)
	}
	return result
}
