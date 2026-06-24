
FROM golang:1.26-alpine AS builder
 
WORKDIR /src
 
COPY backend/proto/go.mod backend/proto/go.sum ./backend/proto/
COPY backend/docs/go.mod backend/docs/go.sum ./backend/docs/
COPY backend/transaction-service/go.mod backend/transaction-service/go.sum ./backend/transaction-service/

RUN go work init ./backend/transaction-service ./backend/proto ./backend/docs

RUN go mod download

COPY backend/proto ./backend/proto
COPY backend/docs ./backend/docs
COPY backend/transaction-service ./backend/transaction-service

 
# 5. Build
WORKDIR /src/backend/transaction-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/app ./cmd/main/main.go
 
# Run stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /out/app /app/app
EXPOSE 8082
ENTRYPOINT ["/app/app"]
 