package config

import "github.com/BurntSushi/toml"

// Функция создающая все конфиги и парсящая в них все значения необходимые для работы приложения
func AllConfigSetup(configPathApi *string) (*ConfigAPI, error) {
	// Конфиг для менеджера
	managerConfig := NewManagerConfig()
	_, err := toml.DecodeFile(*configPathApi, managerConfig)
	if err != nil {
		return nil, err
	}

	// Конфиг для БД
	storageConfig := NewStorageConfig()
	_, err = toml.DecodeFile(*configPathApi, storageConfig)
	if err != nil {
		return nil, err
	}

	// Конфиг для приложения
	apiConfig := NewConfigAPI(storageConfig, managerConfig)
	_, err = toml.DecodeFile(*configPathApi, apiConfig)
	if err != nil {
		return nil, err
	}

	return apiConfig, nil
}
