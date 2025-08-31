# --- build stage ---
FROM golang:1.23.4-bullseye AS builder
WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/main.go

FROM debian:bullseye-slim AS runtime
WORKDIR /app

COPY --from=builder /app/api .

COPY .env .env
COPY config config
COPY migrations migrations
COPY static static

CMD ["./api"]
