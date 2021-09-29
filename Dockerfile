# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.16-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./
COPY commands ./commands
COPY internal ./internal
COPY pkg ./pkg

RUN go build -o /elastic-copy

##
## Deploy
##
FROM alpine:3

WORKDIR /

COPY --from=build /elastic-copy /usr/bin/elastic-copy

RUN chmod +x /usr/bin/elastic-copy
