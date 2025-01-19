package main

import (
	"embed"
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/urfave/cli/v2"
)

//go:embed templates/*
var templates embed.FS

func main() {
	app := &cli.App{
		Name:  "constructor",
		Usage: "make something",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "name",
				Usage: "name of the thing to create a constructor for",
				Value: "Client",
			},
			&cli.StringFlag{
				Name:  "package",
				Usage: "name of the package the generated code should belong to",
				Value: "main",
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

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

type tmplInput struct {
	PackageName string
	Name        string
}

func mainAction(c *cli.Context) error {
	if c.String("name") == "" {
		return fmt.Errorf("provided \"name\" value is invalid: '%s' - Use '--name' to set it", c.String("name"))
	}
	if c.String("package") == "" {
		return fmt.Errorf("provided \"package\" value is invalid: '%s' - Use '--package' to set it", c.String("package"))
	}
	funcMap := template.FuncMap{
		"title": func(s string) string {
			return cases.Title(language.English).String(s)
		},
		"lower_first_letter": func(s string) string {
			if len(s) < 1 {
				return cases.Lower(language.English).String(s)
			}
			return cases.Lower(language.English).String(s[:1]) + s[1:]
		},
	}

	tmpl, err := determineTemplate(c, funcMap)
	if err != nil {
		return fmt.Errorf("error determining template: %w", err)
	}

	input := tmplInput{
		PackageName: c.String("package"),
		Name:        c.String("name"),
	}

	output, err := determineOutput(c)
	if err != nil {
		return fmt.Errorf("error determining output destination: %w", err)
	}

	err = tmpl.Execute(output, input)
	if err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	return nil
}

func determineTemplate(c *cli.Context, funcMap template.FuncMap) (*template.Template, error) {
	var tmplContent []byte
	var err error
	if c.Bool("test") {
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

func determineOutput(c *cli.Context) (io.Writer, error) {
	var output io.Writer
	if c.String("output") == "" || c.String("output") == "-" {
		output = os.Stdout
	} else {
		err := os.MkdirAll(filepath.Dir(c.String("output")), os.ModePerm)
		if err != nil {
			return nil, err
		}
		output, err = os.Create(c.String("output"))
		if err != nil {
			return nil, err
		}
	}

	return output, nil
}
