package auth

import (
	"jwt/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/exp/rand"
)

// Метод для создания Access Token
func (manager *Manager) NewAccessTokenJWT(userID string, ip string, tokenConnect string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, models.UserClaim{
		USER_ID: userID,
		IP:      ip,
		CON:     tokenConnect,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(ttl),
			},
		},
	})

	tokenString, err := token.SignedString([]byte(manager.SigningKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Метод для создания Refresh Token
func (manager *Manager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(uint64(time.Now().Unix()))
	r := rand.New(s)

	_, err := r.Read(b)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
