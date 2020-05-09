package types

import (
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/gulducat/hashi-releases/util"
)

type Build struct {
	Product  string `json:"name"`
	Version  string `json:"version"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Filename string `json:"filename"`
	URL      string `json:"url"`
}

func (b *Build) String() string {
	return b.Filename
}

func (b *Build) DownloadAndCheck() ([]byte, error) {
	bts, err := util.HTTPGetBody(b.URL)
	if err != nil {
		return nil, err
	}
	err = util.CheckBytes(b.Filename, bts)
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
	filePath, err := util.ExtractZip(product, dir, bts)
	if err != nil {
		return "", err
	}
	log.Printf("extracted: %s", filePath)
	return filePath, err
}

func (b *Build) Install() error {
	binDir, err := util.BinDir(b.Product)
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
	binDir, err := util.BinDir(b.Product)
	if err != nil {
		return err
	}
	filePath := path.Join(binDir, b.Version)
	err = os.Remove(filePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	// TODO: only remove symlink if uninstalling current active version - another reason to make a separate link-reading function.
	return util.RemoveLink(b.Product)
}

// TODO: copy file instead of symlink.
func (b *Build) Link() error {
	binDir, err := util.BinDir(b.Product)
	if err != nil {
		return err
	}
	filePath := path.Join(binDir, b.Version)
	link := util.LinkPath(b.Product)
	// TODO: check if filePath exists, if not, suggest `install` ?
	util.RemoveLink(b.Product)
	log.Printf("Creating symlink %s -> %s\n", link, filePath)
	return os.Symlink(filePath, link)
}
