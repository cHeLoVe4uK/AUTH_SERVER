package auth

import (
	"errors"
	"jwt/internal/app/config"
)

// Менеджер для работы с токенами в приложении
type Manager struct {
	SigningKey string // секретный код, которым будут подписываться токены
}

// Конструктор, возвращающий экземпляр Manager
func NewManager(config *config.ManagerConfig) (*Manager, error) {
	if config.SigningKey == "" {
		return nil, errors.New("empty SigningKey")
	}
	return &Manager{SigningKey: config.SigningKey}, nil
}
