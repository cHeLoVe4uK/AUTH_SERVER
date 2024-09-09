package api

import (
	"encoding/base64"
	"jwt/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Хэндлер, отвечающий за работу приложения по пути "/user/auth/:GUID" (отвечает за выдачу токенов пользователю)
func (api *API) Auth(ctx *gin.Context) {
	log.Println("User do 'POST:Auth /user/auth/:GUID'")

	// Считываем значение из запроса и проверяем верно ли оно по формату
	guid := ctx.Params.ByName("GUID")
	_, err := uuid.Parse(guid)
	if err != nil {
		log.Println("User provide incorrected GUID: ", err)
		message := responseMessage{
			Message: "Sorry, you provide incorrected value of GUID",
		}
		ctx.JSON(http.StatusBadRequest, message)
		return
	}

	// Создаем сущность пользователя и заполняем его
	var user models.User
	user.USER_ID = guid
	user.TOKEN_CONNECT = uuid.NewString()

	// Создаем токены
	accessToken, err := api.manager.NewAccessTokenJWT(guid, ctx.ClientIP(), user.TOKEN_CONNECT, api.config.AccessTokenLive)
	if err != nil {
		log.Println("Trouble with creating Access Token: ", err)
		message := responseMessage{
			Message: "Sorry, trouble with server. Try later",
		}
		ctx.JSON(http.StatusInternalServerError, message)
		return
	}
	refreshToken, err := api.manager.NewRefreshToken()
	if err != nil {
		log.Println("Trouble with creating Refresh Token: ", err)
		message := responseMessage{
			Message: "Sorry, trouble with server. Try later",
		}
		ctx.JSON(http.StatusInternalServerError, message)
		return
	}

	user.CREATED_AT = time.Now().Format(time.RFC3339)
	user.EXPIRATION_TIME = time.Now().Add(api.config.RefreshTokenLive).Format(time.RFC3339)

	// Если токены выбились успешно шифруем рефреш токен
	encryptedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), 14)
	if err != nil {
		log.Println("Trouble with encrypt Refresh Token: ", err)
		message := responseMessage{
			Message: "Sorry, trouble with server. Try later",
		}
		ctx.JSON(http.StatusInternalServerError, message)
		return
	}
	user.REFRESH_TOKEN = string(encryptedRefreshToken)

	// Создаем пользователя в БД (или добавляем еще одну запись)
	_, err = api.storage.User().CreateUser(&user)
	if err != nil {
		log.Println("Trouble while creating user in DB: ", err)
		message := responseMessage{
			Message: "Sorry, trouble with access to DB. Try later",
		}
		ctx.JSON(http.StatusInternalServerError, message)
		return
	}
	log.Printf("User add in table (%s): %v", api.config.Storage.UserTable, user)

	// Если пользователь был успешно добавлен в БД возвращаем ему токены
	refreshTokenBase64 := base64.StdEncoding.EncodeToString([]byte(refreshToken))
	message := responseMessage{
		Message: "Request successfully done",
		responseToken: responseToken{
			AccessToken:  accessToken,
			RefreshToken: refreshTokenBase64,
		},
	}
	ctx.JSON(http.StatusOK, message)
}
