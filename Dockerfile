# Stage 1: Build binary
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git dan ca-certificates
RUN apk add --no-cache git ca-certificates

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server .

# Stage 2: Runtime image
FROM alpine:3.21

WORKDIR /app

# Install ca-certificates dan tzdata untuk koneksi TLS dan timezone Jakarta
RUN apk --no-cache add ca-certificates tzdata

# Copy binary dan openapi.json
COPY --from=builder /app/server /app/server
COPY --from=builder /app/openapi.json /app/openapi.json

EXPOSE 8080

CMD ["/app/server"]
