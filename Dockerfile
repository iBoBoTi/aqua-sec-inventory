# Start from the latest golang base image
FROM golang:1.23.1-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the CLI server command
RUN CGO_ENABLED=0 go build -o aqua-sec-cloud-inventory ./cmd/server

# Final image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/aqua-sec-cloud-inventory /app/aqua-sec-cloud-inventory

EXPOSE 8080
ENTRYPOINT ["/app/aqua-sec-cloud-inventory", "server"]
