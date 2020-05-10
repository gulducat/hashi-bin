package types

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/gulducat/hashi-releases/util"
)

type Index struct {
	Products map[string]*Product
}

func NewIndex(IndexURL string) (Index, error) {
	var index Index
	resp, err := util.HTTPGet(IndexURL)
	if err != nil {
		return index, err
	}
	defer resp.Body.Close()
	etag := resp.Header.Get("Etag")
	etag = strings.Trim(etag, "\"")
	if etag == "" {
		return index, errors.New("no etag found")
	}

	// TODO: cache expiration or purge
	tmpDir := path.Join(os.TempDir(), "hashi-bin")
	cacheFilePath := path.Join(tmpDir, etag, etag+".index")

	b, err := ioutil.ReadFile(cacheFilePath)
	if err != nil {
		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return index, err
		}
		if err = os.MkdirAll(path.Dir(cacheFilePath), 0700); err != nil {
			return index, err
		}
		if err = ioutil.WriteFile(cacheFilePath, b, 0600); err != nil {
			return index, err
		}
	}

	if err = json.Unmarshal(b, &index.Products); err != nil {
		return index, err
	}

	// massage the datas
	for _, p := range index.Products {
		if err = p.sortVersions(); err != nil {
			return index, err
		}
		// populate children's parent fields for convenience
		// surely there is a better way?
		p.index = &index
		for _, v := range p.Versions {
			v.product = p
			for _, b := range v.Builds {
				b.version = v
			}
		}
	}

	return index, nil
}

func (i *Index) GetVersion(product string, version string) (*Version, error) {
	p, err := i.GetProduct(product)
	if err != nil {
		return nil, err
	}
	v, err := p.GetVersion(version)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (i *Index) GetProduct(name string) (*Product, error) {
	product, ok := i.Products[name]
	if !ok {
		return nil, errors.New("invalid product name")
	}
	return product, nil
}

func (i *Index) ListProducts() []string {
	products := make([]string, len(i.Products))
	var idx int
	for k, _ := range i.Products {
		products[idx] = k
		idx++
	}
	return products
}
