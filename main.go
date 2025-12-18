package main

import (
	"fmt"
	"os"

	c "github.com/monkeymatt0/gomcp-cli/commands"
)

func main() {

	var command = c.InitCommand
	if err := command.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
