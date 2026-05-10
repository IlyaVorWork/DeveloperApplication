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