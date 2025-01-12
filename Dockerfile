# Start from the latest golang base image
FROM golang:1.23.1-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the CLI server command for main-service
RUN CGO_ENABLED=0 go build -o aqua-sec-cloud-inventory ./cmd/server/main-service

# Build the CLI server command for notification-service
RUN CGO_ENABLED=0 go build -o aqua-sec-cloud-inventory-notification ./cmd/server/notification-service

# Final image
FROM alpine:latest AS main-service
WORKDIR /app
COPY --from=builder /app/aqua-sec-cloud-inventory /app/aqua-sec-cloud-inventory
EXPOSE 8080
ENTRYPOINT ["/app/aqua-sec-cloud-inventory", "main-server"]

# Final image
FROM alpine:latest AS notification-service
WORKDIR /app
COPY --from=builder /app/aqua-sec-cloud-inventory-notification /app/aqua-sec-cloud-inventory-notification

EXPOSE 8081
ENTRYPOINT ["/app/aqua-sec-cloud-inventory-notification", "notification-server"]
