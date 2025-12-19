package commands

import (
	"log"
	"os"
	"strings"
	"text/template"

	c "github.com/monkeymatt0/gomcp-cli/constants"
	"github.com/spf13/cobra"
)

type data struct {
	FunctionName string
	ToolName     string
}

var RegisterCommand = &cobra.Command{
	Use:   "gomcp init <project_name>",
	Short: "Init your MCP project with the given name",
	Long:  "Create an MCP project respectin a specfici scaffholding for a standardized building",
	Run:   Generate,
}

func Generate(cmd *cobra.Command, args []string) {

	data := data{
		FunctionName: strings.ToUpper(args[1][:1]) + strings.ToLower(args[1][1:]),
		ToolName:     args[1],
	}
	// 1. Create a new file in internal/tools folder using the passed name
	_bp, err := os.Getwd()
	if err != nil {
		log.Fatalf("Getwd failed: %v", err)
	}

	els, err2 := os.ReadDir(_bp)
	if err2 != nil {
		log.Fatalf("ReadDir failed: %v", err2)
	}

	pf := ""

	for _, el := range els {
		if strings.Contains(el.Name(), "mcp_") {
			pf = el.Name()
		}
	}

	bp := strings.Join([]string{_bp, pf}, "/")

	// 1.2 Generating the file tool
	src := strings.Join([]string{_bp, c.Template, c.Ttcustom}, "/")
	dst := strings.Join([]string{bp, c.Internal, c.Tools, strings.Join([]string{args[1], "go"}, ".")}, "/")

	tf, err := os.Create(dst)
	if err != nil {
		log.Fatalf("file copy tool registration failed: %v", err)
	}

	tplf, err2 := template.ParseFiles(src)
	if err2 != nil {
		log.Fatalf("parse files failed: %e", err2)
	}

	if err := tplf.Execute(tf, data); err != nil {
		log.Fatalf("failed to execute template: %e", err)
	}

	// 2. Append the tools on the MCP server
	fp := strings.Join([]string{bp, c.Internal, c.Registry, c.Tregistry[:len(c.Tregistry)-4]}, "/")
	f, err := os.ReadFile(fp)
	if err != nil {
		log.Fatalf("read file failed: %v", err)
	}

	content := string(f)
	marker := "// gomcp:tools"

	if !strings.Contains(content, marker) {
		log.Fatalf("marker missing, want: %v", marker)
	}

	line := "\nserver.AddTool(tools." + strings.ToUpper(args[1][:1]) + strings.ToLower(args[1][1:]) + "(), nil)\n"
	if strings.Contains(content, line) {
		log.Fatal("tool already exists")
	}

	content = strings.Replace(content, marker, line+marker, 1)

	if err := os.WriteFile(fp, []byte(content), 0644); err != nil {
		log.Fatal("Not able to wirte in the file")
	}
}
