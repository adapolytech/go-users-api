package utils

import "golang.org/x/crypto/bcrypt"

type PasswordUtils struct {
	Cost int
}

func (u *PasswordUtils) Hash(plainText string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), u.Cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (*PasswordUtils) Compare(hash string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func NewPasswordUtils(cost int) *PasswordUtils {
	return &PasswordUtils{Cost: cost}
}
