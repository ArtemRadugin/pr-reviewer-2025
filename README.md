# PR Reviewer Assignment Service (осень 2025)

Сервис автоматического назначения ревьюеров на Pull Request’ы внутри команды.

Решение тестового задания на позицию стажёра Backend.

## Запуск:

docker-compose up --build

После успешного старта сервис будет доступен по адресу:

http://localhost:8080

PostgreSQL поднимется внутри контейнера и автоматически прогонит миграции.

## Стек
- Go 1.25 
- Chi — роутинг HTTP-запросов 
- Viper — конфигурации 
- PostgreSQL — основная БД 
- sqlx — работа с SQL 
- Docker + Docker Compose 
- migrate/migrate — система миграций 
- Makefile

## План расширения (если нужно)
- Добавить авторизацию (JWT)
- Добавить очередь задач (RabbitMQ)
- Добавить метрики (Prometheus)
- Покрыть сервис unit-тестами
- Добавить линтеры (golangci-lint)