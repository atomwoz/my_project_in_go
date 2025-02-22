# Build stage with CGO enabled
FROM golang:1.24.0-alpine AS builder

# Install build dependencies for CGO
RUN apk add --no-cache build-base

RUN mkdir /app
WORKDIR /app

# Cache Go module downloads
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# RUN TESTS
# RUN go test ./...


# Set CGO enabled globally
ENV CGO_ENABLED=1

# Build executables with CGO support
RUN mkdir -p /app/build
RUN GOOS=linux go build -o /app/build/main ./cmd/main
RUN GOOS=linux go build -o /app/build/data_importer ./cmd/data_importer

# ================================================
# Final stage (Runner image)
FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN mkdir /app
WORKDIR /app

# Copy binaries and additional data
COPY --from=builder /app/build/ .
COPY --from=builder /app/data/* .

EXPOSE 8080
CMD ["sh", "-c", "./data_importer && ./main"]
