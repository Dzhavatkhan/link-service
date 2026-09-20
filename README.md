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
go build -o link-service cmd/url-shortener/main.go
```

## Запуск

```bash
./link-service
```

## Пример запроса

```bash
curl -X POST -d '{"url":"https://github.com/oscript-library/go-link-service"}' http://localhost:8080/api/v1/url
```

## Пример ответа

```json
{
  "id": 1,
  "alias": "github.com",
  "url": "https://github.com/oscript-library/go-link-service"
}
