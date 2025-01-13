FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o rate-limit ./cmd/server/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/rate-limit .
COPY --from=builder /app/.env .env

EXPOSE 8080

CMD ["./rate-limit"]