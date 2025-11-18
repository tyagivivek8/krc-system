# Build stage
FROM golang:1.25.4 AS builder

WORKDIR /app

# Copy go.mod and go.sum first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build a static Linux binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o audit ./cmd/audit

# Runtime stage
FROM alpine:3.20

WORKDIR /app

# Copy the binary from the builder
COPY --from=builder /app/audit .

# Just to be safe, ensure it is executable
RUN chmod +x ./audit

CMD ["./audit"]
