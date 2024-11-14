package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// panicOnError will panic off error is not nil
func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

// copyFileToTrash will create a copy of a file in the Trash directory and the remove
// the original file
func copyFileToTrash(sourceFile string) (err error) {
	reader, err := os.Open(sourceFile)
	panicOnError(err)
	defer func() {
		err := reader.Close()
		panicOnError(err)
	}()

	outputFileName := appendUnixTime(sourceFile)
	writer, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := writer.Close()
		panicOnError(err)
	}()
	_, err = writer.ReadFrom(reader)
	panicOnError(err)
	if err == nil {
		err := os.Remove(sourceFile)
		panicOnError(err)
	}
	return nil
}

// appendUnixTime will append a file name with double underscore _^_ and then
// the current unix time.
//
// For example `foo.txt` would become `foo.text__1731550759`
func appendUnixTime(sourceFile string) string {
	unixTime := "__" + strconv.FormatInt(time.Now().Unix(), 10)
	return "./Trash/" + sourceFile + unixTime
}

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Println("Please provide a file.")
		return
	}
	fileName := args[1]

	err := copyFileToTrash(fileName)
	panicOnError(err)
}
