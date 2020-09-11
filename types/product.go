package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/gulducat/hashi-bin/util"
	"github.com/hashicorp/go-version"
)

type Product struct {
	index    *Index              // parent
	Name     string              `json:"name"`
	Versions map[string]*Version `json:"versions"`
	Sorted   version.Collection
	isSorted bool
}

func NewProduct(indexURL string) (*Product, error) {
	var product Product

	b, err := GetIndexBody(indexURL, false)
	if err != nil {
		return &product, err
	}
	if err = json.Unmarshal(b, &product); err != nil {
		return &product, err
	}

	if err = product.sortVersions(); err != nil {
		return &product, err
	}

	for _, v := range product.Versions {
		v.product = &product
		for _, b := range v.Builds {
			b.version = v
		}
	}

	return &product, nil
}

func (p *Product) String() string {
	return p.Name
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
	versions := p.ListVersions()
	version := versions[len(versions)-1]
	return p.Versions[version]
}

func (p *Product) ListVersions() []string {
	var versions []string

	// exclude some versions based on options
	opts := util.GetOptions()

	for _, v := range p.Sorted {
		vString := v.Original()

		// hide vault ".hsm" files - they're not useful for us here.
		if strings.HasSuffix(vString, ".hsm") {
			continue
		}

		if opts.All { // show all the things
			versions = append(versions, vString)
			continue
		}

		v, _ := p.GetVersion(vString)
		if !opts.Beta && v.IsBeta() { // hide -beta* and -rc* if not -with-beta
			continue
		}
		if !opts.Ent && v.IsEnterprise() { // hide +ent if not -enterprise
			continue
		}
		if opts.Ent && !v.IsEnterprise() { // show only +ent if -enterprise
			continue
		}

		versions = append(versions, vString)
	}

	return versions
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
