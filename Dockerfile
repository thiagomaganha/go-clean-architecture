FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -mod=vendor -o ordersystem cmd/ordersystem/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/ordersystem .
COPY --from=builder /app/.env .
COPY --from=builder /app/sql/migrations ./sql/migrations

CMD ["./ordersystem"]
