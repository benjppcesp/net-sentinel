#----------------
# Builder
ARG GO_VERSION=1.21.5
FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /src

# Install git (needed for some module downloads) and ca-certificates if needed
RUN apk add --no-cache git ca-certificates

# 1. Copy module files and download dependencies (cache friendly)
COPY go.mod go.sum ./
# (opcional) si usas vendor, COPY vendor ./vendor y usar -mod=vendor
# RUN go mod download   <-- esta línea puede fallar si go.mod es inválido

# 2. Copy the rest of the source
COPY . .

# 3. Build a static linux binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /net-sentinel ./cmd/net-sentinel

#-----------------
# RUNTIME 
FROM alpine:3.19

# Install certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /net-sentinel /usr/local/bin/net-sentinel

ENTRYPOINT ["/usr/local/bin/net-sentinel"]
