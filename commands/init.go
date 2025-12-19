package commands

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"

	c "github.com/monkeymatt0/gomcp-cli/constants"
	"github.com/spf13/cobra"
)

type mainTemplateData struct {
	Module      string
	ProjectName string
}

var InitCommand = &cobra.Command{
	Use:   "gomcp init <project_name>",
	Short: "Init your MCP project with the given name",
	Long:  "Create an MCP project respectin a specfici scaffholding for a standardized building",
	Run:   Init,
}

func Init(cmd *cobra.Command, args []string) {

	// Populating template data
	data := mainTemplateData{
		Module:      args[0],
		ProjectName: args[0],
	}

	// Project's folder creation
	if err := os.Mkdir(args[0], 0755); err != nil {
		log.Fatalf("mkdir failed: %v", err)
	}

	// Setting new base path with the created folder
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("base path failed: %v", err)
		return
	}

	bp := strings.Join([]string{wd, args[0]}, "/")

	// Calling go mod init
	gmi := exec.Command("go", "mod", "init", args[0])
	gmi.Dir = bp
	if err := gmi.Run(); err != nil {
		log.Fatalf("go mod init failed: %v", err)
	}

	// Calling go get
	gi := exec.Command("go", "get", "github.com/modelcontextprotocol/go-sdk", args[0])
	gi.Dir = bp
	if err := gi.Run(); err != nil {
		log.Fatalf("go get failed: %v", err)
	}

	/*
	* Creation of the folder internal inside the project folder
	 */
	bpi := strings.Join([]string{bp, c.Internal}, "/")
	if err := os.Mkdir(bpi, 0755); err != nil {
		log.Fatalf("mkdir internal failed: %v", err)
	}

	/*
	* Creation of the following 2 folders:
	* - internal/tools
	* - internal/registry
	 */
	bpit := strings.Join([]string{bpi, c.Tools}, "/")
	if err := os.Mkdir(bpit, 0755); err != nil {
		log.Fatalf("mkdir tools failed: %v", err)
	}

	bpir := strings.Join([]string{bpi, c.Registry}, "/")
	if err := os.Mkdir(bpir, 0755); err != nil {
		log.Fatalf("mkdir registry failed: %v", err)
	}

	/*
	* Creation of the file:
	*  - main.go
	 */
	_src, e := os.Executable()
	if e != nil {
		log.Fatalf("getwd failed: %v", e)
	}
	_ssrc := strings.Split(_src, "/")
	_bp := strings.Join(_ssrc[:len(_ssrc)-1], "/")
	_src = _bp
	src := strings.Join([]string{_src, c.Template, c.Tmain}, "/")
	dst := strings.Join([]string{bp, c.Tmain}, "/")

	mf, err := os.Create(dst[:len(dst)-4])
	if err != nil {
		log.Fatalf("generate main failed: %v", err)
	}
	defer mf.Close()

	tpl, err := template.ParseFiles(src)
	if err != nil {
		log.Fatalf("failed to parse template files: %v", err)
	}

	err = tpl.Execute(mf, data)
	if err != nil {
		log.Fatalf("template interpolation failed: %v", err)
	}

	/*
	* Creation of the file:
	*  - internal/tools/example.go
	 */
	src = strings.Join([]string{_src, c.Template, c.Ttools}, "/")

	dst = strings.Join([]string{bp, c.Internal, c.Tools, c.Ttools[:len(c.Ttools)-4]}, "/")
	if err := copyFile(src, dst); err != nil {
		log.Fatalf("copy main failed: %v", err)
	}

	/*
	* Creation of the file:
	*  - internal/registry/registry.go
	 */
	src = strings.Join([]string{_src, c.Template, c.Tregistry}, "/")

	dst = strings.Join([]string{bp, c.Internal, c.Registry, c.Tregistry[:len(c.Tregistry)-4]}, "/")
	if err := copyFile(src, dst); err != nil {
		log.Fatalf("registry copy failed: %v", err)
	}

	/*
	* Creation of the file:
	*  - internal/registry/loader.go
	 */
	src = strings.Join([]string{_src, c.Template, c.Tloader}, "/")

	dst = strings.Join([]string{bp, c.Internal, c.Registry, c.Tloader[:len(c.Tloader)-4]}, "/")
	if err := copyFile(src, dst); err != nil {
		log.Fatalf("registry copy failed: %v", err)
	}
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, data, 0644)
}
