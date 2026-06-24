FROM golang:1.26-alpine AS builder  

WORKDIR /src

COPY backend/proto/go.mod backend/proto/go.sum ./backend/proto/
COPY backend/docs/go.mod backend/docs/go.sum ./backend/docs/
COPY backend/media-service/go.mod backend/media-service/go.sum ./backend/media-service/

RUN go work init ./backend/media-service ./backend/proto ./backend/docs

RUN go mod download

COPY backend/proto ./backend/proto
COPY backend/docs ./backend/docs
COPY backend/media-service ./backend/media-service


# 5. Собираем
WORKDIR /src/backend/media-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/app ./cmd/main/main.go

# Run stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /out/app /app/app
EXPOSE 8084 9094
ENTRYPOINT ["/app/app"]