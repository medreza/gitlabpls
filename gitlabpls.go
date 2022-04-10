package main

import (
	"embed"
	"fmt"
	"github.com/medreza/gitlabpls/pkg/generator"
	"github.com/pelletier/go-toml"
	"github.com/pkg/browser"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

//go:embed gitlabpls.config.toml
var cfgRes embed.FS

func main() {
	var vars string
	availableCmds := []string{"url", "browser"}
	cfgBytes, _ := cfgRes.ReadFile("gitlabpls.config.toml")
	cfg, _ := toml.LoadBytes(cfgBytes)
	varsConfig := cfg.Get("vars").(*toml.Tree)

	uri := generator.New(
		cfg.Get("app.gitlab.project.base.url").(string),
		cfg.Get("app.gitlab.repo").(string),
		cfg.Get("app.gitlab.branch").(string))

	app := &cli.App{
		Usage:                "GitLab, please! A CLI app to help you start new GitLab CI/CD pipeline",
		HideHelpCommand:      true,
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:    "url",
				Aliases: []string{"u"},
				Usage:   "generate URL",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "vars",
						Required:    true,
						Usage:       "Variables group defined in the config file",
						Destination: &vars,
					},
				},
				Action: func(c *cli.Context) error {
					doProcess(vars, varsConfig, uri, false)
					return nil
				},
				BashComplete: func(c *cli.Context) {
					if c.NArg() > 0 {
						return
					}
					for _, t := range availableCmds {
						fmt.Println(t)
					}
				},
			},
			{
				Name:    "browser",
				Aliases: []string{"b"},
				Usage:   "generate URL and open it in your browser",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "vars",
						Required:    true,
						Usage:       "Variables group defined in the config file",
						Destination: &vars,
					},
				},
				Action: func(c *cli.Context) error {
					doProcess(vars, varsConfig, uri, true)
					return nil
				},
				BashComplete: func(c *cli.Context) {
					if c.NArg() > 0 {
						return
					}
					for _, t := range availableCmds {
						fmt.Println(t)
					}
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func doProcess(vars string, varsConfig *toml.Tree, uri *generator.Generator, shouldOpenBrowser bool) {
	if !varsConfig.Has(vars) {
		fmt.Println("vars not found")
	}
	generated, _ := uri.Generate(varsConfig.ToMap()[vars].(map[string]interface{}))
	fmt.Printf("URL: %s\n", generated)
	if shouldOpenBrowser {
		_ = browser.OpenURL(generated)
	}
}
