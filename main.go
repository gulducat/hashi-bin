package main

import (
	"log"
	"os"

	"github.com/gulducat/hashi-releases/types"
	"github.com/gulducat/hashi-releases/vars"
	"github.com/mitchellh/cli"
)

var (
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

func main() {
	log.SetFlags(0) // remove timestamp from log messages
	c := cli.NewCLI("hashi-bin", "0.1.0")
	// c.HelpFunc = HelpyHelp(c.Name)
	c.Args = os.Args[1:]
	// c.GlobalFlags = ......
	index := types.NewIndex(vars.ReleasesURL + "/index.json")
	c.Commands = GetCommands(c, &index)
	// note: i'm not using this quite right... TopLevelHelp is actually handling this.
	c.HiddenCommands = GetHiddenCommands(c)
	// c.Commands["version"] = versionFactory
	exitStatus, err := c.Run()
	if err != nil {
		panic(err)
	}
	os.Exit(exitStatus)
}

// func Route() *chi.Mux {
// 	r := chi.NewRouter()
// 	r.Use(middleware.RealIP)
// 	r.Use(middleware.Logger)
// 	r.Use(middleware.SetHeader("Content-Type", "application/json"))
// 	r.Get("/", handleRoot)
// 	r.Get("/latest/{product}", handleProductLatest)
// 	r.Get("/versions/{product}", handleListVersions)
// 	r.Get("/list", handleListProducts)
// 	return r
// }
