FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o temporal-server ./main.go

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/temporal-server /usr/local/bin/temporal-server
# Copy any necessary config files if we had them, but we'll use env vars and default config

# Expose ports (Frontend, History, Matching, Worker)
EXPOSE 7233 7234 7235 7239

ENTRYPOINT ["temporal-server"]
