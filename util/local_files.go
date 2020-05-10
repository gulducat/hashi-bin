package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

func CurrentActive(product string) (version string, link string, target string) {
	link = LinkPath(product)
	target, err := os.Readlink(link)
	if err != nil {
		return "", "", ""
	}
	_, version = path.Split(target)
	return version, link, target
}

func ListInstalled(product string) ([]string, error) {
	var i []string

	// ls hashi-bin/{product}/ to discover installed versions
	binDir, err := BinDir(product)
	if err != nil {
		log.Println(err)
		return i, err
	}
	fileInfo, err := ioutil.ReadDir(binDir)
	if err != nil {
		log.Println(err)
		return i, err
	}

	current, link, target := CurrentActive(product)
	if current != "" {
		log.Printf("%s -> %s\n", link, target)
	}

	for _, file := range fileInfo {
		name := file.Name()
		if name == current {
			i = append(i, fmt.Sprintf("%s (current)", name))
		} else {
			i = append(i, name)
		}
	}

	return i, nil
}

func BinDir(product string) (string, error) {
	binDir, ok := os.LookupEnv("HASHI_BIN")
	if ok {
		binDir = path.Join(binDir, product)
	} else {
		binDir = path.Join(os.Getenv("HOME"), ".hashi-bin", product)
	}
	err := os.MkdirAll(binDir, 0755)
	if err != nil && !os.IsExist(err) {
		return "", err
	}
	return binDir, nil
}

func LinkPath(product string) string {
	dir, ok := os.LookupEnv("HASHI_LINKS")
	if !ok {
		dir = "/usr/local/bin"
	}
	return path.Join(dir, product)
}

func RemoveLink(product string) error {
	err := os.Remove(LinkPath(product))
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
