package main

import (
	"log"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/mitchellh/cli"
)

var (
	index      = NewIndex()
	ListenPort = ":8080"
	// appVersion         = "v1.0.1"
	versionCommandHelp = "output the version of this application"
)

// type versionCommand struct{}

// func (vc versionCommand) Help() string {
// 	return versionCommandHelp
// }

// func (vc versionCommand) Run(args []string) int {
// 	fmt.Println(appVersion)
// 	return 0
// }

// func (vc versionCommand) Synopsis() string {
// 	return versionCommandHelp
// }

// func versionFactory() (cli.Command, error) {
// 	return versionCommand{}, nil
// }

// TODO: help suffix?
func HelpyHelp(app string) cli.HelpFunc {
	// basic := cli.BasicHelpFunc(app)
	return func(commands map[string]cli.CommandFactory) string {
		return strings.TrimSpace(`
		Usage: hashi-bin [--version] [--help] <product> [version] [-all] [-with-beta] [-only-enterprise]

		`)
	}
}

func main() {
	log.SetFlags(0) // remove timestamp from log messages
	c := cli.NewCLI("hashi-bin", "0.1.0")
	// c.HelpFunc = HelpyHelp(c.Name)
	c.Args = os.Args[1:]
	// c.GlobalFlags = ......
	c.Commands = GetCommands(c, &index)
	c.HiddenCommands = GetHiddenCommands(c)
	// c.Commands["version"] = versionFactory
	exitStatus, err := c.Run()
	if err != nil {
		panic(err)
	}
	os.Exit(exitStatus)
}

func Route() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))
	r.Get("/", handleRoot)
	r.Get("/latest/{product}", handleProductLatest)
	r.Get("/versions/{product}", handleListVersions)
	r.Get("/list", handleListProducts)
	return r
}
