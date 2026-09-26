# Планировщик задач — дипломный проект Go

Веб-сервер для TODO-планировщика из финального проекта Яндекс Практикума. Сервер хранит задачи в SQLite, отдаёт готовый фронтенд из каталога `web` и реализует REST API для создания, просмотра, редактирования, выполнения и удаления задач.

## Реализовано

- SQLite база `scheduler.db`, таблица `scheduler`, индекс по `date`.
- `GET /api/nextdate`.
- Повторения `d`, `y`, а также задания со звёздочкой `w` и `m`.
- `POST /api/task` — добавление задачи.
- `GET /api/tasks` — список ближайших задач.
- `GET /api/task?id=...` — получение задачи.
- `PUT /api/task` — изменение задачи.
- `DELETE /api/task?id=...` — удаление задачи.
- `POST /api/task/done?id=...` — выполнение задачи.
- Поиск по `title` и `comment` через `GET /api/tasks?search=...`.
- Переменные окружения `TODO_PORT`, `TODO_DBFILE`, `TODO_PASSWORD`.
- Аутентификация через `POST /api/signin` и JWT-подобный HS256 токен в cookie `token`.
- Dockerfile.

## Подготовка шаблона

В корне проекта должны оставаться исходные каталоги `web` и `tests` из официального шаблона `Yandex-Practicum/go_final_project`.

После копирования файлов выполните:

```bash
go mod tidy
```

## Локальный запуск

```bash
go run .
```

По умолчанию приложение доступно по адресу:

```text
http://localhost:7540/
```

Опциональные переменные окружения:

```text
TODO_PORT=7540
TODO_DBFILE=scheduler.db
TODO_PASSWORD=
```

Если `TODO_PASSWORD` пустой, авторизация отключена и базовые тесты работают без токена.

## Тесты

Сначала запустите сервер:

```bash
go run .
```

В другом терминале:

```bash
go test ./tests
```

Для проверки отдельных шагов:

```bash
go test -run ^TestDB$ ./tests
go test -run ^TestNextDate$ ./tests
go test -run ^TestAddTask$ ./tests
go test -run ^TestTasks$ ./tests
go test -run ^TestTask$ ./tests
go test -run ^TestEditTask$ ./tests
go test -run ^TestDone$ ./tests
go test -run ^TestDelTask$ ./tests
```

В `tests/settings.go` для полной проверки правил повторения можно установить:

```go
var FullNextDate = true
```

Если тестируете авторизацию, установите `TODO_PASSWORD`, выполните `/api/signin`, затем укажите полученный токен в `tests/settings.go`.

## Docker

Сборка:

```bash
docker build -t go-scheduler .
```

Linux/macOS:

```bash
docker run --rm -p 7540:7540 \
  -e TODO_DBFILE=/data/scheduler.db \
  -v "$(pwd):/data" \
  go-scheduler
```

PowerShell:

```powershell
docker run --rm -p 7540:7540 `
  -e TODO_DBFILE=/data/scheduler.db `
  -v "${PWD}:/data" `
  go-scheduler
```
