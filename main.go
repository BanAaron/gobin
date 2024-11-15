package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

const trashDir string = "./Trash"
const filesDir string = "./Trash/files"
const infoDir string = "./Trash/info"

type fileInfo struct {
	UUID        string
	FileName    string
	FilePath    string
	DeletedDate int64
	FileSize    int
}

func (fi fileInfo) String() string {
	return fmt.Sprintf("%s %s %s %d %d", fi.UUID, fi.FileName, fi.FilePath, fi.DeletedDate, fi.FileSize)
}

// copyFileToTrash will create a copy of a file in the Trash directory and then
// remove the original file
func copyFileToTrash(fileName string) (err error) {
	// create uuid
	uid, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("copyFileToTrash: Could not generate UUID. %s", err)
	}
	// get unix time
	unixTime := time.Now().Unix()
	// get file path
	filePath, err := filepath.Abs(fileName)
	if err != nil {
		return fmt.Errorf("copyFileToTrash: Could not get absolute path. %s", err)
	}
	// create json
	fi := fileInfo{uid.String(), fileName, filePath, unixTime, 0}
	fiJson, err := json.Marshal(fi)
	if err != nil {
		panic(err)
	}
	fmt.Println(fiJson)

	return nil
}

// init sets up the test files for local development. This should be removed
func init() {
	_, err := os.Create("test.txt")
	if err != nil {
		panic(err)
	}
	testFileContents := []byte("Hello, World!\n")
	err = os.WriteFile("test.txt", testFileContents, 0644)
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println(trashDir, filesDir, infoDir)

	args := os.Args
	if len(args) == 1 {
		fmt.Println("Please provide a file.")
		return
	}
	fileName := args[1]

	err := copyFileToTrash(fileName)
	if err != nil {
		panic(err)
	}
}
