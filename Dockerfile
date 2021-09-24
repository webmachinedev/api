FROM alpine:latest

RUN apk add git
RUN apk add go

COPY . .

RUN go mod download
RUN go run cmd/server/main.go
