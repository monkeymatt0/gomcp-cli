package commands

import (
	"github.com/spf13/cobra"
)

var RegisterCommand = &cobra.Command{
	Use:   "gomcp init <project_name>",
	Short: "Init your MCP project with the given name",
	Long:  "Create an MCP project respectin a specfici scaffholding for a standardized building",
	Run:   Generate,
}

func Generate(cmd *cobra.Command, args []string) {
	// 1. Create a new file in internal/tools folder using the passed name
	// 2. Append the tools on the MCP server
}
