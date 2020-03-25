package main

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path"
)

func InArray(arr []string, str string) bool {
	for _, x := range arr {
		if x == str {
			return true
		}
	}
	return false
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

func ExtractZip(product, parentDir string, bts []byte) (string, error) {
	zipReader, err := zip.NewReader(bytes.NewReader(bts), int64(len(bts)))
	if err != nil {
		return "", err
	}
	finalPath := path.Join(parentDir, product)
	outFile, err := os.OpenFile(finalPath, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return "", err
	}
	var content io.ReadCloser
	for _, f := range zipReader.File {
		if f.Name == product {
			zipFile, err := f.Open()
			if err != nil {
				return "", err
			}
			content = zipFile
		}
	}
	_, err = io.Copy(outFile, content)
	if err != nil {
		return "", nil
	}
	return finalPath, nil
}
