package types

import (
	"regexp"

	"github.com/gulducat/hashi-releases/util"
	"github.com/gulducat/hashi-releases/vars"
)

type Version struct {
	product    *Product // parent
	Product    string   `json:"name"`
	Version    string   `json:"version"`
	SHASums    string   `json:"shasums"`
	SHASumsSig string   `json:"shasums_signature"`
	Builds     []*Build `json:"builds"`
}

func (v *Version) String() string {
	return v.Version
}

func (v *Version) GetBuild(os string, arch string) *Build {
	// TODO: feels bad, darwin arches for vagrant .dmg downloads...
	arches := []string{arch}
	if vars.LocalOS == "darwin" {
		arches = []string{"amd64", "x86_64"}
	}
	for _, b := range v.Builds {
		for _, a := range arches {
			if b.OS == os && b.Arch == a {
				return b
			}
		}
	}
	return nil
}

func (v *Version) GetBuildForLocal() *Build {
	return v.GetBuild(vars.LocalOS, vars.LocalArch)
}

func (v *Version) IsActive() bool {
	current, _, _ := util.CurrentActive(v.Product)
	return v.Version == current
}

func (v *Version) IsBeta() bool {
	re := regexp.MustCompile(`-(beta|rc)`)
	return re.FindStringIndex(v.Version) != nil
}

func (v *Version) IsEnterprise() bool {
	re := regexp.MustCompile(`\+ent`)
	return re.FindStringIndex(v.Version) != nil
}
