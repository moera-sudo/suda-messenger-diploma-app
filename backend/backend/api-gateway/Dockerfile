# Build stage
FROM golang:1.26-alpine AS builder  
# (Или 1.25, если ты уверен, что такой образ есть, но 1.25 еще не вышел в релиз)
# Если у тебя кастомный образ 1.25.4, ок. Если нет - ставь 1.22

WORKDIR /src

# 1. Копируем файлы модулей (Gateway + Proto)
# Нам нужны только те, от которых зависит этот сервис
COPY backend/proto/go.mod backend/proto/go.sum ./backend/proto/
COPY backend/api-gateway/go.mod backend/api-gateway/go.sum ./backend/api-gateway/

# 2. Создаем новый go.work ВНУТРИ образа
# Вместо COPY go.work .
RUN go work init ./backend/api-gateway ./backend/proto

# 3. Скачиваем зависимости
RUN go mod download

# 4. Копируем исходный код
COPY backend/proto ./backend/proto
COPY backend/api-gateway ./backend/api-gateway

# 5. Собираем
WORKDIR /src/backend/api-gateway
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/app ./cmd/main/main.go

# Run stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /out/app /app/app
EXPOSE 8080
ENTRYPOINT ["/app/app"]