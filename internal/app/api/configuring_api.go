package api

import (
	"jwt/internal/app/auth"
	"jwt/storage"

	"github.com/gin-gonic/gin"
)

// Конфигурируем роутер приложения
func (api *API) configureRouterField() {
	router := gin.Default()
	router.POST("/user/auth/:GUID", api.Auth)
	router.GET("/user/refresh", api.RefreshTokens)

	api.router = router
}

// Конфигурируем менеджер токенов приложения
func (api *API) configureManagerField() error {
	manager, err := auth.NewManager(api.config.Manager)
	if err != nil {
		return err
	}

	api.manager = manager
	return nil
}

// Конфигурируем хранилище нашего приложения
func (api *API) configureStorageField() error {
	storage := storage.New(api.config.Storage)

	err := storage.Open()
	if err != nil {
		return err
	}

	api.storage = storage
	return nil
}
