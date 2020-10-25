package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	configFile := flag.String("configFile", "config.json", "File to read settings for the build")
	useLinuxSlash := flag.Bool("useLinuxSlash", true, "True means use linux slash, flash means use windows slash")
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
		byteBuffer.WriteString(" ")
		return nil
	})
	if err != nil {
		fmt.Println("Error walking the filesystem path:", configInfo.RootCodeDirectory, "Error:", err)
		return
	}

	for _, buildFlag := range configInfo.WebBuildFlags {
		byteBuffer.WriteString("-D ")
		byteBuffer.WriteString(buildFlag)
		byteBuffer.WriteString(" ")
	}

	byteBuffer.WriteString(configInfo.EmFlags)
	byteBuffer.WriteString(" ")

	byteBuffer.WriteString("-I")
	byteBuffer.WriteString(configInfo.IncludeDirectory)
	byteBuffer.WriteString(" ")

	err = filepath.Walk(configInfo.PreloadDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Given an error when walking the PreloadPath:", configInfo.PreloadDirectory, "Error: ", err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		byteBuffer.WriteString("--preload-file ")
		byteBuffer.WriteString(path)
		byteBuffer.WriteString(" ")
		return nil
	})
	if err != nil {
		fmt.Println("Error walking the filesystem preload path:", configInfo.PreloadDirectory, "Error:", err)
		return
	}

	byteBuffer.WriteString("--post-js ")
	byteBuffer.WriteString(configInfo.PostfixFile)

	commandToRun := byteBuffer.String()
	if *useLinuxSlash {
		commandToRun = strings.Replace(commandToRun, "\\", "/", -1)
	}
	fmt.Println("Running command:", commandToRun, "\n")
	out, err := exec.Command(commandToRun).Output()
	if err != nil {
		fmt.Println("Error running the above command!", err, "out:", string(out))
		return
	}

	fmt.Println("Success!", string(out))
}
