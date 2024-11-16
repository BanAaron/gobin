package FileManager

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

const filesDir string = "./Trash/files"
const infoDir string = "./Trash/info"

// FileInfo is the data stored in Trash/info to manage deleting and restoring
// files.
type FileInfo struct {
	UUID               uuid.UUID
	FileName           string
	FilePath           string
	BackupFileLocation string
	DeletedDate        int64
}

// String returns a string representing FileInfo. This is for debugging purposes
// only.
func (fi FileInfo) String() string {
	return fmt.Sprintf("%s %s %s %s %d", fi.UUID, fi.FileName, fi.FilePath, fi.BackupFileLocation, fi.DeletedDate)
}

// ToJson converts FileInfo into json format.
func (fi FileInfo) ToJson() (jsonData []byte, err error) {
	jsonData, err = json.Marshal(fi)
	if err != nil {
		return []byte{}, fmt.Errorf("FileInfo.ToJson(): Could not marshal JSON. %s", err)
	}
	return jsonData, nil
}

// NewFileInfo creates an instance of FileInfo
func NewFileInfo(absPath string) (FileInfo, error) {
	if !filepath.IsAbs(absPath) {
		return FileInfo{}, fmt.Errorf("provided filepath is not absolute: %s", absPath)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return FileInfo{}, fmt.Errorf("failed to generate UUID: %w", err)
	}

	backupFileLocation, err := filepath.Abs(filesDir)
	if err != nil {
		return FileInfo{}, fmt.Errorf("failed to determine backup file location: %w", err)
	}

	return FileInfo{
		UUID:               id,
		FileName:           filepath.Base(absPath),
		FilePath:           absPath,
		BackupFileLocation: filepath.Join(backupFileLocation, id.String()),
		DeletedDate:        time.Now().Unix(),
	}, nil
}

// MoveFile moves the file from its original location to ./Trash/files
//
// The file will be renamed to an uuid matching to a file in ./Trash/info
func MoveFile(fi FileInfo) (err error) {
	reader, err := os.Open(fi.FilePath)
	defer func() {
		err := reader.Close()
		if err != nil {
			panic(err)
		}
	}()
	writer, err := os.Create(path.Join(filesDir, fi.UUID.String()))
	if err != nil {
		panic(err)
	}
	defer func() {
		err := writer.Close()
		if err != nil {
			panic(err)
		}
	}()
	_, err = writer.ReadFrom(reader)
	if err != nil {
		return fmt.Errorf("moveFile: Could not move file. %s", err)
	} else {
		err := os.Remove(fi.FilePath)
		if err != nil {
			return fmt.Errorf("moveFile: Could not remove file. %s", err)
		}
	}
	return nil
}
