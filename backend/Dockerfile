FROM golang:1.26.10-alpine3.24@sha256:5cb67f19c5e7994dcb80c074f7cfca8aeeada45bbdca035f73df269cefdcce61 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /trickreport ./cmd/api

FROM alpine:3.24.3@sha256:b3413572228e34bd1c9245f5565e6cd127005655e8a7899aef33d8e83255c9e6 AS runtime
RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S appgroup && adduser -S -G appgroup appuser
WORKDIR /app
COPY --from=builder /trickreport .
USER appuser
EXPOSE 8080
HEALTHCHECK --interval=10s --timeout=5s --retries=5 --start-period=15s \
  CMD wget -qO- http://127.0.0.1:8080/ready || exit 1
ENTRYPOINT ["/app/trickreport"]
