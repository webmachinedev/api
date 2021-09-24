FROM alpine:latest

RUN apk add git
RUN apk add go

ENV PORT 80

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY main.go main.go
COPY data data
CMD go run main.go
