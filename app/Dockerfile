#----------------
# Builder
ARG GO_VERSION=1.21.5
FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /src

# Install git (needed for some module downloads) and ca-certificates if needed
RUN apk add --no-cache git ca-certificates

# Copy module files (keeps caching possible) and skip explicit go mod download
COPY go.mod go.sum ./

# Copy the rest of the source
COPY . .

# Build a static linux binary (go build will fetch modules if needed)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /net-sentinel ./cmd/net-sentinel

#-----------------
# RUNTIME 
FROM alpine:3.19

# Install certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage (path matches the -o above)
COPY --from=builder /net-sentinel /usr/local/bin/net-sentinel

ENTRYPOINT ["/usr/local/bin/net-sentinel"]
