# backend-practice

Servicio HTTP de práctica.

Requisitos

- Go 1.20+ (se probó con 1.22.3)

Ejecutar

Desde la raíz del módulo (donde está `go.mod`):

```bash
# Ejecuta usando el Makefile (usa CGO deshabilitado por defecto)
make run
```

Alternativa sin Makefile:

```bash
# Ejecutar directamente (deshabilita cgo si tienes problemas con dyld en macOS)
CGO_ENABLED=0 go run ./cmd
```

Notas

- Si dependes de C libraries o necesitas cgo, asegúrate de tener las Xcode Command Line Tools instaladas:

```bash
xcode-select --install
```

Si tienes problemas con el enlazador dinámico (`dyld`) al usar `go run`, intenta compilar con `go build` y ejecutar el binario resultante para obtener más información.

## Docker

Puedes construir y ejecutar la imagen Docker usando los targets del `Makefile`:

````bash
# backend-practice

Pequeña API en Go diseñada como proyecto de práctica y demostración técnica.

Descripción
-----------
Servicio HTTP mínimo que expone un endpoint `/health`. El objetivo es mostrar una estructura de proyecto Go ordenada, buenas prácticas de build y despliegue con Docker.

Requisitos
---------
- Go 1.22.3 (fijado en `go.mod`)
- Docker (opcional, para construir y ejecutar imágenes)

Estructura clave
----------------
- `cmd/` — puntos de entrada (main)
- `internal/infra/transport` — router y handlers
- `Makefile` — tareas comunes (run, build, docker-* )
- `dockerfile` — Dockerfile multi-stage para builds reproducibles

Ejecutar localmente
-------------------
1. Desde la raíz del proyecto (donde está `go.mod`):

```bash
# Ejecuta la aplicación (Makefile usa CGO deshabilitado por defecto)
make run
````

2. Comprobá el endpoint:

```bash
curl http://localhost:8080/health
# -> {"status":"ok"}
```

Alternativa sin Makefile:

```bash
CGO_ENABLED=0 go run ./cmd
```

## Notas sobre CGO y macOS

En macOS (especialmente en M1/M2) la compilación con `cgo` activado puede causar errores del enlazador (dyld). Para reproducibilidad y evitar esos errores en desarrollo usamos `CGO_ENABLED=0`. Si necesitás cgo por dependencias nativas, compilá dentro de un contenedor Linux o asegurate de tener las Xcode Command Line Tools instaladas (`xcode-select --install`).

## Docker

El proyecto incluye un `dockerfile` multi-stage que compila el binario (estático, sin cgo) y produce una imagen final mínima (distroless nonroot).

## Comandos útiles (Makefile)

```bash
# Construir imagen (etiqueta: auth-service)
make docker-build

# Ejecutar la imagen en background (mapea puerto 8080)
make docker-run

# Parar y eliminar el contenedor de ejemplo
make docker-stop
```

Durante la sesión se creó además una etiqueta de imagen `auth-service:secure` (imagen optimizada/minimal). Para auditoría, escaneá la imagen con `trivy` o `grype`:

```bash
# ejemplo con trivy
trivy image auth-service:secure
```
