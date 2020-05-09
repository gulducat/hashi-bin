package vars

import (
	"runtime"

	"github.com/hashicorp/go-hclog"
)

const (
	ReleasesURL = "https://releases.hashicorp.com"
	LocalOS     = runtime.GOOS
	LocalArch   = runtime.GOARCH
)

var (
	logger hclog.Logger // TODO: logger...?

	// for the general public, only show these.
	// -all or HASHI_ALL env var will show _all_ products
	CoreProducts = []string{
		"consul",
		"nomad",
		"packer",
		"sentinel",
		"terraform",
		"vagrant",
		"vault",
	}
	// not available as .zip on OSX, only .dmg
	DmgOnly = []string{
		"vagrant",
	}
)
