package helper

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const (
	upperChars   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowerChars   = "abcdefghijklmnopqrstuvwxyz"
	numberChars  = "0123456789"
	specialChars = "!@#$%^&*()-_=+[]{}|;:,.<>?/~"
	allChars     = upperChars + lowerChars + numberChars + specialChars
)

// GenerateDefaultPassword generates a random password with a specified length
func GenerateDefaultPassword(length int) (string, error) {
	if length < 8 { // Ensure minimum password length is 8
		return "", fmt.Errorf("password length must be at least 8 characters")
	}

	// Ensure the password includes at least one of each character type
	password := make([]byte, length)

	// Add one of each character type
	types := []string{upperChars, lowerChars, numberChars, specialChars}
	for i, t := range types {
		char, err := randomChar(t)
		if err != nil {
			return "", err
		}
		password[i] = char
	}

	// Fill the remaining characters with random values
	for i := len(types); i < length; i++ {
		char, err := randomChar(allChars)
		if err != nil {
			return "", err
		}
		password[i] = char
	}

	// Shuffle the password to randomize the order
	shuffle(password)

	return string(password), nil
}

// randomChar selects a random character from a string
func randomChar(charSet string) (byte, error) {
	index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charSet))))
	if err != nil {
		return 0, err
	}
	return charSet[index.Int64()], nil
}

// shuffle randomly shuffles a byte slice
func shuffle(data []byte) {
	for i := range data {
		j, _ := rand.Int(rand.Reader, big.NewInt(int64(len(data))))
		data[i], data[j.Int64()] = data[j.Int64()], data[i]
	}
}
package helper

import "golang.org/x/crypto/bcrypt"

func CheckPassword(inputPassword, storedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(inputPassword))
	return err == nil
}

func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}
