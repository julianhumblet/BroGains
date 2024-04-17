package main

import (
	"BroGains/configfile"
	"BroGains/logfile"
	"BroGains/webserver"
	"fmt"
	"log"
	"os"
)

func init() {

	// Set the desired configuration settings here
	pathLogfile := "./logfile.log"
	pathConfigfile := "./config.json"

	// Initialize logfile
	_, err := logfile.InitLogfile(pathLogfile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Initialize configfile
	err = configfile.InitConfigfile(pathConfigfile)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	webserver.SetupSessions()

	webserverPort := configfile.Configuration.Webserver.Port

	// Start the webserver
	webserver.StartWebserver(webserverPort)
}
