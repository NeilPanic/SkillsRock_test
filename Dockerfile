FROM golang:1.22-alpine AS builder

RUN apk add --no-cache upx

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o todo-api ./cmd/api



FROM alpine:3.20

COPY --from=builder /app/todo-api /usr/bin/todo-api

ENV PORT=8080 \
    DATABASE_DSN="postgres://auth_user:auth_password@db:5433/auth_db?sslmode=disable"


EXPOSE ${PORT}

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s \
  CMD wget -qO- http://localhost:${PORT}/tasks || exit 1

ENTRYPOINT ["todo-api"]
