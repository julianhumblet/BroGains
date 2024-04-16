package hashing

import (
	"encoding/base64"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %s", err)
	}

	hashedPassword := base64.StdEncoding.EncodeToString(hashedBytes)

	return hashedPassword, nil
}

func CheckHashPassword(password, hash string) bool {

	decodedHash, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		log.Printf("error decoding string: %s", err)
		return false
	}

	err = bcrypt.CompareHashAndPassword(decodedHash, []byte(password))

	return err == nil
}
