package main

import (
	"errors"
	"fmt"
	"regexp"
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
	// exclude some versions based on options
	opts := GetOptions()
	versions := make(map[string]*Version)
	reBeta := regexp.MustCompile(`-(beta|rc)`)
	reEnt := regexp.MustCompile(`\+ent`)
	for s, v := range p.Versions {
		// hide -beta* and -rc* if not -with-beta
		if !opts.beta && reBeta.FindStringIndex(s) != nil {
			continue
		}
		// hide +ent if not -only-enterprise
		if !opts.ent && reEnt.FindStringIndex(s) != nil {
			continue
		}
		// show only +ent if -only-enterprise
		if opts.ent && reEnt.FindStringIndex(s) == nil {
			continue
		}
		versions[s] = v
	}

	// do sorting
	collection := make(version.Collection, len(versions))
	var idx int
	for k, _ := range versions {
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
