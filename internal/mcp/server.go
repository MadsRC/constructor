package mcp

import (
	"fmt"
	"log/slog"

	"github.com/MadsRC/constructor/internal/generator"
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

	mcpServer := server.NewMCPServer("constructor", opts.Version)

	return &Server{
		server:  mcpServer,
		options: &opts,
	}, nil
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
