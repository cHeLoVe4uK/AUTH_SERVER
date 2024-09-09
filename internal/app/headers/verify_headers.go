package headers

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

// Функция для проверки, что токен доступа, переданный пользователем в Authorization Header, удовлетворяет запросам
func VerifiesHeader(c *gin.Context) (string, error) {
	// Считываем заголовок
	accessToken := c.GetHeader("Authorization")
	// Если пустой возвращаем ошибку
	if accessToken == "" {
		return "", errors.New("empty auth header")
	}
	// Проверяем верный ли формат токена предоставлен
	accessTokenParts := strings.Split(accessToken, " ")
	if len(accessTokenParts) != 2 || accessTokenParts[0] != "Bearer" {
		return "", errors.New("invalid header format")
	}
	if accessTokenParts[1] == "" {
		return "", errors.New("token is empty")
	}
	return accessTokenParts[1], nil
}
