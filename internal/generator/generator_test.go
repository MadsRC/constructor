package generator

import (
	"io"
	"log/slog"
	"testing"
	
	"github.com/stretchr/testify/require"
)

func TestNewGenerator(t *testing.T) {
	discardLogger := slog.New(slog.NewTextHandler(io.Discard, nil))
	tests := []struct {
		name    string
		options []GeneratorOption
		want    *Generator
		wantErr bool
	}{
		{
			name:    "Create with default logger",
			options: []GeneratorOption{},
			want: &Generator{
				options: &generatorOptions{
					Logger: slog.Default(),
				},
			},
			wantErr: false,
		},
		{
			name:    "Create with custom logger",
			options: []GeneratorOption{WithGeneratorLogger(discardLogger)},
			want: &Generator{
				options: &generatorOptions{
					Logger: discardLogger,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGenerator(tt.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGenerator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.options.Logger != tt.want.options.Logger {
				t.Errorf("NewGenerator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTitleFunc(t *testing.T) {
	g, err := NewGenerator()
	require.NoError(t, err)

	tests := []struct {
		input    string
		expected string
	}{
		{"client", "Client"},
		{"fiskePinde", "FiskePinde"},
		{"a", "A"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := g.options.TitleFunc(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func BenchmarkTitleFunc(b *testing.B) {
	g, _ := NewGenerator()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		g.options.TitleFunc("fiskePinde")
	}
}

func TestNewGenerator_GlobalOptions(t *testing.T) {
	tests := []struct {
		name        string
		options     []GeneratorOption
		inputLogger *slog.Logger
	}{
		{
			name:        "Global options are applied",
			options:     []GeneratorOption{},
			inputLogger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GlobalGeneratorOptions = []GeneratorOption{
				WithGeneratorLogger(tt.inputLogger),
			}
			got1, _ := NewGenerator(tt.options...)
			got2, _ := NewGenerator(tt.options...)
			if got1.options.Logger != tt.inputLogger {
				t.Errorf("NewGenerator() = %v, want %v", got1, tt.inputLogger)
			}
			if got2.options.Logger != tt.inputLogger {
				t.Errorf("NewGenerator() = %v, want %v", got2, tt.inputLogger)
			}
			if got1.options.Logger != got2.options.Logger {
				t.Errorf("NewGenerator() = %v, want %v", got1, got2)
			}
			GlobalGeneratorOptions = []GeneratorOption{}
			got3, _ := NewGenerator(tt.options...)
			if got3.options.Logger == tt.inputLogger {
				t.Errorf("NewGenerator() = %v, want %v", got3, slog.Default())
			}
		})
	}
}
