package config

import "time"

// Конфиг для приложения
type ConfigAPI struct {
	BindAddr         string         `toml:"bind_addr"`          // порт, на котором будет работать приложение
	AccessTokenLive  time.Duration  `toml:"acces_token_live"`   // время жизни токена доступа
	RefreshTokenLive time.Duration  `toml:"refresh_token_live"` // время жизни рефреш токена
	Storage          *StorageConfig // конфиг для БД
	Manager          *ManagerConfig // конфиг для менеджера токенов
}

// Конструктор конфига для приложения
func NewConfigAPI(storage *StorageConfig, manager *ManagerConfig) *ConfigAPI {
	return &ConfigAPI{
		Storage: storage,
		Manager: manager,
	}
}

// Конфиг для менеджера токенов
type ManagerConfig struct {
	SigningKey string `toml:"secret_key"` // секретный код, которым будут подписываться токены
}

// Конструктор конфига для менеджера токенов
func NewManagerConfig() *ManagerConfig {
	return &ManagerConfig{}
}

// Конфиг для БД
type StorageConfig struct {
	DataBaseURI string `toml:"database_uri"` // строка, содержащая параметры для подключения к БД
	DriverName  string `toml:"driver_name"`  // имя драйвера для работы с БД
	UserTable   string `toml:"user_table"`   // название таблицы для работы с пользователями
}

// Конструктор конфига для БД
func NewStorageConfig() *StorageConfig {
	return &StorageConfig{}
}
