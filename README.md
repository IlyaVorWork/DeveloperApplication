# Команды для генерации

## sqlc

```
sqlc generate -f ./internal/adapters/repository/postgres/sqlc.yaml
```

## gRPC

```
buf generate --template buf.gen.application_service.yaml --path internal/adapters/grpc/application_service/application.proto
buf generate --template buf.gen.auto_verification.yaml --path internal/adapters/grpc/auto_verification/auto_verification.proto
```

## Swagger

```
swag init -g ./cmd/app/main.go
```

## ApplicationUpload

### Обязательные при загрузке

- `code_name` — кодовое имя приложения
- `category_id` — ID категории (из application_service), по умолчанию 0 (Uncategorized)
- `android_package_name` — Android package name (напр. `com.TPU.DIVE_E9`); при публикации передаётся как `launch_url` в CreateRepository
- `default_locale` — локаль по умолчанию (напр. `ru`)
- `name` — название приложения (→ `ApplicationTranslation.name`)
- `short_title` — краткий заголовок (→ `ApplicationTranslation.short_title`)
- `version` — версия в формате `major.minor.micro` (напр. `1.0.0`)

### Опциональные при загрузке, обязательные при публикации

- `inapp_background_image` — фоновое изображение для in-app / launcher (1000×664, 1200×806)

### Опциональные

**Медиа (→ `Application`)**
- `web_video` — ссылка на видео для web
- `inapp_video` — ссылка на видео для in-app
- `web_background_image` — фоновое изображение для web
- `video_cover` — обложка видео

**Перевод (→ `ApplicationTranslation`)**
- `description`
- `goals` — цели
- `tasks` — задачи
- `results` — результаты
- `challenges`
- `location`
- `safety` — техника безопасности