# ─── Build stage ────────────────────────────────────────────────────────────
FROM golang:1.23.4-alpine AS builder

WORKDIR /app

# Cache deps separately from source
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /trickreport ./cmd/api/main.go

# ─── Runtime stage ──────────────────────────────────────────────────────────
FROM alpine:3.20 AS runtime

RUN apk add --no-cache ca-certificates tzdata

# Non-root user for the runtime container
RUN addgroup -S appgroup && adduser -S -G appgroup appuser

WORKDIR /app

# Copy the binary
COPY --from=builder /trickreport .

# Run as the non-root user
USER appuser

EXPOSE 8080

# Health check using the /health endpoint (busybox wget is present in alpine)
HEALTHCHECK --interval=10s --timeout=5s --retries=5 --start-period=15s \
  CMD wget -qO- http://localhost:8080/health || exit 1

ENTRYPOINT ["/app/trickreport"]
