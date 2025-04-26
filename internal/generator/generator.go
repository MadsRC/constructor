package generator

import (
	"log/slog"
)

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
	Logger     *slog.Logger
}

var defaultGeneratorOptions = generatorOptions{
	Logger:  slog.Default(),
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
