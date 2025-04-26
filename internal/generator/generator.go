package generator

import (
	"embed"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

//go:embed templates/*
var embeddedTemplates embed.FS

type Generator struct {
    options *generatorOptions
}

// NewGenerator creates a new [Generator].
func NewGenerator(options ...GeneratorOption) (*Generator, error) {
	opts := defaultGeneratorOptions
	for _, opt := range GlobalGeneratorOptions {
		opt.apply(&opts)
	}
	for _, opt := range options {
		opt.apply(&opts)
	}

	return &Generator{
		options: &opts,
	}, nil
}

type generatorOptions struct {
	Logger              *slog.Logger
	TitleFunc           func(string) string
	LowerFirstLetterFunc func(string) string
}

var defaultGeneratorOptions = generatorOptions{
	Logger:  slog.Default(),
	TitleFunc: func(s string) string {
		return cases.Title(language.English).String(s[0:1]) + s[1:]
	},
	LowerFirstLetterFunc: func(s string) string {
		if len(s) < 1 {
			return strings.ToLower(s)
		}
		return strings.ToLower(s[:1]) + s[1:]
	},
}

// GlobalGeneratorOptions is a list of [GeneratorOption]s that are applied to all [Generator]s.
var GlobalGeneratorOptions []GeneratorOption

// GeneratorOption is an option for configuring a [Generator].
type GeneratorOption interface {
	apply(*generatorOptions)
}

// funcGeneratorOption is a [GeneratorOption] that calls a function.
// It is used to wrap a function, so it satisfies the [GeneratorOption] interface.
type funcGeneratorOption struct {
	f func(* generatorOptions)
}

func (fdo *funcGeneratorOption) apply(opts *generatorOptions) {
	fdo.f(opts)
}

func newFuncGeneratorOption(f func(*generatorOptions)) *funcGeneratorOption {
	return &funcGeneratorOption{
		f: f,
	}
}

// WithGeneratorLogger returns a [GeneratorOption] that uses the provided logger.
func WithGeneratorLogger(logger *slog.Logger) GeneratorOption {
	return newFuncGeneratorOption(func(opts *generatorOptions) {
		opts.Logger = logger
	})
}

// WithTitleFunc returns a [GeneratorOption] that uses the provided title function.
func WithTitleFunc(fn func(string) string) GeneratorOption {
	return newFuncGeneratorOption(func(opts *generatorOptions) {
		opts.TitleFunc = fn
	})
}

// WithLowerFirstLetterFunc returns a [GeneratorOption] that uses the provided lower first letter function.
func WithLowerFirstLetterFunc(fn func(string) string) GeneratorOption {
	return newFuncGeneratorOption(func(opts *generatorOptions) {
		opts.LowerFirstLetterFunc = fn
	})
}

// Generate generates code based on the provided parameters.
func (g *Generator) Generate(pkgName, typeName string, isTest bool, output string) error {
	tmplName := "main.go.tmpl"
	if isTest {
		tmplName = "main_test.go.tmpl"
	}

	tmplContent, err := embeddedTemplates.ReadFile("templates/" + tmplName)
	if err != nil {
		return fmt.Errorf("failed to read template: %w", err)
	}

	funcMap := template.FuncMap{
		"title":              g.options.TitleFunc,
		"lower_first_letter": g.options.LowerFirstLetterFunc,
	}

	tmpl, err := template.New("constructor").Funcs(funcMap).Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("template parsing failed: %w", err)
	}

	// Handle output destination
	var w io.Writer = os.Stdout
	if output != "" && output != "-" {
		if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
		f, err := os.Create(output)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer f.Close()
		w = f
	}

	return tmpl.Execute(w, struct {
		PackageName string
		Name        string
	}{
		PackageName: pkgName,
		Name:        typeName,
	})
}
