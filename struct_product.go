package main

import (
	"errors"
	"fmt"
	"sort"

	"github.com/hashicorp/go-version"
)

type Product struct {
	Name     string              `json:"name"`
	Versions map[string]*Version `json:"versions"`
	Sorted   version.Collection
	isSorted bool
}

func (p *Product) GetVersion(version string) (*Version, error) {
	if version == "latest" {
		return p.LatestVersion(), nil
	}
	v, ok := p.Versions[version]
	if !ok {
		return nil, errors.New(fmt.Sprintf("invalid version for %s", p.Name))
	}
	return v, nil
}

func (p *Product) LatestVersion() *Version {
	versionString := p.Sorted[len(p.Sorted)-1].Original()
	return p.Versions[versionString]
}

func (p *Product) sortVersions() error {
	collection := make(version.Collection, len(p.Versions))
	var idx int
	for k, _ := range p.Versions {
		v, err := version.NewVersion(k)
		if err != nil {
			return err
		}
		collection[idx] = v
		idx++
	}
	p.Sorted = collection
	sort.Sort(p.Sorted)
	p.isSorted = true
	return nil
}
