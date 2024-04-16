package hashing

import (
	"encoding/base64"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func CheckHashPassword(password, hash string) bool {

	decodedHash, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		log.Printf("error decoding string: %s", err)
		return false
	}

	err = bcrypt.CompareHashAndPassword(decodedHash, []byte(password))

	return err == nil
}
