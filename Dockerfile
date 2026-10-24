# Stage 1: Build binary
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server .

# Stage 2: Runtime image
FROM alpine:3.20

WORKDIR /app

# Install ca-certificates dan tzdata untuk koneksi TLS dan timezone Jakarta
RUN apk --no-cache add ca-certificates tzdata

# Copy binary dan openapi.json
COPY --from=builder /app/server /app/server
COPY --from=builder /app/openapi.json /app/openapi.json

EXPOSE 8080

CMD ["/app/server"]
