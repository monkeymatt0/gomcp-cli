package registry

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// LoadTools aggiunge tutti i tool registrati al server MCP
func LoadTools(server *mcp.Server) {
	for _, tool := range GetAll() {
		server.AddTool(tool)
	}
}
