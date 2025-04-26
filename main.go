package main

import (
	"context"
	"embed"
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/urfave/cli/v3"
)

//go:embed templates/*
var templates embed.FS
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

type tmplInput struct {
	PackageName string
	Name        string
}

func titleFunc(s string) string {
	return cases.Title(language.English).String(s[0:1]) + s[1:]
}

func lowerFirstLetterFunc(s string) string {
	if len(s) < 1 {
		return cases.Lower(language.English).String(s)
	}
	return cases.Lower(language.English).String(s[:1]) + s[1:]
}

func mainAction(ctx context.Context, cmd *cli.Command) error {
	if cmd.String("name") == "" {
		return fmt.Errorf("provided \"name\" value is invalid: '%s' - Use '--name' to set it", cmd.String("name"))
	}
	if cmd.String("package") == "" {
		return fmt.Errorf("provided \"package\" value is invalid: '%s' - Use '--package' to set it", cmd.String("package"))
	}
	funcMap := template.FuncMap{
		"title":              titleFunc,
		"lower_first_letter": lowerFirstLetterFunc,
	}

	tmpl, err := determineTemplate(cmd, funcMap)
	if err != nil {
		return fmt.Errorf("error determining template: %w", err)
	}

	input := tmplInput{
		PackageName: cmd.String("package"),
		Name:        cmd.String("name"),
	}

	output, err := determineOutput(cmd)
	if err != nil {
		return fmt.Errorf("error determining output destination: %w", err)
	}

	err = tmpl.Execute(output, input)
	if err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	return nil
}

func determineTemplate(cmd *cli.Command, funcMap template.FuncMap) (*template.Template, error) {
	var tmplContent []byte
	var err error
	if cmd.Bool("test") {
		tmplContent, err = templates.ReadFile("templates/main_test.go.tmpl")
	} else {
		tmplContent, err = templates.ReadFile("templates/main.go.tmpl")
	}
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("constructor").Funcs(funcMap).Parse(string(tmplContent))

	return tmpl, err
}

func determineOutput(cmd *cli.Command) (io.Writer, error) {
	var output io.Writer
	if cmd.String("output") == "" || cmd.String("output") == "-" {
		output = os.Stdout
	} else {
		err := os.MkdirAll(filepath.Dir(cmd.String("output")), os.ModePerm)
		if err != nil {
			return nil, err
		}
		output, err = os.Create(cmd.String("output"))
		if err != nil {
			return nil, err
		}
	}

	return output, nil
}
