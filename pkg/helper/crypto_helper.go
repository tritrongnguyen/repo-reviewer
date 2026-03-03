package helper

import "golang.org/x/crypto/bcrypt"

func HashPassword(p string) (string, error) {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(p), 14)
	return string(bytes), nil
}

func CheckPasswordHash(p, h string) bool {
	return bcrypt.CompareHashAndPassword([]byte(h), []byte(p)) == nil
}
