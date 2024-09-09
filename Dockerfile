FROM ubuntu:24.04
RUN ["mkdir","/serv_dir"]
RUN ["mkdir","/serv_dir/api_config"]
WORKDIR /serv_dir
COPY /api_config/api.toml ./api_config
COPY auth-server .
ENTRYPOINT ["./auth-server"]
# Для работы требуется иметь предварительно запущенный контейнер PostrgreSQL (либо же просто запущенный PostgreSQL) со всеми настройками, что используется в конфиге для приложения.
