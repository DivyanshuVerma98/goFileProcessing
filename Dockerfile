FROM golang:1.20.0

WORKDIR /usr/srv/app

RUN go install github.com/cosmtrek/air@latest

COPY . .
RUN go mod tidy