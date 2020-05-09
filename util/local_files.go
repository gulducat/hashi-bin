package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

func ListInstalled(product string) ([]string, error) {
	// get current symlink target if present
	var current string
	link := LinkPath(product)
	target, err := os.Readlink(link)
	if err == nil {
		log.Printf("%s -> %s\n", link, target)
		_, current = path.Split(target)
	}

	// ls hashi-bin/{product}/ to discover installed versions
	binDir, err := BinDir(product)
	if err != nil {
		log.Println(err)
		return []string{}, err
	}
	fileInfo, err := ioutil.ReadDir(binDir)
	if err != nil {
		log.Println(err)
		return []string{}, err
	}

	// build list of versions with "*" indicating currently-in-use
	installed := []string{}
	for _, file := range fileInfo {
		name := file.Name()
		if name == current {
			installed = append(installed, fmt.Sprintf("%s *", name))
		} else {
			installed = append(installed, name)
		}
	}
	return installed, nil
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
