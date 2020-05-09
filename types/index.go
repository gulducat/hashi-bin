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
	tmpDir := path.Join(os.TempDir(), "hashi-releases") // # TODO: rename
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

	// this intermediary `products` var is so the Index only gets core Products unless -all
	var products map[string]*Product
	if err = json.Unmarshal(b, &products); err != nil {
		return index, err
	}

	index = Index{
		Products: make(map[string]*Product),
	}
	for n, p := range products {
		index.Products[n] = p
	}

	for _, v := range index.Products {
		if err = v.sortVersions(); err != nil {
			return index, err
		}
	}
	return index, nil
}

func (i *Index) GetProductVersion(name string, version string) (*Product, *Version, error) {
	p, err := i.GetProduct(name)
	if err != nil {
		return nil, nil, err
	}
	v, err := p.GetVersion(version)
	if err != nil {
		return p, nil, err
	}
	return p, v, nil
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
