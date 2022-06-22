# syntax=docker/dockerfile:1

FROM golang:1.18.2-alpine3.15

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /hStream-router

EXPOSE 8080

CMD [ "/hStream-router" ]