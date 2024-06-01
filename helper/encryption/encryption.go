package encryption

import "golang.org/x/crypto/bcrypt"

func EncryptWithHash(password []byte) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func VerifyPassword(hashedPassword string, password []byte) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), password)
}
