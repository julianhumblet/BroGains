package webserver

import (
	"BroGains/configfile"
	"BroGains/database"
	"BroGains/hashing"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore
var sessionName string

// Constant values
const (
	methodNotAllowed   = "Method is not allowed"
	unauthorized       = "Unauthorized"
	invalidRequestBody = "Invalid request body"
	internalError      = "Internal server error"
	usernameNotFound   = "Username does not exist"
	usernameFoundOften = "Username exists more than once"
	invalidCredentials = "Invalid credentials"
	usernameExists     = "Username is already in use"
)

// Function to assign the session info to the package level variables
func SetupSessions() {

	sessionName = configfile.Configuration.Webserver.SessionName

	store = sessions.NewCookieStore([]byte(
		configfile.Configuration.Webserver.SecretSessionKey,
	))
}

func StartWebserver(port int) {

	// Create the fileservers
	fileserverLogin := http.FileServer(
		http.Dir("webserver/static/login"),
	)
	fileserverDashboard := http.FileServer(
		http.Dir("webserver/static/dashboard"),
	)

	// Host the fileservers
	http.Handle(
		"/",
		fileserverLogin,
	)
	http.Handle(
		"/dashboard/",
		isAuthenticated(
			http.StripPrefix(
				"/dashboard/",
				fileserverDashboard,
			),
		),
	)

	// Host the handlers
	http.HandleFunc("/login", loginHandler)
	http.Handle("/logout", isAuthenticated(http.HandlerFunc(logoutHandler)))
	http.HandleFunc("/register", registrationHandler)

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
		log.Printf("error making db conn on login: %s", err)
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

	// Start a session for the authenticated user
	session, err := store.Get(r, sessionName)
	if err != nil {
		log.Printf("error getting session: %s", err)
		http.Error(w, internalError, http.StatusInternalServerError)
		return
	}

	// Set session value to indicate user is authenticated
	session.Values["authenticated"] = true
	session.Values["username"] = user.Username // Store additional user data if needed
	session.Save(r, w)

	// Send HTTP code 200, login is succesful
	w.WriteHeader(http.StatusOK)

	// Redirect the user to the dashboard page
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {

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

	// Retrieve the session for the current request
	session, err := store.Get(r, sessionName)
	if err != nil {
		log.Printf("error getting session: %s", err)
		http.Error(w, internalError, http.StatusInternalServerError)
		return
	}

	// Clear the session
	session.Values["authenticated"] = false
	session.Values["username"] = nil

	// Set MaxAge to -1 to delete the session cookie after browser close
	session.Options.MaxAge = -1

	// Save the session
	err = session.Save(r, w)
	if err != nil {
		log.Printf("error saving session: %s", err)
		http.Error(w, internalError, http.StatusInternalServerError)
		return
	}

	// Send HTTP code 200, logout is succesful
	w.WriteHeader(http.StatusOK)

	// Redirect the user to the login page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func registrationHandler(w http.ResponseWriter, r *http.Request) {

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

	// Insert post data into UserRegister struct
	var user database.UserRegistration
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, invalidRequestBody, http.StatusBadRequest)
		return
	}

	// Make a connection to the database
	dbCon, err := database.OpenDBCon()
	if err != nil {
		log.Printf("error making db conn on registration: %s", err)
		http.Error(w, internalError, http.StatusInternalServerError)
		return
	}
	defer dbCon.Close()

	// Check if the username is already present in the database
	queryUserCheck := "SELECT COUNT(*) FROM accounts WHERE username = ?"
	var userCount int
	err = dbCon.QueryRow(queryUserCheck, user.Username).Scan(&userCount)
	if err != nil {
		log.Printf("error checking user existence on login: %s", err)
		http.Error(w, internalError, http.StatusInternalServerError)
		return
	}

	// Check the count to determine if the username already exists in the database
	if userCount > 0 {
		http.Error(w, usernameExists, http.StatusConflict)
		return
	}

	// Hash the password from the user
	userPasswordHash, err := hashing.HashPassword(user.Password)
	if err != nil {
		log.Println(err)
		http.Error(w, internalError, http.StatusInternalServerError)
		return
	}

	// Add the user to the database
	queryAddUserToDB := "INSERT INTO accounts (username, password) VALUES (?, ?)"
	_, err = dbCon.Exec(queryAddUserToDB, user.Username, userPasswordHash)
	if err != nil {
		log.Printf("error inserting user into database: %s", err)
		http.Error(w, internalError, http.StatusInternalServerError)
		return
	}

	// Send HTTP code 200, registration is succesful
	w.WriteHeader(http.StatusOK)

	// Redirect the user to the login page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Check if the user is logged in
func isAuthenticated(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the session for the current request
		session, err := store.Get(r, sessionName)
		if err != nil {
			log.Printf("error getting session: %s", err)
			http.Error(w, internalError, http.StatusInternalServerError)
			return
		}

		// Check if the user is authenticated
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Error(w, unauthorized, http.StatusUnauthorized)
			return
		}

		// Proceed to the next handler if the user is authenticated
		next.ServeHTTP(w, r)
	})
}

// Checking if the given given APIKey is valid
func isAPIKeyValid(apiKey string) bool {

	validApiKey := configfile.Configuration.Webserver.APIKey

	return apiKey == validApiKey
}
