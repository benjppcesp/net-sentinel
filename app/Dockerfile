#----------------
# Builder
ARG GO_VERSION=1.21.5
FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /src

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates

# Copy module files
COPY go.mod go.sum ./

# Download dependencies (sin verify para evitar problemas)
RUN go mod download

# Copy source code
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /net-sentinel ./cmd/net-sentinel

#-----------------
# RUNTIME 
FROM alpine:3.19

# Install certificates
RUN apk --no-cache add ca-certificates

# Copy binary
COPY --from=builder /net-sentinel /usr/local/bin/net-sentinel

ENTRYPOINT ["/usr/local/bin/net-sentinel"]
