package main

type Version struct {
	Product    string   `json:"name"`
	Version    string   `json:"version"`
	SHASums    string   `json:"shasums"`
	SHASumsSig string   `json:"shasums_signature"`
	Builds     []*Build `json:"builds"`
}

func (v *Version) GetBuild(os string, arch string) *Build {
	// TODO: feels bad, dumb darwin arches for vagrant .dmg downloads...
	arches := []string{arch}
	if localOS == "darwin" {
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
	return v.GetBuild(localOS, localArch)
}
