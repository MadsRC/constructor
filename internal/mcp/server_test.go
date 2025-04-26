package mcp

import (
	"io"
	"log/slog"
	"testing"

	"github.com/MadsRC/constructor/internal/generator"
)

func TestNewServer(t *testing.T) {
	discardLogger := slog.New(slog.NewTextHandler(io.Discard, nil))
	gen, _ := generator.NewGenerator() // Create generator instance
	tests := []struct {
		name    string
		options []ServerOption
		want    *Server
		wantErr bool
	}{
		{
			name:    "Create with default logger",
			options: []ServerOption{WithServerGenerator(gen)},
			want: &Server{
				options: &serverOptions{
					Logger:    slog.Default(),
					Generator: gen,
				},
			},
			wantErr: false,
		},
		{
			name:    "Create with custom logger",
			options: []ServerOption{WithServerGenerator(gen), WithServerLogger(discardLogger)},
			want: &Server{
				options: &serverOptions{
					Logger:    discardLogger,
					Generator: gen,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewServer(tt.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.options.Logger != tt.want.options.Logger {
				t.Errorf("NewServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewServer_GlobalOptions(t *testing.T) {
	gen, _ := generator.NewGenerator()
	tests := []struct {
		name        string
		options     []ServerOption
		inputLogger *slog.Logger
	}{
		{
			name:        "Global options are applied",
			options:     []ServerOption{WithServerGenerator(gen)},
			inputLogger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GlobalServerOptions = []ServerOption{
				WithServerLogger(tt.inputLogger),
				WithServerGenerator(gen), // Add generator to global options
			}
			got1, _ := NewServer(tt.options...)
			got2, _ := NewServer(tt.options...)
			if got1.options.Logger != tt.inputLogger {
				t.Errorf("NewServer() = %v, want %v", got1, tt.inputLogger)
			}
			if got2.options.Logger != tt.inputLogger {
				t.Errorf("NewServer() = %v, want %v", got2, tt.inputLogger)
			}
			if got1.options.Logger != got2.options.Logger {
				t.Errorf("NewServer() = %v, want %v", got1, got2)
			}
			GlobalServerOptions = []ServerOption{}
			got3, _ := NewServer(tt.options...)
			if got3.options.Logger == tt.inputLogger {
				t.Errorf("NewServer() = %v, want %v", got3, slog.Default())
			}
		})
	}
}
