# Todo-app

Простейший REST-сервис. Реализован CRUD для задач, валидация статуса.

## 🔑 Переменные окружения

Создайте файл `.env` (или экспортируйте в shell):

# DATABASE_DSN**обязательна** – файл env лежит в корне, можно прописать в переменные окружения.
DATABASE_DSN=postgres://auth_user:auth_password@localhost:5433/auth_db?sslmode=disable

# HTTP-порт
PORT=8080

# параметры для миграций
PG_DATABASE_NAME=auth_db
PG_USER=auth_user
PG_PASSWORD=auth_password
PG_PORT=5433
MIGRATION_DIR=./migrations

🚀 Запуск (локально)
bash
Копировать
Редактировать
# 1. установить deps (goose + golangci-lint)
make install-deps           
             
# 3. запустить API
make run-local              
# либо вручную:
#   set -o allexport; source .env; set +o allexport
#   go run ./cmd/api

📚 Основные Make-цели
make lint	golangci-lint на всём дереве
make migrate-up	применить все миграции
make migrate-down	откатить последнюю
make migrate-status	показать версию БД
**make run-local	миграции → запуск API без Docker**

Сервис готов к проверке: применяете миграции, задаёте DATABASE_DSN, запускаете — CRUD-ручки работают.

<img width="1280" height="995" alt="image" src="https://github.com/user-attachments/assets/6d69596a-0469-4bfc-9571-e51fdfddc21f" />

<img width="1278" height="997" alt="image" src="https://github.com/user-attachments/assets/9198aff6-9f93-4a26-a6ea-26b4d2ad8304" />

<img width="1280" height="997" alt="image" src="https://github.com/user-attachments/assets/910a262b-60ba-484e-a0eb-fb9498cef405" />

<img width="1280" height="988" alt="image" src="https://github.com/user-attachments/assets/5b88b6c7-2519-4646-8116-aff005c7f218" />

<img width="1280" height="869" alt="image" src="https://github.com/user-attachments/assets/e30de2d5-0a8c-4aaa-91a7-9e9fc8ce3a3c" />
