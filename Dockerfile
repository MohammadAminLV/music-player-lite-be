
# Build stage
FROM golang:1.25-alpine AS builder
RUN apk add --no-cache git build-base

WORKDIR /src
# Copy modules first to leverage cache
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Build a static-ish binary; strip symbol table to reduce size.
ENV CGO_ENABLED=0
RUN go build -ldflags "-s -w" -o /out/music-player-lite-be ./ 

# Final stage
FROM alpine:3.18

# create unprivileged user
RUN addgroup -S app && adduser -S -G app app

# ca-certificates in case your app needs outgoing HTTPS (e.g., fetching remote resources)
RUN apk add --no-cache ca-certificates && update-ca-certificates

WORKDIR /app
# Copy binary from builder
COPY --from=builder /out/music-player-lite-be /app/music-player-lite-be

# Copy default data.json into image so image is self-contained (override at run time with a volume if desired)
COPY data.json.example /app/data.json

# Ensure correct ownership
RUN chown -R app:app /app

# Use non-root user
USER app

EXPOSE 8000

# Healthcheck for container orchestrators (adjust timeout/interval as needed)
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget -qO- --timeout=3 http://127.0.0.1:8000/api/tracks >/dev/null || exit 1

# Entrypoint - respects PORT and DATA_FILE environment variables set at runtime
ENV PORT=8000
ENV DATA_FILE=/app/data.json

CMD ["/app/music-player-lite-be"]