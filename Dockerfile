FROM docker.io/library/golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ordersystem cmd/ordersystem/main.go

FROM docker.io/library/alpine:latest

WORKDIR /app

COPY --from=builder /app/ordersystem .
COPY --from=builder /app/.env .
COPY --from=builder /app/sql/migrations ./sql/migrations

CMD ["./ordersystem"]
