# ---- Etapa 1: Compilación (Build) ----
# Usamos la imagen oficial de Go (misma versión que tu go.mod)
FROM golang:1.24-alpine AS builder

WORKDIR /src

# Copiamos los archivos de dependencias y las descargamos
COPY go.mod go.sum ./
RUN go mod download

# Copiamos todo el código fuente
COPY . .

# Compilamos la aplicación para Linux
# Esto crea un binario llamado 'main' en la carpeta /app
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /app/main ./cmd/web

# ---- Etapa 2: Producción (Final) ----
# Usamos una imagen mínima y ligera
FROM alpine:latest

WORKDIR /app

# Copiamos SÓLO el binario compilado desde la etapa 'builder'
COPY --from=builder /app/main .

# ¡Importante! Copiamos tus plantillas y archivos estáticos
COPY static ./static
COPY templates ./templates

# Exponemos el puerto 8080 (el que usa tu app)
EXPOSE 8080

# El comando para iniciar tu aplicación
CMD ["./main"]