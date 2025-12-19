package main

import (
	"fmt"
	"log"
	"os"

	c "github.com/monkeymatt0/gomcp-cli/commands"
	"github.com/spf13/cobra"
)

func main() {

	var command *cobra.Command

	if len(os.Args) < 2 {
		fmt.Println("Usage: gomcp <command>")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		command = c.InitCommand
	case "register":
		command = c.RegisterCommand
	default:
		fmt.Println("Unknown command: ", os.Args[1])
	}

	if err := command.Execute(); err != nil {
		log.Fatalf("Error while executing the command: %v", err)
	}

}
