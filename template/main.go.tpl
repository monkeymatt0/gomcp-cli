package main

import (
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	r "{{ .Module }}/internal/registry"
)

func main() {

	// Your server MCP: change the name
	server := mcp.NewServer(&mcp.Implementation{Name: "{{ .ProjectName }}", Version: "v1.0.0"}, nil)

	r.LoadTools(server)

	// HTTP mode
	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, nil)

	if err := http.ListenAndServe("localhost:3000", handler); err != nil {
		log.Fatalf("Failed to start up server: %v", err)
	}
}
