package types

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gulducat/hashi-releases/util"
)

var (
	CacheFilePath = path.Join(os.TempDir(), "hashicorp.releases.json")
	CacheMaxAge   = 60 // minutes
)

type Index struct {
	Products map[string]*Product
}

func NewIndex(indexURL string) (Index, error) {
	var index Index

	b, err := GetIndexBody(indexURL)
	if err != nil {
		return index, err
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

func GetIndexBody(indexURL string) ([]byte, error) {
	if err := ExpireCache(); err != nil {
		return nil, err
	}

	b, err := ioutil.ReadFile(CacheFilePath)
	if err != nil {
		resp, err := util.HTTPGet(indexURL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		etag := resp.Header.Get("Etag")
		etag = strings.Trim(etag, "\"")
		if etag == "" {
			return nil, errors.New("no etag found in http response headers")
		}

		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		if err = os.MkdirAll(path.Dir(CacheFilePath), 0700); err != nil {
			return nil, err
		}
		if err = ioutil.WriteFile(CacheFilePath, b, 0600); err != nil {
			return nil, err
		}
	}

	return b, nil
}

func ExpireCache() error {
	stat, err := os.Stat(CacheFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	age := time.Since(stat.ModTime()).Minutes()
	if int(age) >= CacheMaxAge {
		err = os.Remove(CacheFilePath)
		if err != nil {
			return err
		}
	}
	return nil
}
