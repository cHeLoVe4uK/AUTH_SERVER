package main

import (
	"flag"
	"jwt/internal/app/api"
	"jwt/internal/app/config"
	"log"
)

// Переменная, в которой хранится путь до файла с конфигом
var configPathApi string

// Функция определяющая новый флаг для приложения (путь до конфига)
func init() {
	flag.StringVar(&configPathApi, "path", "config/api.toml", "path to config file in .toml format")
}

func main() {
	// Парсим конфиг для приложения, путь был указан пользователем через cmd или использован по-умолчанию
	flag.Parse()

	// Настраиваем конфиги для приложения
	apiConfig, err := config.AllConfigSetup(&configPathApi)
	if err != nil {
		log.Printf("Configure file not found, server will not be started: %v", err)
		return
	}

	// Создаем сервер аутентификации
	server := api.NewAPI(apiConfig)

	// Запускаем работу сервера
	err = server.Start()
	if err != nil {
		log.Printf("An error occurred when starting server: %v", err)
	}
}
