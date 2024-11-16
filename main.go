package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/banaaron/gobin/FileManager"
)

var Red = "\033[31m"
var Green = "\033[32m"

// init sets up the test files for local development. This should be removed
func init() {
	testFileContents := []byte("Hello, World!\n")
	err := os.WriteFile("testetxt", testFileContents, 0644)
	if err != nil {
		panic(err)
	}
}

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Println("Please provide a filePath.")
		return
	}
	if len(args) >= 3 {
		fmt.Println("Only one filePath at a time for now.")
	}

	filePath := args[1]
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		fmt.Println(Red+"Unable to determine absolute path for", filePath, "Error:", err)
		os.Exit(1)
	}

	fileInfo, err := FileManager.NewFileInfo(absPath)
	if err != nil {
		fmt.Println(Red+"Unable to create FileInfo for", filePath, "Error:", err)
		os.Exit(1)
	}

	jsonData, err := fileInfo.ToJson()
	if err != nil {
		fmt.Println(Red+"Unable to convert FileInfo to Json for", filePath, "Error:", err)
		os.Exit(1)
	}

	path := filepath.Join("Trash/info/", fileInfo.UUID.String()) + ".json"
	err = os.WriteFile(path, jsonData, 0644)
	if err != nil {
		fmt.Println(Red+"Unable to write Json for", filePath, "Error:", err)
		os.Exit(1)
	}
	err = FileManager.MoveFile(fileInfo)
	if err != nil {
		fmt.Println(Red+"Unable to move file", filePath, "Error:", err)
		os.Exit(1)
	}
	os.Exit(0)
}
