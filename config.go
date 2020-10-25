package main

import (
	"encoding/json"
	"io/ioutil"
)

type configInformation struct {
	RootCodeDirectory string   `json:"RootCodeDirectory"`
	PostfixFile       string   `json:"PostfixFile"`
	PreloadDirectory  string   `json:"PreloadDirectory"`
	EmFlags           string   `json:"EmFlags"`
	IncludeDirectory  string   `json:"IncludeDirectory"`
	WebBuildFlags     []string `json:"WebBuildFlags"`
}

func readConfigInformation() (configInformation, error) {
	return readConfigByFile("config.json")
}

func readConfigByFile(filePath string) (configInformation, error) {
	info := configInformation{}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return info, err
	}

	if err = json.Unmarshal(data, &info); err != nil {
		return info, err
	}

	return info, nil
}
