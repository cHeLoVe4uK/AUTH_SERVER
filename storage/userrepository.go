package storage

import (
	"fmt"
	"jwt/models"
	"log"
	"time"
)

// Модельный репозиторий
type UserRepository struct {
	storage *Storage // хранит в себе БД, т.к. общение с БД реализовано посредством репозитория (небольшое замыкание)
}

// Метод для добавления пользователя в БД
func (u *UserRepository) CreateUser(user *models.User) (*models.User, error) {
	query := fmt.Sprintf("INSERT INTO %s (USER_ID, TOKEN_CONNECT, REFRESH_TOKEN, CREATED_AT, EXPIRATION_TIME, USED_AT) VALUES ($1, $2, $3, $4, $5, $6)", u.storage.config.UserTable)
	fmt.Println(query)
	_, err := u.storage.db.Exec(query, user.USER_ID, user.TOKEN_CONNECT, user.REFRESH_TOKEN, user.CREATED_AT, user.EXPIRATION_TIME, user.USED_AT)
	if err != nil {
		log.Println("An error occurred when calling the method CreateUser")
		return nil, err
	}
	return user, nil
}

// Метод для получения пользователя из БД
func (u *UserRepository) GetUser(guid string, connect string) (*models.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE USER_ID=$1 AND TOKEN_CONNECT=$2", u.storage.config.UserTable)

	var user models.User

	err := u.storage.db.QueryRow(query, guid, connect).Scan(&user.USER_ID, &user.TOKEN_CONNECT, &user.REFRESH_TOKEN, &user.CREATED_AT, &user.EXPIRATION_TIME, &user.USED_AT)
	if err != nil {
		log.Println("An error occurred when calling the method GetUser")
		return nil, err
	}

	return &user, nil
}

// Метод для установки у рефреш токена пользователя в БД статус на использованный
func (u *UserRepository) SetUserColumnUsedAt(guid string, connect string) error {
	query := fmt.Sprintf("UPDATE %s SET USED_AT=$1 WHERE USER_ID=$2 AND TOKEN_CONNECT=$3", u.storage.config.UserTable)

	_, err := u.storage.db.Exec(query, time.Now().Format(time.RFC3339), guid, connect)
	if err != nil {
		log.Println("An error occurred when calling the method SetUserColumnUsedAt")
		return err
	}

	return nil
}
