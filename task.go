Используемые технологии:
* GO
* JWT
* PostgreSQL

Задание:
Написать часть сервиса аутентификации
Для REST маршрута:
* Первый маршрут выдает пару Access, Refresh токенов для пользователя с идентификатором (GUID) указанным в параметре запроса
* Второй маршрут выполняет Refresh операцию на пару Access, Refresh токенов

Требования:
Access токен тип JWT, алгоритм SHA512, хранить в базе строго запрещено.

Refresh токен тип произвольный, формат передачи base64, хранится в базе исключительно в виде bcrypt хеша, должен быть защищен от изменения на стороне клиента и попыток повторного использования.

Access, Refresh токены обоюдно связаны, Refresh операцию для Access токена можно выполнить только тем Refresh токеном, который был выдан вместе с ним.

Payload токенов должен содержать сведения об ip адресе клиента, которому он был выдан. В случае, если ip адрес изменился, при рефреш операции нужно послать email warning на почту юзера (для упрощения можно использовать моковые данные).

Результат:
Результат выполнения задания нужно предоставить в виде исходного кода на GitHub. Будет плюсом, если получится использовать Docker и покрыть код тестами.
