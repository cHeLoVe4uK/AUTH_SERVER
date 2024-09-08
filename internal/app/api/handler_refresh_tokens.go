package api

import (
	"jwt/internal/app/headers"
	"jwt/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Хэндлер, отвечающий за работу приложения по пути "/user/refresh" (отвечает за сброс токенов и выдачу новых)
func (api *API) RefreshTokens(ctx *gin.Context) {
	log.Println("User do 'POST:RefreshTokens /user/refresh'")

	// Проверяем что Auth Header с токеном доступа был передан как надо
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
			Message: "Sorry, you provide incorrect data format",
		}
		ctx.JSON(http.StatusBadRequest, message)
		return
	}
	// Если все прошло хорошо у нас есть и токен доступа и рефреш токен
	// Получаем поллезные данные из токена
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
	tokenConnect := claims["TOKEN_CONNECT"]

	// Теперь получаем пользователя из БД с необходимым guid и найдем связанный с ним рефреш токен
	userDB, err := api.storage.User().GetUser(userID.(string), tokenConnect.(string))
	if err != nil {
		log.Println("Cannot get user from DB: ", err)
		message := responseMessage{
			Message: "Sorry, trouble with server. Try later",
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
	err = api.storage.User().SetUserColumnUsedAt(userID.(string), tokenConnect.(string))
	if err != nil {
		log.Println("Trouble with set Refresh Token used status: ", err)
		message := responseMessage{
			Message: "Sorry, trouble with server. Try later",
		}
		ctx.JSON(http.StatusInternalServerError, message)
		return
	}

	// И только теперь выдаем пользователю 2 новых токена

}
