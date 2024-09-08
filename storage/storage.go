package storage

import (
	"database/sql"
	"jwt/internal/app/config"

	_ "github.com/lib/pq"
)

// Хранилище для приложения
type Storage struct {
	// Поля неэкспортируемые (конфендициальная информация)
	config         *config.StorageConfig // параметры для работы БД
	db             *sql.DB               // сама БД
	userRepository *UserRepository       // модельный репозиторий для работы с пользователями
}

// Конструктор, возвращающий БД
func New(config *config.StorageConfig) *Storage {
	return &Storage{
		config: config,
	}
}

// Метод, открывающий соединение между нашим приложением и БД
func (storage *Storage) Open() error {
	db, err := sql.Open(config.NewStorageConfig().DriverName, config.NewStorageConfig().DataBaseURI)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}

	storage.db = db
	return nil
}

// Метод, закрывающий наше соединение с БД
func (storage *Storage) Close() {
	storage.db.Close()
}

// Метод, создающий репозиорий для работы с пользователями
func (storage *Storage) User() *UserRepository {
	if storage.userRepository != nil {
		return storage.userRepository
	}

	storage.userRepository = &UserRepository{
		storage: storage,
	}

	return storage.userRepository
}
