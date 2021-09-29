FROM alpine:latest

RUN apk add go

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY main.go main.go
COPY data data
CMD go run main.go
