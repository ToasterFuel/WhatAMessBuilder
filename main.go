package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	configFile := flag.String("configFile", "config.json", "File to read settings for the build")
	flag.Parse()
	configInfo, err := readConfigByFile(*configFile)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Config read:", configInfo, "\n")
	var byteBuffer bytes.Buffer
	byteBuffer.WriteString("emcc ")

	err = filepath.Walk(configInfo.RootCodeDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Given an error when walking the RootCodeDirectory:", configInfo.RootCodeDirectory, "Error: ", err)
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".cpp") {
			return nil
		}
        byteBuffer.WriteString(path)
		return nil
	})
	if err != nil {
		fmt.Println("Error walking the filesystem:", err)
		return
	}
	fmt.Println("Waa:", cppFiles)
}
