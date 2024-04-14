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
	logfile, err := os.OpenFile(pathLogfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)
	if err != nil {
		return nil, fmt.Errorf("logfile error: %s", err)
	}

	// Define the created logfile as standard output for logs
	log.SetOutput(logfile)

	return logfile, nil
}

// Check if the given directory exists
func dirExists(dir string) bool {

	_, err := os.Stat(dir)

	return !os.IsNotExist(err)
}
