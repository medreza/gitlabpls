package app

import (
	"fmt"

	"github.com/medreza/gitlabpls/pkg/generator"
	"github.com/pelletier/go-toml"
	"github.com/pkg/browser"
	"github.com/urfave/cli/v2"
)

func New(varsCfg *toml.Tree, urlGenerator *generator.Generator, availableCmds []string) *cli.App {
	return &cli.App{
		Name:                 "gitlabpls",
		Usage:                "GitLab, please! A CLI tool to help you run new GitLab CI/CD pipeline",
		Version:              "1.0.0",
		HideHelpCommand:      true,
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:         "url",
				Aliases:      []string{"u"},
				Usage:        "generate URL",
				Subcommands:  getSubcommands(varsCfg, urlGenerator, false),
				BashComplete: getBashCompletionFunc(availableCmds),
			},
			{
				Name:         "browser",
				Aliases:      []string{"b"},
				Usage:        "generate URL and open it in your browser",
				Subcommands:  getSubcommands(varsCfg, urlGenerator, true),
				BashComplete: getBashCompletionFunc(availableCmds),
			},
		},
	}
}

func actionHandler(vars string, varsCfg *toml.Tree, urlGenerator *generator.Generator, shouldOpenBrowser bool) error {
	if !varsCfg.Has(vars) {
		fmt.Println("vars not found")
	}
	generated, err := urlGenerator.Generate(varsCfg.ToMap()[vars].(map[string]interface{}))
	if err != nil {
		return err
	}
	fmt.Printf("URL: %s\n", generated)
	if shouldOpenBrowser {
		fmt.Println("Opening browser...")
		return browser.OpenURL(generated)
	}
	return nil
}

func getSubcommands(varsCfg *toml.Tree, urlGenerator *generator.Generator, shouldOpenBrowser bool) (subCmds []*cli.Command) {
	for _, k := range varsCfg.Keys() {
		subCmd := &cli.Command{
			Name:  k,
			Usage: fmt.Sprintf("generate url with variables configured in 'vars.%s'", k),
			Action: func(context *cli.Context) error {
				return actionHandler(k, varsCfg, urlGenerator, shouldOpenBrowser)
			},
		}
		subCmds = append(subCmds, subCmd)
	}
	return subCmds
}

func getBashCompletionFunc(availableCmds []string) func(c *cli.Context) {
	return func(c *cli.Context) {
		if c.NArg() > 0 {
			return
		}
		for _, t := range availableCmds {
			fmt.Println(t)
		}
	}
}
