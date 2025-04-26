package mcp

import (
	"log/slog"

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

	mcpServer := server.NewMCPServer("constructor", opts.Version)

	return &Server{
		server:  mcpServer,
		options: &opts,
	}, nil
}

type serverOptions struct {
	Logger  *slog.Logger
	Version string
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

// ServeStdio starts the server and serves requests from stdin and writes responses to stdout.
func (s *Server) ServeStdio() error {
	return server.ServeStdio(s.server)
}
