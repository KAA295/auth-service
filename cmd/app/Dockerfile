FROM golang:1.24.2 AS build

WORKDIR /build
COPY . .

RUN go build -o cmd/app/bin/main cmd/app/main.go

FROM ubuntu:22.04

WORKDIR app

COPY --from=build /build/cmd/app/bin/main ./main
COPY --from=build /build/config/docker.yml ./config/docker.yml

CMD ["./main", "--config=./config/docker.yml"]
