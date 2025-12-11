#----------------
# Builder
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install git (needed for some module downloads) and ca-certificates if needed
RUN apk add --no-cache git

# 1. Copy module files and download dependencies (cache friendly)
COPY go.mod go.sum ./
RUN go mod download

# 2. Copy the rest of the source
COPY . .

# 3. Build a static linux binary
ENV CGO_ENABLED=0
RUN GOOS=linux GOARCH=amd64 go build -o /app/net-sentinel ./main.go

#-----------------
# RUNTIME 
FROM alpine:3.18

# Install certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

EXPOSE 2112

# Copy the binary from the builder stage
COPY --from=builder /app/net-sentinel /usr/local/bin/net-sentinel

# Create non-root user and switch
RUN addgroup -S app && adduser -S -G app app
USER app

ENTRYPOINT ["/usr/local/bin/net-sentinel"]