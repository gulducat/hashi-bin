package main

import (
	"io/ioutil"
	"net/http"
)

const UserAgent = "HashiCorp hashi-bin CLI utility"

func HTTPGet(url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", UserAgent)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func HTTPGetBody(url string) ([]byte, error) {
	resp, err := HTTPGet(url)
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
