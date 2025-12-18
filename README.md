# gomcp-cli

CLI to develop MCP server following this structure:

```
mcp-hello-fs/
├── go.mod
├── main.go               
└── internal/
    ├── tools/            
    │   └── example.go
    └── registry/         
        ├── registry.go   
        └── loader.go     
```  

## tools
Folder containing all your tools, basically the implementation of the function that the tool will use

## registry
In this folder we have 2 files:

- registry: register the functions
- loader: will add all the tools to the MCP server


## COMMANDS

### init (WIP)
```bash
gomcp init <project_name>
```

This command will init you MCP project with the basic command go mod init and also create the scaffholding for the project.

### register (TODO)
```bash
gomcp register <tool_name>
```

This command will updates:

- register.go -> Here will register the function using a map
- loader.go -> Here will add to MCP server
- tools -> will add a new file: tool_name.go with a boilerplate


WIP