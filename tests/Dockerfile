FROM golang:1.24.2 AS build

WORKDIR /build

COPY . .

CMD ["go", "test", "./tests/auth_test.go"]
