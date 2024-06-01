package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var signedKey = []byte("test")

func CreateRefreshToken(userID string, tokenExpiry time.Duration) (string, *Payload, error) {
	return createToken(userID, tokenExpiry)
}

func CreateAccessToken(userID string, tokenExpiry time.Duration) (string, *Payload, error) {
	return createToken(userID, tokenExpiry)
}

func createToken(userID string, tokenExpiry time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(userID, tokenExpiry)
	if err != nil {
		return "", nil, err // Tambahkan penanganan kesalahan yang sesuai
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	// Buat token yang ditandatangani
	tokenString, err := token.SignedString(signedKey)
	if err != nil {
		return "", nil, err
	}

	return tokenString, payload, nil
}

func VerifyToken(tokenString string) (*Payload, error) {
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return signedKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}

	userIdInterface, ok := claims["user_id"]
	if !ok {
		return nil, fmt.Errorf("user ID claim not found in token")
	}

	userID, ok := userIdInterface.(string)
	if !ok {
		return nil, fmt.Errorf("user ID claim is not a string")
	}

	payload := &Payload{
		UserID: userID,
	}

	return payload, nil
}
