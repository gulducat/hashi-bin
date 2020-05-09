package main

import (
	"log"
	"os"

	"github.com/gulducat/hashi-releases/types"
	"github.com/gulducat/hashi-releases/vars"
	"github.com/mitchellh/cli"
)

const version = "0.1.0"

// TODO: help suffix?

func main() {
	log.SetFlags(0) // remove timestamp from log messages
	c := cli.NewCLI("armon", version)
	// c.HelpFunc = HelpyHelp(c.Name)
	c.Args = os.Args[1:]
	// c.GlobalFlags = ......
	index, err := types.NewIndex(vars.ReleasesURL + "/index.json")
	if err != nil {
		log.Fatal(err)
	}
	c.Commands = GetCommands(c, &index)
	// note: i'm not using this quite right... TopLevelHelp is actually handling this.
	c.HiddenCommands = GetHiddenCommands(c)
	// c.Commands["version"] = versionFactory
	exitStatus, err := c.Run()
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(exitStatus)
}
