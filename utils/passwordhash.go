package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPasswordMD5(password string) string {
	pwd := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(pwd, 14)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}
func CheckHash(hpass string, pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hpass), []byte(pass))
	return err == nil
}
