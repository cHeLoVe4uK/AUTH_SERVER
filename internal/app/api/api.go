package api

import (
	"jwt/internal/app/auth"
	"jwt/internal/app/config"
	"jwt/storage"
	"log"

	"github.com/gin-gonic/gin"
)

// Само приложение (сервер)
type API struct {
	config  *config.ConfigAPI // конфиг, который будет использоваться для настройки работы приложения
	storage *storage.Storage  // поле, через которое будет осуществляться работа с БД
	router  *gin.Engine       // поле, которое отвечает за роутинг в приложении
	manager *auth.Manager     // поле, через которое будет осуществляться работа с токенами
}

// Конструктор, возвращающий экземпляр API (приложения)
func NewAPI(config *config.ConfigAPI) *API {
	return &API{
		config: config,
	}
}

// Метод, конфигурирующий и запускающий сервер
func (api *API) Start() error {
	err := api.configureManagerField()
	if err != nil {
		return err
	}
	log.Println("Manager succsessfully configured")

	api.configureRouterField()
	log.Println("Router succsessfully configured")

	err = api.configureStorageField()
	if err != nil {
		return err
	}
	log.Println("Storage succsessfully configured")

	log.Printf("Starting on port: %v", api.config.BindAddr)

	err = api.router.Run(":" + api.config.BindAddr)
	if err != nil {
		log.Fatalf("While server is working: %v", err)
	}

	return nil
}
