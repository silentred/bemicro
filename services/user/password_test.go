package user

import (
	"fmt"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestBcrypt(t *testing.T) {
	password := []byte("111111")

	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(hashedPassword))

	// Comparing the password with the hash
	hashedPassword = []byte("$2y$10$zPJu5NXsWfc0ctGv0xXGZOCQreXjOe4dMoIJHQcfEwCEFFe6Jc5ey")
	err = bcrypt.CompareHashAndPassword(hashedPassword, password)
	fmt.Println(err) // nil means it is a match
}
