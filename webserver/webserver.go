package webserver

import (
	"log"
	"net/http"
	"strconv"
)

// Constant values
const (
	methodNotAllowed = "Method is not allowed"
)

func StartWebserver(port int) {

	// Serve the static files
	http.Handle("/",
		http.FileServer(
			http.Dir("webserver/static"),
		),
	)

	// Host the handlers
	http.HandleFunc("/login", loginHandler)

	// Start the webserver
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Fatalf("error starting the webserver: %s", err)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	// Add login handling

	// Check if the made request is POST
	if r.Method != http.MethodPost {
		http.Error(w, methodNotAllowed, http.StatusMethodNotAllowed)
		return
	}
}
