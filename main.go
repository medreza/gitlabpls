package main

import (
	"embed"
	"log"
	"os"

	"github.com/medreza/gitlabpls/internal/app"
	"github.com/medreza/gitlabpls/pkg/generator"
	"github.com/pelletier/go-toml"
)

//go:embed gitlabpls.config.toml
var cfgRes embed.FS

func main() {
	availableCmds := []string{"url", "browser"}
	cfgBytes, _ := cfgRes.ReadFile("gitlabpls.config.toml")
	cfg, _ := toml.LoadBytes(cfgBytes)
	varsCfg := cfg.Get("vars").(*toml.Tree)

	urlGenerator := generator.New(
		cfg.Get("app.gitlab.project.base.url").(string),
		cfg.Get("app.gitlab.repo").(string),
		cfg.Get("app.gitlab.branch").(string))

	cliApp := app.New(varsCfg, urlGenerator, availableCmds)

	err := cliApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
