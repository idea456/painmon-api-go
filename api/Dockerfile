FROM golang:latest as build

# Environment variables will be put here
ENV GO111MODULE=on
ENV PORT=8080

RUN mkdir /build
WORKDIR /build

COPY . .

RUN go mod tidy

EXPOSE 8080