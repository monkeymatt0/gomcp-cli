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
	_bp, err := os.Executable()
	if err != nil {
		log.Fatalf("Getwd failed: %v", err)
	}

	t := strings.Split(_bp, "/")
	_bp = strings.Join(t[:len(t)-1], "/")

	// 1.2 Generating the file tool
	src := strings.Join([]string{_bp, c.Template, c.Ttcustom}, "/")

	// project dst
	_wd, err2 := os.Getwd()
	if err2 != nil {
		log.Fatalf("Getwd failed: %v", err2)
	}

	dirs, err3 := os.ReadDir(_wd)
	if err3 != nil {
		log.Fatalf("Failed to read a directory: %v", err3)
	}

	pn := ""
	for _, dir := range dirs {
		if strings.Contains(dir.Name(), "mcp_") {
			pn = dir.Name()
		}
	}

	dst := strings.Join([]string{_wd, pn, c.Internal, c.Tools, strings.Join([]string{args[1], "go"}, ".")}, "/")

	tf, err := os.Create(dst)
	if err != nil {
		log.Fatalf("file copy tool registration failed: %v", err)
	}

	tplf, err4 := template.ParseFiles(src)
	if err4 != nil {
		log.Fatalf("parse files failed: %e", err4)
	}

	if err5 := tplf.Execute(tf, data); err5 != nil {
		log.Fatalf("failed to execute template: %e", err5)
	}

	// 2. Append the tools on the MCP server
	fp := strings.Join([]string{_wd, pn, c.Internal, c.Registry, c.Tregistry[:len(c.Tregistry)-4]}, "/")
	f, err6 := os.ReadFile(fp)
	if err6 != nil {
		log.Fatalf("read file failed: %v", err6)
	}

	content := string(f)
	marker := "\t// gomcp:tools"

	if !strings.Contains(content, marker) {
		log.Fatalf("marker missing, want: %v", marker)
	}

	line := "\tserver.AddTool(tools." + strings.ToUpper(args[1][:1]) + strings.ToLower(args[1][1:]) + "(), nil)\n"
	if strings.Contains(content, line) {
		log.Fatal("tool already exists")
	}

	content = strings.Replace(content, marker, line+marker, 1)

	if err := os.WriteFile(fp, []byte(content), 0644); err != nil {
		log.Fatal("Not able to wirte in the file")
	}
}
