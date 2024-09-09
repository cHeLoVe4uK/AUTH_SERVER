package api

import (
	"encoding/base64"
	"jwt/internal/app/headers"
	"jwt/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Хэндлер, отвечающий за работу приложения по пути "/user/refresh" (отвечает за сброс токенов и выдачу новых)
func (api *API) RefreshTokens(ctx *gin.Context) {
	log.Println("User do 'POST:RefreshTokens /user/refresh'")

	// Проверяем что Auth Header с токеном доступа был передан как полагается
	accessHeaderToken, err := headers.VerifiesHeader(ctx)
	if err != nil {
		log.Println("Trouble with verify header: ", err)
		message := responseMessage{
			Message: "Sorry, you provide incorrect Authorization Header",
		}
		ctx.JSON(http.StatusBadRequest, message)
		return
	}

	// Если все прошло хорошо валидируем полученный токен
	accessToken, err := api.manager.ParseAccessToken(accessHeaderToken)
	if err != nil {
		log.Println("Trouble with parsing Access Token: ", err)
		message := responseMessage{
			Message: "Sorry, trouble with server. Try later",
		}
		ctx.JSON(http.StatusInternalServerError, message)
		return
	}

	// Если токен доступа прошел валидацию переходим к обработке рефреш токена
	var user models.User
	err = ctx.ShouldBindJSON(&user)
	if err != nil {
		log.Println("User provide incorrect json in body request: ", err)
		message := responseMessage{
			Message: "Sorry, you provide incorrect refresh token",
		}
		ctx.JSON(http.StatusBadRequest, message)
		return
	}
	if user.REFRESH_TOKEN == "" {
		log.Println("User provide empty refresh token: ", err)
		message := responseMessage{
			Message: "Sorry, you provide empty refresh token",
		}
		ctx.JSON(http.StatusBadRequest, message)
		return
	}

	// Если все прошло хорошо у нас есть и токен доступа и рефреш токен
	// Получаем полезные данные из токена
	claims, ok := accessToken.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("Cannot recieve claims from token")
		message := responseMessage{
			Message: "Sorry, trouble with server. Try later",
		}
		ctx.JSON(http.StatusInternalServerError, message)
		return
	}
	userID := claims["USER_ID"]
	tokenConnect := claims["CON"]
	tokenIP := claims["IP"]

	// Теперь получаем пользователя из БД с необходимым guid и найдем связанный с ним рефреш токен
	userDB, err := api.storage.User().GetUser(userID.(string), tokenConnect.(string))
	if err != nil {
		log.Println("Cannot get user from DB: ", err)
		message := responseMessage{
			Message: "Sorry, trouble with access to DB. Try later",
		}
		ctx.JSON(http.StatusInternalServerError, message)
		return
	}

	// Если все прошло хорошо значит guid и token_connect из токена доступа уже совпадают с таковыми у пользователя из БД, поэтому занимаемся только проверкой рефреш токена
	err = api.manager.ParseRefreshToken(user.REFRESH_TOKEN, userDB)
	if err != nil {
		log.Println("Trouble with parsing Refresh Token: ", err)
		message := responseMessage{
			Message: "Sorry, trouble with server. Try later",
		}
		ctx.JSON(http.StatusInternalServerError, message)
		return
	}

	// Если все прошло хорошо устанавливаем использованному рефреш токену статус использованный
	err = api.storage.User().SetUserColumnUsedAt(userDB.USER_ID, userDB.TOKEN_CONNECT)
	if err != nil {
		log.Println("Trouble with set Refresh Token used status: ", err)
		message := responseMessage{
			Message: "Sorry, trouble with access to DB. Try later",
		}
		ctx.JSON(http.StatusInternalServerError, message)
		return
	}

	// И только теперь выдаем пользователю 2 новых токена
	// Подготавливаем новую запись о пользователе для выдачи новых токенов
	var newUser models.User
	newUser.USER_ID = userDB.USER_ID
	newUser.TOKEN_CONNECT = uuid.NewString()
	newAccessToken, err := api.manager.NewAccessTokenJWT(newUser.USER_ID, ctx.ClientIP(), newUser.TOKEN_CONNECT, api.config.AccessTokenLive)
	if err != nil {
		log.Println("Trouble with creating Access Token: ", err)
		message := responseMessage{
			Message: "Sorry, trouble with server. Try later",
		}
		ctx.JSON(http.StatusInternalServerError, message)
		return
	}
	newRefreshToken, err := api.manager.NewRefreshToken()
	if err != nil {
		log.Println("Trouble with creating Refresh Token: ", err)
		message := responseMessage{
			Message: "Sorry, trouble with server. Try later",
		}
		ctx.JSON(http.StatusInternalServerError, message)
		return
	}
	newUser.CREATED_AT = time.Now().Format(time.RFC3339)
	newUser.EXPIRATION_TIME = time.Now().Add(api.config.RefreshTokenLive).Format(time.RFC3339)

	// Если токены выбились успешно шифруем рефреш токен
	newEncryptedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(newRefreshToken), 14)
	if err != nil {
		log.Println("Trouble with encrypt Refresh Token: ", err)
		message := responseMessage{
			Message: "Sorry, trouble with server. Try later",
		}
		ctx.JSON(http.StatusInternalServerError, message)
		return
	}
	newUser.REFRESH_TOKEN = string(newEncryptedRefreshToken)

	// Создаем новую запись для пользователя в БД
	_, err = api.storage.User().CreateUser(&newUser)
	if err != nil {
		log.Println("Trouble while creating user in DB: ", err)
		message := responseMessage{
			Message: "Sorry, trouble with access to DB. Try later",
		}
		ctx.JSON(http.StatusInternalServerError, message)
		return
	}
	log.Printf("User add in table (%s): %v", api.config.Storage.UserTable, newUser)

	// Возвращаем новые токены
	newRefreshTokenBase64 := base64.StdEncoding.EncodeToString([]byte(newRefreshToken))
	message := responseMessage{
		Message: "Request successfully done",
		responseToken: responseToken{
			AccessToken:  newAccessToken,
			RefreshToken: newRefreshTokenBase64,
		},
	}
	ctx.JSON(http.StatusOK, message)

	// В конце проверяем совпадают ли IP в полученном токене доступа и текущий пользователя, если нет, отправляем предупреждение на почту)
	if tokenIP.(string) != ctx.ClientIP() {
		api.sendMessageOnEmail()
		// Логика обработки ошибки
	}
}
