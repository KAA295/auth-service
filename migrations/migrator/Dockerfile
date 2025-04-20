FROM golang:1.24.2

WORKDIR /app

COPY /migrations ./migrations

RUN go install github.com/rubenv/sql-migrate/...@latest