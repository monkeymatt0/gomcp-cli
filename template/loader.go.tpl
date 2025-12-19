package registry

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// LoadTools add all the registered tools to the server
func LoadTools(server *mcp.Server) {
	for _, tool := range GetAll() {
		server.AddTool(tool)
	}
}
