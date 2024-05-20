
# User Balance Service

Этот микросервис предназначен для работы с балансом пользователей. Он предоставляет HTTP API для зачисления и списания средств, перевода средств между пользователями, а также получения баланса пользователей. Все запросы и ответы обрабатываются в формате JSON.

## Требования

- Язык программирования: Golang
- Фреймворк: Gin
- База данных: PostgreSQL
- Docker и docker-compose для развертывания dev-среды

## Установка

1. Установите Golang: https://golang.org/doc/install
2. Установите Docker: https://docs.docker.com/get-docker/
3. Клонируйте репозиторий:

   ```bash
   git clone https://github.com/ваш_пользователь/user-balance-service.git
   ```

4. Перейдите в каталог проекта:

   ```bash
   cd user-balance-service
   ```

5. Запустите контейнеры Docker:

   ```bash
   docker-compose up -d
   ```

6. Инициализируйте базу данных:

   ```bash
   docker exec -i $(docker-compose ps -q db) psql -U user -d user_balance_db < init.sql
   ```

## Использование

### Запуск сервера

1. Запустите сервер:

   ```bash
   go run main.go
   ```

### HTTP API

- `POST /balance/credit`: Зачисление средств на баланс пользователя.
- `POST /balance/debit`: Списание средств с баланса пользователя.
- `POST /balance/transfer`: Перевод средств между пользователями.
- `GET /balance/:userId`: Получение баланса пользователя.

Подробнее о параметрах и формате запросов/ответов смотрите в коде и в Swagger файле.
