package mcpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/blumsicle/bcli/cmd/bcli/create"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Server wraps the MCP server and bcli execution settings.
type Server struct {
	cfg    Config
	runner Runner
}

// Runner runs external commands for MCP tool handlers.
type Runner interface {
	Run(ctx context.Context, request CommandRequest) (CommandResult, error)
}

// CommandRequest describes an external command invocation.
type CommandRequest struct {
	Command    string
	Args       []string
	WorkingDir string
}

// CommandResult describes an external command result.
type CommandResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

// ExecRunner runs commands with os/exec.
type ExecRunner struct{}

// Run executes an external command and captures stdout and stderr.
func (ExecRunner) Run(ctx context.Context, request CommandRequest) (CommandResult, error) {
	cmd := exec.CommandContext(ctx, request.Command, request.Args...)
	cmd.Dir = request.WorkingDir
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			return CommandResult{}, err
		}
	}

	return CommandResult{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: exitCode,
	}, nil
}

// New constructs a bcli MCP server wrapper.
func New(cfg Config, runner Runner) *Server {
	if runner == nil {
		runner = ExecRunner{}
	}

	return &Server{
		cfg:    cfg,
		runner: runner,
	}
}

// MCP constructs the MCP server with all bcli tools registered.
func (s *Server) MCP() *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "bcli-mcp",
		Version: "dev",
	}, nil)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name: "create_go_cli_project",
			Description: "Create a new Go command-line tool using the installed bcli generator. " +
				"Use this when the user asks to create, scaffold, or generate a new Go CLI project.",
		},
		s.createGoCLIProject,
	)

	return server
}

// RunStdio runs the MCP server over stdio.
func (s *Server) RunStdio(ctx context.Context) error {
	return s.MCP().Run(ctx, &mcp.StdioTransport{})
}

type createProjectInput struct {
	Name             string `json:"name"                         jsonschema:"Name of the new Go CLI project"`
	Description      string `json:"description"                  jsonschema:"Description for the generated CLI"`
	WorkingDirectory string `json:"working_directory,omitempty"  jsonschema:"Directory where bcli should run"`
	RootPath         string `json:"root_path,omitempty"          jsonschema:"Directory to create the new project under"`
	ProjectDirPrefix string `json:"project_dir_prefix,omitempty" jsonschema:"Prefix to prepend to the generated project directory name"`
	GitLocation      string `json:"git_location,omitempty"       jsonschema:"Git host and owner prefix for the generated module path"`
	BCLIConfigFile   string `json:"bcli_config_file,omitempty"   jsonschema:"Path to the bcli config file"`
	InPlace          bool   `json:"inplace,omitempty"            jsonschema:"Create the project in the working directory"`
	SkipGoGetUpdate  bool   `json:"skip_go_get_update,omitempty" jsonschema:"Skip the go get -u ./... post step"`
	SkipGoModTidy    bool   `json:"skip_go_mod_tidy,omitempty"   jsonschema:"Skip the go mod tidy post step"`
	SkipGitInit      bool   `json:"skip_git_init,omitempty"      jsonschema:"Skip the git init post step"`
	SkipGitCommit    bool   `json:"skip_git_commit,omitempty"    jsonschema:"Skip the git commit post step"`
}

type createProjectOutput struct {
	create.CreateResult
	Command    []string `json:"command"`
	WorkingDir string   `json:"working_directory,omitempty"`
	Stderr     string   `json:"stderr,omitempty"`
}

func (s *Server) createGoCLIProject(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input createProjectInput,
) (*mcp.CallToolResult, createProjectOutput, error) {
	output, err := s.CreateGoCLIProject(ctx, input)
	if err != nil {
		return nil, createProjectOutput{}, err
	}

	text := fmt.Sprintf("Created Go CLI project %q at %s", output.Project, output.TargetPath)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: text},
		},
	}, output, nil
}

// CreateGoCLIProject runs bcli create and returns parsed JSON output.
func (s *Server) CreateGoCLIProject(
	ctx context.Context,
	input createProjectInput,
) (createProjectOutput, error) {
	if input.Name == "" {
		return createProjectOutput{}, fmt.Errorf("name is required")
	}
	if input.Description == "" {
		return createProjectOutput{}, fmt.Errorf("description is required")
	}

	args := createArgs(input)
	request := CommandRequest{
		Command:    s.cfg.BCLICommand,
		Args:       args,
		WorkingDir: input.WorkingDirectory,
	}

	runCtx, cancel := context.WithTimeout(ctx, s.cfg.Timeout)
	defer cancel()

	result, err := s.runner.Run(runCtx, request)
	if err != nil {
		return createProjectOutput{}, fmt.Errorf("run bcli: %w", err)
	}
	if result.ExitCode != 0 {
		return createProjectOutput{}, fmt.Errorf(
			"bcli exited with code %d: %s",
			result.ExitCode,
			strings.TrimSpace(result.Stderr),
		)
	}

	var output createProjectOutput
	if err := json.Unmarshal([]byte(result.Stdout), &output); err != nil {
		return createProjectOutput{}, fmt.Errorf("parse bcli json output: %w", err)
	}
	output.Command = append([]string{s.cfg.BCLICommand}, args...)
	output.WorkingDir = input.WorkingDirectory
	output.Stderr = result.Stderr

	return output, nil
}

func createArgs(input createProjectInput) []string {
	args := []string{}
	if input.BCLIConfigFile != "" {
		args = append(args, "--config-file", input.BCLIConfigFile)
	}

	args = append(args, "create", "--json")
	if input.RootPath != "" && !input.InPlace {
		args = append(args, "--root-path", input.RootPath)
	}
	if input.ProjectDirPrefix != "" && !input.InPlace {
		args = append(args, "--project-dir-prefix", input.ProjectDirPrefix)
	}
	if input.GitLocation != "" {
		args = append(args, "--git-location", input.GitLocation)
	}
	if input.InPlace {
		args = append(args, "--inplace")
	}
	if input.SkipGoGetUpdate {
		args = append(args, "--no-go-get-update")
	}
	if input.SkipGoModTidy {
		args = append(args, "--no-go-mod-tidy")
	}
	if input.SkipGitInit {
		args = append(args, "--no-git-init")
	}
	if input.SkipGitCommit {
		args = append(args, "--no-git-commit")
	}

	args = append(args, input.Name, input.Description)
	return args
}
