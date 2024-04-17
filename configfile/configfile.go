package configfile

import (
	"BroGains/logfile"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Global configuration variable
var Configuration configuration

// Configuration struct
type configuration struct {
	Webserver struct {
		Port             int    `json:"port"`
		SecretSessionKey string `json:"secretsessionkey"`
		SessionName      string `json:"sessionname"`
	} `json:"webserver"`
	Database struct {
		Host         string `json:"host"`
		Port         int    `json:"port"`
		Username     string `json:"username"`
		Password     string `json:"Password"`
		DatabaseName string `json:"databasename"`
	} `json:"database"`
}

func InitConfigfile(pathConfigfile string) error {

	if !fileExists(pathConfigfile) {
		// Create a configfile, if it does not exist in the directory
		configfile, err := logfile.CreateFile(pathConfigfile)
		if err != nil {
			return err
		}

		// Convert the configuration struct into JSON format
		encodedJsonStruct, err := json.MarshalIndent(Configuration, "", "\t")
		if err != nil {
			return err
		}

		// Write the JSON data to the created configfile
		_, err = configfile.Write(encodedJsonStruct)
		if err != nil {
			return err
		}

		return fmt.Errorf("configure the settings in the configfile")
	}

	// Open the configfile
	configfile, err := openFile(pathConfigfile)
	if err != nil {
		return err
	}
	defer configfile.Close()

	// Define the configfile and assign to global variable
	err = json.NewDecoder(configfile).Decode(&Configuration)
	if err != nil {
		return err
	}

	return nil
}

// Check if the file in the given path exists
func fileExists(path string) bool {

	_, err := os.Stat(path)
	if err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		log.Fatalf("error checking file existence: %s", err)
	}

	return false
}

// Open the file within the given path
func openFile(path string) (*os.File, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %s", err)
	}

	return file, nil
}
