package commands

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

// getProjectRoot trova la root del progetto in modo affidabile
func getProjectRoot() (string, error) {
	// Prova prima con os.Executable() (funziona quando compilato)
	execPath, err := os.Executable()
	if err == nil {
		// Se il percorso contiene "go-build", siamo in modalità go run
		if !strings.Contains(execPath, "go-build") {
			// Eseguibile compilato: risali alla directory del binario
			execDir := filepath.Dir(execPath)
			// Verifica se esiste la cartella template (siamo nella root)
			if _, err := os.Stat(filepath.Join(execDir, c.Template)); err == nil {
				return execDir, nil
			}
			// Altrimenti risali di un livello (se il binario è in una subdirectory)
			parent := filepath.Dir(execDir)
			if _, err := os.Stat(filepath.Join(parent, c.Template)); err == nil {
				return parent, nil
			}
		}
	}

	// Fallback: usa runtime.Caller per trovare il percorso del file sorgente
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", os.ErrNotExist
	}

	// Risali dalla directory del file corrente alla root del progetto
	// Il file è in commands/init.go, quindi risali di 1 livello
	dir := filepath.Dir(filename)
	projectRoot := filepath.Dir(dir) // Risali da commands/ a root

	// Verifica che esista la cartella template
	if _, err := os.Stat(filepath.Join(projectRoot, c.Template)); err != nil {
		return "", err
	}

	return projectRoot, nil
}

func Init(cmd *cobra.Command, args []string) {

	// Populating template data
	folder_name := "mcp_" + args[1]

	data := mainTemplateData{
		Module:      folder_name,
		ProjectName: args[1],
	}

	// Project's folder creation
	if err := os.Mkdir(folder_name, 0755); err != nil {
		log.Fatalf("mkdir failed: %v", err)
	}

	// Setting new base path with the created folder
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("base path failed: %v", err)
		return
	}

	bp := filepath.Join(wd, folder_name)

	// Calling go mod init
	gmi := exec.Command("go", "mod", "init", folder_name)
	gmi.Dir = bp
	if err := gmi.Run(); err != nil {
		log.Fatalf("go mod init failed: %v", err)
	}

	// Calling go get
	gi := exec.Command("go", "get", "github.com/modelcontextprotocol/go-sdk", folder_name)
	gi.Dir = bp
	if err := gi.Run(); err != nil {
		log.Fatalf("go get failed: %v", err)
	}

	/*
	* Creation of the folder internal inside the project folder
	 */
	bpi := filepath.Join(bp, c.Internal)
	if err := os.Mkdir(bpi, 0755); err != nil {
		log.Fatalf("mkdir internal failed: %v", err)
	}

	/*
	* Creation of the following 2 folders:
	* - internal/tools
	* - internal/registry
	 */
	bpit := filepath.Join(bpi, c.Tools)
	if err := os.Mkdir(bpit, 0755); err != nil {
		log.Fatalf("mkdir tools failed: %v", err)
	}

	bpir := filepath.Join(bpi, c.Registry)
	if err := os.Mkdir(bpir, 0755); err != nil {
		log.Fatalf("mkdir registry failed: %v", err)
	}

	/*
	* Get project root directory (where templates are located)
	 */
	projectRoot, err := getProjectRoot()
	if err != nil {
		log.Fatalf("failed to find project root: %v", err)
	}

	/*
	* Creation of the file:
	*  - main.go
	 */
	src := filepath.Join(projectRoot, c.Template, c.Tmain)
	dst := filepath.Join(bp, c.Tmain)

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
	src = filepath.Join(projectRoot, c.Template, c.Ttools)

	dst = filepath.Join(bp, c.Internal, c.Tools, c.Ttools[:len(c.Ttools)-4])
	if err := copyFile(src, dst); err != nil {
		log.Fatalf("registry copy failed: %v", err)
	}

	/*
	* Creation of the file:
	*  - internal/registry/registry.go
	 */
	src = filepath.Join(projectRoot, c.Template, c.Tregistry)

	dst = filepath.Join(bp, c.Internal, c.Registry, c.Tregistry[:len(c.Tregistry)-4])

	tf, err := os.Create(dst)
	if err != nil {
		log.Fatalf("registry creation failed: %v", err)
	}

	defer tf.Close()

	ttpl, err := template.ParseFiles(src)
	if err != nil {
		log.Fatalf("parse registry template failed: %v", err)
	}

	if err := ttpl.Execute(tf, data); err != nil {
		log.Fatalf("Failed to execute the template: %v", err)
	}
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, data, 0644)
}
