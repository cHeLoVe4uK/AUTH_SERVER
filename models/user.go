package models

import (
	"github.com/golang-jwt/jwt/v5"
)

// Модель пользователя, которая будет использоваться в приложении и для работы с БД
type User struct {
	USER_ID         string // guid пользователя
	TOKEN_CONNECT   string // рандоманый uuid, который будет использоваться для создания связи между 2 токенами
	REFRESH_TOKEN   string `json:"refresh_token"` // рандомная строка представляющая собой рефреш токен
	CREATED_AT      string // время создания рефреш токена
	EXPIRATION_TIME string // время когда рефреш токен перестанет работать
	USED_AT         string // по этому полю будет отслеживаться использовался ли рефреш токен (просто временная галочка)
}

// Модель для работы с полезной нагрузкой в токенах
type UserClaim struct {
	USER_ID              string // guid пользователя
	IP                   string // IP пользователя
	CON                  string // рандоманый uuid, который будет использоваться для создания связи между 2 токенами
	jwt.RegisteredClaims        // стандартные claims для токенов
}
