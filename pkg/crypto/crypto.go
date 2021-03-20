package crypto

import "golang.org/x/crypto/bcrypt"

// HashPassword encrypts a given password using bcrypt algorithm.
func HashPassword(password string, cost int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// CheckPasswordHash checks if the given passwords matches.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
