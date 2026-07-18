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

WORKDIR /app

# Copy the binary
COPY --from=builder /trickreport .

EXPOSE 8080

ENTRYPOINT ["/app/trickreport"]
