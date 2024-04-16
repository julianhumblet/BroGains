package database

import (
	"BroGains/configfile"
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRegistration struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func OpenDBCon() (*sql.DB, error) {

	// Database connection creds
	dbConCreds := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		configfile.Configuration.Database.Username,
		configfile.Configuration.Database.Password,
		configfile.Configuration.Database.Host,
		strconv.Itoa(configfile.Configuration.Database.Port),
		configfile.Configuration.Database.DatabaseName,
	)

	// Open a new database connection
	dbCon, err := sql.Open("mysql", dbConCreds)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %s", err)
	}

	// Ping the connection to the database
	err = dbCon.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %s", err)
	}

	return dbCon, nil
}
