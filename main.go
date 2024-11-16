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
	UUID        uuid.UUID
	FileName    string
	FilePath    string
	DeletedDate int64
	FileSize    int
}

func (fi fileInfo) String() string {
	return fmt.Sprintf("%s %s %s %d %d", fi.UUID, fi.FileName, fi.FilePath, fi.DeletedDate, fi.FileSize)
}

func newFileInfo(fileName string) (fi fileInfo, err error) {
	uid, err := uuid.NewUUID()
	if err != nil {
		return fi, fmt.Errorf("copyFileToTrash: Could not generate UUID. %s", err)
	}
	unixTime := time.Now().Unix()
	filePath, err := filepath.Abs(fileName)
	if err != nil {
		return fi, fmt.Errorf("copyFileToTrash: Could not get absolute path. %s", err)
	}
	fi = fileInfo{uid, fileName, filePath, unixTime, 0}
	return fi, nil
}

func (fi fileInfo) Marshall() (FileInfoJson []byte, err error) {
	FileInfoJson, err = json.Marshal(fi)
	if err != nil {
		return []byte{}, fmt.Errorf("fileInfo.Marshall(): Could not marshal JSON. %s", err)
	}
	return FileInfoJson, nil
}

// buildInfoPath will create the file path for fileInfo objects
//
// for example: ./Trash/info/b5b8cc50-a3ad-11ef-b204-0a01357c9fb2.json
func buildInfoPath(fi fileInfo) string {
	return infoDir + "/" + fi.UUID.String() + ".json"
}

// copyFileToTrash will create a copy of a file in the Trash directory and then
// remove the original file
func copyFileToTrash(fileName string) (err error) {
	// create uuid
	fi, err := newFileInfo(fileName)
	if err != nil {
		return fmt.Errorf("copyFileToTrash: Could not create file info. %s", err)
	}

	fileInfoJson, err := fi.Marshall()
	if err != nil {
		return fmt.Errorf("copyFileToTrash: Could not marshal file info. %s", err)
	}

	infoPath := buildInfoPath(fi)
	err = os.WriteFile(infoPath, fileInfoJson, 0644)
	if err != nil {
		return fmt.Errorf("copyFileToTrash: Could not write file info. %s", err)
	}

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
	if len(args) != 3 {
		fmt.Println("Only one file at a time for now.")
	}

	fileName := args[1]
	err := copyFileToTrash(fileName)
	if err != nil {
		panic(err)
	}
}
