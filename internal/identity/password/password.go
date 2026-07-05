package password

import "golang.org/x/crypto/bcrypt"

func Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 16)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
