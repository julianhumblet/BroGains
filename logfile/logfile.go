package logfile

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Create or open the logfile and set it as default output for logs
func InitLogfile(pathLogfile string) (*os.File, error) {

	// Get the directory of the given logfile path
	dirLogfile := filepath.Dir(pathLogfile)

	// Check if the directory exists
	if !dirExists(dirLogfile) {
		return nil, fmt.Errorf("logfile error: directory %s does not exist", dirLogfile)
	}

	// Create or open the logfile within the given path
	logfile, err := CreateFile(pathLogfile)
	if err != nil {
		return nil, err
	}

	// Define the created logfile as standard output for logs
	log.SetOutput(logfile)

	return logfile, nil
}

// Check if the given directory exists
func dirExists(dir string) bool {

	_, err := os.Stat(dir)
	if err != nil {
		fmt.Println(
			fmt.Errorf("error checking directory existence: %s", err),
		)
	}

	return !os.IsNotExist(err)
}

// Create a new file
func CreateFile(path string) (*os.File, error) {

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0640)
	if err != nil {
		return nil, fmt.Errorf("error creating file: %s", err)
	}

	return file, nil
}
