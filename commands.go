package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/hashicorp/go-version"
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
				commands: &commands,
				option:   option,
				index:    i,
				synopsis: synopsis,
			}, nil
		}
	}

	for _, product := range i.Products {

		name := product.Name
		versions := product.Sorted

		commands["list-available "+name] = func() (cli.Command, error) {
			return &ListAvailableCommand{
				versions: versions,
			}, nil
		}
		commands["list "+name] = func() (cli.Command, error) {
			return &ListCommand{
				product: name,
			}, nil
		}
		for option, _ := range options {
			option := option
			// TODO: something other than this hard-coded list/list-available exclusion...?
			if option == "list" || option == "list-available" {
				continue
			}
			commands[option+" "+name] = func() (cli.Command, error) {
				return &FancyCommand{
					index:   i,
					product: name,
					command: option,
				}, nil
			}
		}
	}

	return commands
}

// top-level command help

type TopLevelHelp struct {
	cli      *cli.CLI
	commands *map[string]cli.CommandFactory
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
	for cmd, cft := range *hc.commands {
		if strings.HasPrefix(cmd, hc.option+" ") {
			commands[cmd] = cft
		}
	}
	return hc.Help() + "\n\n" + hc.cli.HelpFunc(commands)
}

// list-available reads from releases API

type ListAvailableCommand struct {
	versions version.Collection
}

func (lc *ListAvailableCommand) Help() string {
	return ""
}

func (lc *ListAvailableCommand) Synopsis() string {
	return ""
}

func (lc *ListAvailableCommand) Run(args []string) int {
	for _, v := range lc.versions {
		fmt.Println(v)
	}
	return 0
}

// `list` reads the local filesystem

type ListCommand struct {
	product string
}

func (ic *ListCommand) Synopsis() string {
	return ""
}

func (ic *ListCommand) Help() string {
	return fmt.Sprintf("show installed versions of %s", ic.product)
}

func (ic *ListCommand) Run(args []string) int {
	// TODO: split this stuff out, so Help() can show all products?

	// get current symlink target if present
	current := ""
	link := LinkPath(ic.product)
	target, err := os.Readlink(link)
	if err == nil {
		log.Printf("%s -> %s\n", link, target)
		_, current = path.Split(target)
	}

	// ls hashi-bin/{product}/ to discover installed versions
	binDir, err := BinDir(ic.product)
	if err != nil {
		log.Println(err)
		return 1
	}
	fileInfo, err := ioutil.ReadDir(binDir)
	if err != nil {
		log.Println(err)
		return 1
	}

	// prepend * to current active version
	for _, file := range fileInfo {
		name := file.Name()
		if name == current {
			fmt.Printf("* %s\n", name)
		} else {
			fmt.Printf("  %s\n", name)
		}
	}
	return 0
}

// all other commands are "FancyCommand"s
// download, install, use, uninstall

type FancyCommand struct {
	index   *Index
	product string
	command string
}

func (fc *FancyCommand) Synopsis() string {
	return "" // be vewwy vewwy quiet
}

func (fc *FancyCommand) Help() string {
	return fmt.Sprintf("provide 'latest' or a version from `hashi list-available %s` to %s", fc.product, fc.command)
}

func (fc *FancyCommand) Run(args []string) int {
	var err error

	if len(args) < 1 { // additional args will be swallowed without notice.
		log.Println(fc.Help())
		return 1
	}
	versionString := args[0]

	product, version, err := fc.index.GetProductVersion(fc.product, versionString)
	if err != nil {
		log.Println(err)
		return 1
	}
	build := version.GetBuildForLocal()
	if build == nil {
		log.Println("Failed to find a build for this machine...")
		return 1
	}

	log.Printf("%s-ing %s %s\n", fc.command, product.Name, version.Version)

	switch fc.command {
	case "download":
		// TODO: this feels bad, do something else to download vagrant?
		if localOS == "darwin" && InArray(DmgOnly, product.Name) {
			_, err = build.DownloadAndSave(build.Filename)
		} else {
			_, err = build.DownloadAndExtract("", product.Name)
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
