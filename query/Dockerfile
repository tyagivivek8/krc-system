FROM golang:1.25.4 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o query ./cmd/query

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/query .
CMD ["./query"]
