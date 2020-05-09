package main

import (
	"os"
	"path"
	"runtime"

	"github.com/hashicorp/go-hclog"
)

const (
	ReleasesURL    = "https://releases.hashicorp.com"
	ReleasesDomain = "releases.hashicorp.com"

	indexSuffix = ".index"

	localOS   = runtime.GOOS
	localArch = runtime.GOARCH
)

var (
	tmpDir = path.Join(os.TempDir(), "hashi-releases")

	logger hclog.Logger // TODO: logger...?

	// for the general public, only show these.
	// --all or HASHI_ALL_PRODUCTS env var will show _all_ products
	CoreProducts = []string{
		"consul",
		"nomad",
		"packer",
		"sentinel",
		// "serf", // ?
		"terraform",
		"vagrant",
		"vault",
	}
	// not available as .zip on OSX, only .dmg
	DmgOnly = []string{
		"vagrant",
	}
)
