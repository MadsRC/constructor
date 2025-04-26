package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
	"github.com/yourusername/constructor/internal/generator"
)
var version string
var commit string
var date string

func main() {
	cmd := &cli.Command{
		Name:  "constructor",
		Usage: "A tool to generate constructor functions in the style of the functional options pattern for Go structs.",
		Version: fmt.Sprintf("%s (commit: %s, date: %s)", version, commit, date),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "name",
				Usage: "name of the thing to create a constructor for",
				Value: "Client",
			},
			&cli.StringFlag{
				Name:  "package",
				Usage: "name of the package the generated code should belong to",
				Value: "client",
			},
			&cli.StringFlag{
				Name:  "output",
				Usage: "output file for the generated code. If not provided, stdout will be used",
			},
			&cli.BoolFlag{
				Name:  "test",
				Usage: "output tests for the generated code, instead of the code itself. Uses the output flag to determine the output file",
			},
		},
		Action: mainAction,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}


func mainAction(ctx context.Context, cmd *cli.Command) error {
	name := cmd.String("name")
	if name == "" {
		return fmt.Errorf("provided \"name\" value is invalid: '%s' - Use '--name' to set it", name)
	}
	
	pkg := cmd.String("package")
	if pkg == "" {
		return fmt.Errorf("provided \"package\" value is invalid: '%s' - Use '--package' to set it", pkg)
	}

	// Create generator with default options
	gen, err := generator.NewGenerator()
	if err != nil {
		return fmt.Errorf("generator initialization failed: %w", err)
	}

	// Delegate generation to the package
	return gen.Generate(
		pkg,
		name,
		cmd.Bool("test"),
		cmd.String("output"),
	)
}
