# ─── Build stage ────────────────────────────────────────────────────────────
FROM golang:1.23.4-alpine AS builder

WORKDIR /app

# Install templ CLI
RUN go install github.com/a-h/templ/cmd/templ@v0.3.1020

# Cache deps separately from source
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Generate templ files
RUN templ generate ./internal/interfaces/http/views/

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /trickreport ./cmd/api/main.go

# ─── Runtime stage ──────────────────────────────────────────────────────────
FROM alpine:3.20 AS runtime

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy the binary
COPY --from=builder /trickreport .

# Copy static assets (CSS, JS, images)
COPY --from=builder /app/internal/interfaces/http/static ./internal/interfaces/http/static

EXPOSE 8080

ENTRYPOINT ["/app/trickreport"]
