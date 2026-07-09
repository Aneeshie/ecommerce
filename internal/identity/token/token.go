package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/Aneeshie/ecommerce/internal/identity/domain"
	"github.com/golang-jwt/jwt/v5"
)

func (m *Manager) GenerateAccessToken(user domain.User, AccessTokenTTL time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub":  user.ID.String(),
		"role": user.Role,
		"exp":  time.Now().Add(AccessTokenTTL).Unix(),
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err := token.SignedString(m.secret)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func (m *Manager) GenerateRefreshToken() (string, error) {
	token := make([]byte, 32)

	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}

	base64String := base64.RawURLEncoding.EncodeToString(token)

	return base64String, nil
}

func (m *Manager) HashRefreshToken(token string) string {
	bytesResult := sha256.Sum256([]byte(token))

	hexResult := hex.EncodeToString(bytesResult[:])

	return hexResult
}

func (m *Manager) VerifyAccessToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secret, nil

	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
