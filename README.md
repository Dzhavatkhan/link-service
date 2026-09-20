# Link Service (Учебный проект)

## Описание

Сервис для создания коротких ссылок на любой URL.

## Установка

### Предварительные требования

- PostgreSQL


### Зависимости

- PostgreSQL (pgx)
- Go 1.18+

### Сборка

```bash
go build -o link-service/cmd/url-shortener/main.go
```

## Запуск

```bash
go run -o link-service/cmd/url-shortener/main.go

```

## Пример запроса

```bash
curl -X POST -d '{"url":"https://github.com/oscript-library/go-link-service"}' http://localhost:8080/url
```

## Пример ответа

```json
{
  "alias": "github.com",
  "url": "https://github.com/oscript-library/go-link-service"
}
