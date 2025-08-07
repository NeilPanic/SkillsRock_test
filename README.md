# Todo-app

Реализован CRUD для задач, валидация статуса. Простой REST сервис.
| HTTP-метод | Путь                | Описание                               |
|------------|--------------------|----------------------------------------|
| `POST`     | `/tasks`           | создать задачу                         |
| `GET`      | `/tasks`           | список \*фильтр `status`               |
| `PUT`      | `/tasks/{id}`      | частичное обновление                   |
| `DELETE`   | `/tasks/{id}`      | удалить                                |

\* `status -> {new | in_progress | done}`
## 🔑 Запуск. Переменные окружения
# .env
DATABASE_DSN=postgres://auth_user:auth_password@localhost:5432/auth_db?sslmode=disable
PORT=8080

# Примечание: 
лимит/оффсет в целях демонстрации, в проде я бы использовал другие подходы, вроде key-set пагинации. 
Или лучше использовать курсор (id параметра) и сделать where id > или < cursor limit 10/20/30/x

## 1 · PostgreSQL

```bash
# создать пользователя + базу (один раз)
sudo -u postgres psql <<'SQL'
CREATE ROLE auth_user WITH LOGIN SUPERUSER PASSWORD 'auth_password';
CREATE DATABASE auth_db OWNER auth_user;
SQL
```

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
             
# 2. запустить API          
загрузить переменные окружения
set -o allexport; source .env; set +o allexport

старт сервера
go run ./cmd/api

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
