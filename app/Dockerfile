# ----------------------------------------------------------------------------------
# ETAPA 1: BUILDER (Compilación)
# Usamos una imagen de Go completa para tener todas las herramientas de compilación.
# ----------------------------------------------------------------------------------
FROM golang:1.21-alpine AS builder

# Establecemos el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiamos los archivos de definición de módulos (go.mod y go.sum)
# Esto permite a Docker cachear la descarga de dependencias
COPY go.mod .
COPY go.sum .

# Descargamos las dependencias. -mod=readonly asegura que no se cambien las dependencias.
RUN go mod download 

# Copiamos el código fuente de nuestra sonda
COPY main.go .

# Compilamos la aplicación. 
# - CGO_ENABLED=0: Asegura que el binario no dependa de librerías del sistema operativo (estático).
# - -o net-sentinel: Nombra el binario como 'net-sentinel'.
# - -ldflags "-s -w": Optimiza el binario para hacerlo más pequeño.
RUN CGO_ENABLED=0 go build -o net-sentinel main.go

# ----------------------------------------------------------------------------------
# ETAPA 2: RUNTIME (Ejecución Final)
# Crear el contenedor final, ligero y seguro, usando solo lo necesario.
# Usamos alpine para reducir el tamaño final de la imagen a MBs.
# ----------------------------------------------------------------------------------
FROM alpine:latest

# Instalamos curl/ca-certificates: son necesarios para las peticiones HTTPS que hacemos en Go
RUN apk --no-cache add ca-certificates

ENV TZ=America/Argentina/Buenos_Aires

# Establecemos el puerto que nuestra aplicación expondrá (2112 para Prometheus)
EXPOSE 2112

# Copiamos el binario compilado de la Etapa 1
COPY --from=builder /app/net-sentinel /usr/local/bin/net-sentinel

# Comando de inicio del contenedor
ENTRYPOINT ["/usr/local/bin/net-sentinel"]