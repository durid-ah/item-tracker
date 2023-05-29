package helpers

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	byte, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Fatalln("Unable to hash password")
	}

	return string(byte)
}

func ValidatePassword(password string, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	return err == nil
}
