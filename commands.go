package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mitchellh/cli"
)

type ExtraOptions struct {
	all  bool
	beta bool
	ent  bool
}

func GetOptions() ExtraOptions {
	// feels a bit jank but it works..
	// TODO: surely there's a better way to parse extra flags
	// TODO: regardless, add to -h
	return ExtraOptions{
		all:  os.Getenv("HASHI_ALL") != "" || InArray(os.Args, "-all"),
		beta: os.Getenv("HASHI_BETA") != "" || InArray(os.Args, "-with-beta"),
		ent:  os.Getenv("HASHI_ENTERPRISE") != "" || InArray(os.Args, "-only-enterprise"),
	}
}

func GetCommands(c *cli.CLI, i *Index) map[string]cli.CommandFactory {
	commands := make(map[string]cli.CommandFactory)

	options := map[string]string{
		"list-available": "List available versions of a product.",
		"list":           "List installed versions of a product.",
		"download":       "Download to the current directory.",
		"install":        "Install to ~/.hashi-bin/{product}/{version} (or env $HASHI_BIN)",
		"uninstall":      "Delete ~/.hashi-bin/{product}/{version} and remove symlink.",
		"use":            "Symlink /usr/local/bin/{product} (or env $HASHI_LINKS) -> ~/.hashi-bin/{product}/{version}",
		// TODO: "clean" to remove inactive versions of all products?
	}

	// top-level help
	for option, synopsis := range options {
		option := option
		synopsis := synopsis
		commands[option] = func() (cli.Command, error) {
			return &TopLevelHelp{
				cli:      c,
				option:   option,
				index:    i,
				synopsis: synopsis,
			}, nil
		}
	}

	for _, product := range i.Products {
		product := product
		for option, _ := range options {
			option := option
			commands[option+" "+product.Name] = func() (cli.Command, error) {
				return &FancyCommand{
					index:   i,
					product: product,
					command: option,
				}, nil
			}
		}
	}

	return commands
}

func GetHiddenCommands(c *cli.CLI) []string {
	// exclude all but core products unless -all
	// note: c.Commands must already be populated
	hidden := []string{}
	opts := GetOptions()
	for cmd, _ := range c.Commands {
		parts := strings.Split(cmd, " ")
		if len(parts) != 2 {
			continue
		}
		product := parts[1]
		if !opts.all && !InArray(CoreProducts, product) {
			hidden = append(hidden, cmd)
		}
	}
	return hidden
}

// top-level command help

type TopLevelHelp struct {
	cli      *cli.CLI
	option   string
	index    *Index
	synopsis string
}

func (hc *TopLevelHelp) Synopsis() string {
	return hc.synopsis
}

func (hc *TopLevelHelp) Help() string {
	return hc.synopsis
}

func (hc *TopLevelHelp) Run(args []string) int {
	log.Println(hc.HelpTemplate())
	return 127 // 127 to match normal help, because nothing has been done..
}

// TODO: interesting, this applies only when -h is specified
func (hc *TopLevelHelp) HelpTemplate() string {
	// TODO: this help logic is a bit goofy..?
	commands := make(map[string]cli.CommandFactory)
	for cmd, cft := range hc.cli.Commands {
		if InArray(hc.cli.HiddenCommands, cmd) {
			continue
		}
		if strings.HasPrefix(cmd, hc.option+" ") {
			commands[cmd] = cft
		}
	}
	return hc.Help() + "\n\n" + hc.cli.HelpFunc(commands)
}

type FancyCommand struct {
	index   *Index
	product *Product
	command string
}

func (fc *FancyCommand) Synopsis() string {
	return "" // be vewwy vewwy quiet
}

func (fc *FancyCommand) Help() string {
	return fmt.Sprintf("provide 'latest' or a version from `hashi list-available %s` to %s",
		fc.product.Name, fc.command)
}

func (fc *FancyCommand) Run(args []string) int {
	var err error

	// These commands require no version argument
	switch fc.command {
	case "list-available":
		for _, v := range fc.product.Sorted {
			fmt.Println(v)
		}
		return 0
	case "list":
		installed, err := ListInstalled(fc.product.Name)
		if err != nil {
			log.Println(err)
			return 1
		}
		for _, result := range installed {
			result := result
			fmt.Println(result)
		}
		return 0
	}

	// all remaining commands require version
	if len(args) < 1 { // additional args will be swallowed without notice.
		log.Println(fc.Help())
		return 1
	}
	versionString := args[0]

	version, err := fc.product.GetVersion(versionString)
	if err != nil {
		log.Println(err)
		return 1
	}
	build := version.GetBuildForLocal()
	if build == nil {
		log.Println("Failed to find a build for this machine...")
		return 1
	}

	log.Printf("%s-ing %s %s\n", fc.command, fc.product.Name, version.Version)

	switch fc.command {
	case "download":
		// TODO: this feels bad, do something else to download vagrant?
		if localOS == "darwin" && InArray(DmgOnly, fc.product.Name) {
			_, err = build.DownloadAndSave(build.Filename)
		} else {
			_, err = build.DownloadAndExtract("", fc.product.Name)
		}
	case "install":
		err = build.Install()
	case "uninstall":
		err = build.Uninstall()
	case "use":
		err = build.Link()
	default:
		err = errors.New("NotImplementedError")
	}
	if err != nil {
		log.Println(err)
		return 1
	}
	return 0
}
