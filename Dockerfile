# syntax=docker/dockerfile:1
FROM golang:1.18.3-buster

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /hStream-router

RUN go build  -o /go-binary

CMD [ "/go-binary" ]