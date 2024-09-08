package api

// Просто структуры для более удобного ответа пользователю в запросах
type responseMessage struct {
	Message string `json:"message"` // сообщение для пользователя
	responseToken // токены пользователя
}

type responseToken struct {
	AccessToken  string `json:"access_token,omitempty" `
	RefreshToken string `json:"refresh_token,omitempty"`
}
