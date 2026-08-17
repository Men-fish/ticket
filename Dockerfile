FROM golang:1.23-alpine

WORKDIR /src/app

RUN go install github.com/air-verse/air@latest

COPY . .
COPY .air.toml .air.toml
RUN go mod tidy
