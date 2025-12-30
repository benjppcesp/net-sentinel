FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Ajuste en la ruta del build para apuntar a la nueva carpeta cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/net-sentinel/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
# Exponemos el puerto de métricas que definiste en el código
EXPOSE 2112
CMD ["./main"]