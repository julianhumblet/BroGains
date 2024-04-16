package webserver

import (
	"BroGains/configfile"
	"BroGains/database"
	"BroGains/hashing"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// Constant values
const (
	methodNotAllowed   = "Method is not allowed"
	unauthorized       = "Unauthorized"
	invalidRequestBody = "Invalid request body"
	internalError      = "Internal server error"
	usernameNotFound   = "Username does not exist"
	usernameFoundOften = "Username exists more than once"
	invalidCredentials = "Invalid credentials"
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

	// Check if the made request is POST
	if r.Method != http.MethodPost {
		http.Error(w, methodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	// Check if the given API key is valid
	enteredApiKey := r.Header.Get("Authorization")
	if !isAPIKeyValid(enteredApiKey) {
		http.Error(w, unauthorized, http.StatusUnauthorized)
		return
	}

	// Insert post data into UserLogin struct
	var user database.UserLogin
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, invalidRequestBody, http.StatusBadRequest)
		return
	}

	// Make a connection to the database
	dbCon, err := database.OpenDBCon()
	if err != nil {
		log.Printf("error making db conn on login: ", err)
		http.Error(w, internalError, http.StatusInternalServerError)
		return
	}
	defer dbCon.Close()

	// Check if the username is present in the database
	queryUserCheck := "SELECT COUNT(*) FROM accounts WHERE username = ?"
	var userCount int
	err = dbCon.QueryRow(queryUserCheck, user.Username).Scan(&userCount)
	if err != nil {
		log.Printf("error checking user existence on login: %s", err)
		http.Error(w, internalError, http.StatusInternalServerError)
		return
	}

	if userCount == 0 {
		http.Error(w, usernameNotFound, http.StatusConflict)
		return
	} else if userCount > 1 {
		http.Error(w, usernameFoundOften, http.StatusConflict)
		return
	}

	// Check if the password from user is correct
	queryPasswordCheck := "SELECT password FROM accounts WHERE username = ?"
	var userPasswordHash string
	err = dbCon.QueryRow(queryPasswordCheck, user.Username).Scan(&userPasswordHash)
	if err != nil {
		log.Printf("error getting hashed password from db: %s", err)
		http.Error(w, internalError, http.StatusInternalServerError)
		return
	}

	if !hashing.CheckHashPassword(user.Password, userPasswordHash) {
		http.Error(w, invalidCredentials, http.StatusUnauthorized)
		return
	}

	// Send HTTP code 200, login is succesful
	w.WriteHeader(http.StatusOK)
}

// Checking if the given given APIKey is valid
func isAPIKeyValid(apiKey string) bool {

	validApiKey := configfile.Configuration.Webserver.APIKey

	return apiKey == validApiKey
}
