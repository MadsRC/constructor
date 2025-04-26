package mcp

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/MadsRC/constructor/internal/generator"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type Server struct {
	server  *server.MCPServer
	options *serverOptions
}

// NewServer creates a new [Server].
func NewServer(options ...ServerOption) (*Server, error) {
	opts := defaultServerOptions
	for _, opt := range GlobalServerOptions {
		opt.apply(&opts)
	}
	for _, opt := range options {
		opt.apply(&opts)
	}

	// Add validation check
	if opts.Generator == nil {
		return nil, fmt.Errorf("must provide generator through WithServerGenerator option")
	}

	s := &Server{
		options: &opts,
	}

	mcpServer := server.NewMCPServer("constructor", opts.Version)
	mcpServer.AddTool(constructorTool, s.callConstructor)

	s.server = mcpServer
	return s, nil
}

type serverOptions struct {
	Logger    *slog.Logger
	Version   string
	Generator *generator.Generator // New required field
}

var defaultServerOptions = serverOptions{
	Logger:  slog.Default(),
	Version: "0.0.0",
}

// GlobalServerOptions is a list of [ServerOption]s that are applied to all [Server]s.
var GlobalServerOptions []ServerOption

// ServerOption is an option for configuring a [Server].
type ServerOption interface {
	apply(*serverOptions)
}

// funcServerOption is a [ServerOption] that calls a function.
// It is used to wrap a function, so it satisfies the [ServerOption] interface.
type funcServerOption struct {
	f func(*serverOptions)
}

func (fdo *funcServerOption) apply(opts *serverOptions) {
	fdo.f(opts)
}

func newFuncServerOption(f func(*serverOptions)) *funcServerOption {
	return &funcServerOption{
		f: f,
	}
}

// WithServerLogger returns a [ServerOption] that uses the provided logger.
func WithServerLogger(logger *slog.Logger) ServerOption {
	return newFuncServerOption(func(opts *serverOptions) {
		opts.Logger = logger
	})
}

// WithServerVersion returns a [ServerOption] that uses the provided version.
func WithServerVersion(version string) ServerOption {
	return newFuncServerOption(func(opts *serverOptions) {
		opts.Version = version
	})
}

// WithServerGenerator returns a [ServerOption] that provides a generator instance
func WithServerGenerator(generator *generator.Generator) ServerOption {
	return newFuncServerOption(func(opts *serverOptions) {
		opts.Generator = generator
	})
}

// ServeStdio starts the server and serves requests from stdin and writes responses to stdout.
func (s *Server) ServeStdio() error {
	return server.ServeStdio(s.server)
}

func (s *Server) callConstructor(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	pkgName := request.Params.Arguments["package"].(string)
	typeName := request.Params.Arguments["name"].(string)
	isTest := mcp.ParseBoolean(request, "test", false)
	output := request.Params.Arguments["output"].(string)

	err := s.options.Generator.Generate(pkgName, typeName, isTest, output)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("unable to generate constructor", err), nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{mcp.TextContent{
			Type: "text",
			Text: fmt.Sprintf("Generated constructor for %s in %s", typeName, pkgName),
		}},
	}, nil
}

var constructorTool = mcp.NewTool("constructor",
	mcp.WithDescription("Generate golang constructors using the functional options pattern"),
	mcp.WithString("name",
		mcp.Required(),
		mcp.Description("The name of the struct to generate"),
	),
	mcp.WithString("package",
		mcp.Required(),
		mcp.Description("The package to generate the struct in"),
	),
	mcp.WithString("output",
		mcp.Required(),
		mcp.Description("The output file to write the generated code to"),
	),
	mcp.WithBoolean("test",
		mcp.Description("Generate a test file instead of a source file"),
	),
)
