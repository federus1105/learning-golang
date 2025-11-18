FROM golang:1.23.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go get ./...

RUN go build -o app main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/templates /app/templates

COPY --from=builder /app/app /app/app

ENTRYPOINT [ "./app" ]