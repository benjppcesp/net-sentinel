#----------------
# Builder

FROM golang:1.21-alpine AS builder

WORKDIR /app

# Inatalamos git
RUN apk add --no-cache git

# 1. Copiar archivos del modulo (go.mod y go.sum)
COPY go.mod .
COPY go.sum .

# 2. Copiar main.go
COPY main.go .

# 3. Copilar
RUN CG0_ENABLED=0 go build -o net-sentinel main.go

#-----------------
# RUNTIME 
FROM alpine:latest

# Instalar certificadosy utilidades necesarias para la red
RUN apk --no-cache add ca-certificates
EXPOSE 2112

# Copiamos el binario compilado de la Etapa 1
COPY --from=builder /app/net-sentinel /usr/local/bin/net-sentinel

ENTRYPOINT ["/usr/local/bin/net-sentinel"]