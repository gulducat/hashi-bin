package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

type Build struct {
	Product  string `json:"name"`
	Version  string `json:"version"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Filename string `json:"filename"`
	URL      string `json:"url"`
}

func (b *Build) Download() ([]byte, error) {
	// TODO: set a User Agent https://stackoverflow.com/questions/13263492/set-useragent-in-http-request
	resp, err := http.Get(b.URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bts, nil
}

func (b *Build) DownloadAndCheck() ([]byte, error) {
	bts, err := b.Download()
	if err != nil {
		return nil, err
	}
	err = CheckBytes(b.Filename, bts)
	if err != nil {
		return nil, err
	}
	return bts, nil
}

func (b *Build) DownloadAndSave(filePath string) (string, error) {
	bts, err := b.DownloadAndCheck()
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(filePath, bts, 0644)
	if err != nil {
		return "", err
	}
	log.Printf("downloaded: %s", filePath)
	return filePath, err
}

func (b *Build) DownloadAndExtract(dir string, product string) (string, error) {
	bts, err := b.DownloadAndCheck()
	filePath, err := ExtractZip(product, dir, bts)
	if err != nil {
		return "", err
	}
	log.Printf("extracted: %s", filePath)
	return filePath, err
}

func (b *Build) Install() error {
	binDir, err := BinDir(b.Product)
	if err != nil {
		return err
	}
	// TODO: check if already installed before downloading?
	_, err = b.DownloadAndExtract(binDir, b.Product)
	if err != nil {
		return err
	}
	// TODO: ExtractZip should put the file directly where we want it?
	filePath := path.Join(binDir, b.Product)
	newFilePath := path.Join(binDir, b.Version)
	err = os.Rename(filePath, newFilePath)
	if err != nil {
		return err
	}
	log.Printf("installed to: %s\n", newFilePath)
	log.Printf("to use: `hashi-bin use %s %s`\n", b.Product, b.Version)
	return nil
}

func (b *Build) Uninstall() error {
	binDir, err := BinDir(b.Product)
	if err != nil {
		return err
	}
	filePath := path.Join(binDir, b.Version)
	err = os.Remove(filePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	// TODO: only remove symlink if uninstalling current active version - another reason to make a separate link-reading function.
	return RemoveLink(b.Product)
}

func (b *Build) Link() error {
	binDir, err := BinDir(b.Product)
	if err != nil {
		return err
	}
	filePath := path.Join(binDir, b.Version)
	link := LinkPath(b.Product)
	// TODO: check if filePath exists, if not, suggest `install` ?
	RemoveLink(b.Product)
	log.Printf("Creating symlink %s -> %s\n", link, filePath)
	return os.Symlink(filePath, link)
}
