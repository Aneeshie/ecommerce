package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"

	"github.com/Aneeshie/ecommerce/internal/identity/domain"
	"github.com/golang-jwt/jwt/v5"
)

func (m *Manager) GenerateAccessToken(user domain.User, AccessTokenTTL time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": user.ID.String(),
		"role": user.Role,
		"exp": time.Now().Add(AccessTokenTTL).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err := token.SignedString(m.secret)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func (m *Manager) GenerateRefreshToken() (string, error){
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
