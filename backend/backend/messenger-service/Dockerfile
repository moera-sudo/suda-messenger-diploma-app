
FROM golang:1.26-alpine AS builder  

WORKDIR /src

# 1. Module files
COPY backend/proto/go.mod backend/proto/go.sum ./backend/proto/
COPY backend/docs/go.mod backend/docs/go.sum ./backend/docs/
COPY backend/messenger-service/go.mod backend/messenger-service/go.sum ./backend/messenger-service/

# 2. Go workspace
RUN go work init ./backend/messenger-service ./backend/proto ./backend/docs

# 3. Dependencies
RUN go mod download

# 4. Source code
COPY backend/proto ./backend/proto
COPY backend/docs ./backend/docs
COPY backend/messenger-service ./backend/messenger-service


# 5. Собираем
WORKDIR /src/backend/messenger-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/app ./cmd/main/main.go

# Run stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /out/app /app/app
EXPOSE 8081
ENTRYPOINT ["/app/app"]