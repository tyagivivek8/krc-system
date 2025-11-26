# orchestrator/Dockerfile
FROM golang:1.25.4 AS builder

WORKDIR /src

# Copy the whole repo so replace ../audit and ../query resolve
COPY . .

WORKDIR /src/orchestrator

# Download modules
RUN go mod download

# Build a static linux/amd64 binary and place into /out
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /out/orchestrator ./cmd/orchestrator

# Ensure the binary is executable inside the builder image
RUN chmod +x /out/orchestrator

FROM alpine:3.20

# install ca-certificates for outbound HTTPS and a shell for debugging
RUN apk add --no-cache ca-certificates bash

WORKDIR /app

# Copy the built binary from builder and ensure exec bit remains
COPY --from=builder /out/orchestrator ./orchestrator
RUN chmod +x ./orchestrator

# optional: create unprivileged user
RUN addgroup -S app && adduser -S app -G app
USER app

ENV PATH=/app:$PATH

CMD ["./orchestrator"]
