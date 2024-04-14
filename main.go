package main

import (
	"BroGains/logfile"
	"fmt"
	"os"
)

func init() {

	// Set the desired configuration settings here
	pathLogfile := "./logfile.log"

	// Initialize logfile
	_, err := logfile.InitLogfile(pathLogfile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {

}
