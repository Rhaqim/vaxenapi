# ============================================================
# Vaxen API — Multi-stage Dockerfile
# ============================================================
# Produces a minimal image with three binaries:
#   /app/server   — the API server
#   /app/migrate  — database migration tool
#   /app/seed     — admin/exchange-rate seed tool
#
# Usage:
#   docker build -t vaxen-api .
#   docker run -p 8080:8080 --env-file .env vaxen-api
#   docker run --env-file .env vaxen-api /app/seed admin admin@vaxen.io password
#   docker run --env-file .env vaxen-api /app/migrate
# ============================================================

# --- Build stage ---
FROM golang:1.25-alpine AS builder

WORKDIR /src

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build all three binaries — static, no CGO, works on any Linux
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server  ./main.go
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/migrate ./cmd/migrate/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/seed    ./cmd/seed/main.go

# --- Runtime stage ---
FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

# Run as non-root
RUN adduser -D -u 1000 vaxen
USER vaxen

WORKDIR /app

COPY --from=builder /out/server  /app/server
COPY --from=builder /out/migrate /app/migrate
COPY --from=builder /out/seed    /app/seed

EXPOSE 8080

ENTRYPOINT ["/app/server"]
