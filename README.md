# PR Reviewer Assignment Service

Сервис для автоматического назначения ревьюверов на Pull Request'ы.

## Запуск 

### Docker Compose (Сервис + PostgreSQL)

```bash
make docker-up
# или
docker-compose up -d
```

Сервис доступен на `http://localhost:8080`

### Локальный запуск

1. Скопируйте `.env.example` в `.env`:
```bash
cp .env.example .env
```

2. Запустите PostgreSQL:
```bash
docker-compose up -d postgres
```

3. Запустите приложение:
```bash
make run
# или
go run ./cmd/api
```

## Переменные окружения

Создайте `.env` файл на основе `.env.example`:

```env
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=avito_user
DB_PASSWORD=avito_password
DB_NAME=avito_db
DB_SSLMODE=disable
```

## Makefile команды

```bash
make build              # Собрать приложение
make run                # Запустить локально
make test-unit          # Unit-тесты (handlers, domain)
make test-integration   # Integration-тесты (services, требует БД)
make test-all           # Все тесты
make test-coverage      # Покрытие тестами + HTML отчет
make generate-mocks     # Сгенерировать моки для тестов
make fmt                # Отформатировать код
make lint               # Линтер (golangci-lint)
make docker-up          # Запустить Docker Compose
make docker-down        # Остановить Docker Compose
make loadtest-burst     # Burst нагрузочное тестирование
make loadtest-rampup    # Ramp-up нагрузочное тестирование
make loadtest-all       # Все нагрузочные тесты
```

## Реализовано

### Основные требования
- Все эндпоинты из openapi.yml
- Автоматическое назначение до 2 ревьюверов из команды автора
- Переназначение ревьювера из команды заменяемого
- Запрет изменения после MERGED
- Идемпотентность операций Merge
- Транзакции и консистентность данных (не требовалось, но я сделал всё конкурентно безопасным)
- Индексы БД (не требовалось, но я сделал)

### Тестирование
- Unit-тесты (handlers, domain)
- Integration-тесты (services) 
- См. покрытие make test-coverage 
- Нагрузочное тестирование (uniform, burst, ramp-up)
- Все SLI требования выполнены (99.9% success rate, ≤300ms avg response)

## API Endpoints

- `POST /team/add` - Создать команду
- `GET /team/get?team_name=...` - Получить команду
- `POST /team/deactivate` - Деактивировать команду (доп. задание)
- `POST /users/setIsActive` - Установить флаг активности
- `GET /users/getReview?user_id=...` - Получить PR'ы пользователя
- `POST /pullRequest/create` - Создать PR и назначить ревьюверов
- `POST /pullRequest/merge` - Пометить PR как MERGED
- `POST /pullRequest/reassign` - Переназначить ревьювера
- `GET /stats` - Получить статистику (доп. задание)

Полная спецификация: `openapi.yml` (кроме добавленных эндпоинтов)

## Документация

- `docs/DECISIONS.md` - Принятые мной решения 
- `docs/schema.dbml` - Схема базы данных

## Дополнительные задания

- Нагрузочное тестирование (обычное, burst, ramp-up)  
- Интеграционное тестирование (services + PostgreSQL)  
- Endpoint статистики (`GET /stats`)  
- Массовая деактивация пользователей (`POST /team/deactivate`)  
- Конфигурация линтера (golangci-lint)
