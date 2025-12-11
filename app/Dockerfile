#----------------
# Builder
ARG GO_VERSION=1.21.5
FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /src

# Install git (some packages may need it), ca-certificates if needed
RUN apk add --no-cache git ca-certificates

# Copy module files and the vendor directory (cache-friendly)
COPY go.mod go.sum ./
COPY vendor ./vendor

# Copy the rest of the source
COPY . .

# Build using the vendor directory to avoid network access
ENV CGO_ENABLED=0
RUN GOOS=linux GOFLAGS="-mod=vendor" go build -ldflags="-s -w" -o /net-sentinel ./cmd/net-sentinel

#-----------------
# RUNTIME 
FROM alpine:3.19

# Install certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /net-sentinel /usr/local/bin/net-sentinel

ENTRYPOINT ["/usr/local/bin/net-sentinel"]
