package auth

import (
	"encoding/base64"
	"errors"
	"fmt"
	"jwt/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Метод для верификации токена доступа
func (manager *Manager) ParseAccessToken(accessToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) { return manager.SigningKey, nil })
	if err != nil {
		return nil, err
	}
	_, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	return token, nil
}

// Метод для верификации рефреш токена
func (manager *Manager) ParseRefreshToken(refreshToken string, userDB *models.User) error {
	// Декодируем выданный пользователю рефреш токен
	decryptedRefreshToken, err := base64.StdEncoding.DecodeString(refreshToken)
	if err != nil {
		return fmt.Errorf("cannot decrypt refresh token: %s", err)
	}

	// Сравниваем раскодированный рефреш токен и реыреш токен из БД
	err = bcrypt.CompareHashAndPassword([]byte(userDB.REFRESH_TOKEN), decryptedRefreshToken)
	if err != nil {
		return fmt.Errorf("refresh token provided by user and refresh token from db is not equal: %s", err)
	}

	// Если рефреш токен предоставленный пользователем и рефреш токен из БД одни и теже переходим к проверке его на доступное время
	exp, err := time.Parse(time.RFC3339, userDB.EXPIRATION_TIME)
	if err != nil {
		return fmt.Errorf("cannot parse expiration time of refresh token from DB: %s", err)
	}
	if exp.Unix() < time.Now().Unix() {
		return errors.New("expiration time of refresh token run out")
	}

	// Если все прошло хорошо переходим к последней проверке на то, был ли использован токен или нет
	if userDB.USED_AT != "" {
		return errors.New("this refresh token was used")
	}

	return nil
}
