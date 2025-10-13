# Backend Practice - Go + MySQL + Docker

Este proyecto es una API backend desarrollada en Go, usando MySQL como base de datos y corriendo en contenedores Docker. Está diseñada para ser fácil de levantar y probar, ideal para mostrar tu stack a recruiters o colegas.

## Requisitos

- [Docker](https://www.docker.com/) instalado
- [Docker Compose](https://docs.docker.com/compose/install/) (generalmente viene con Docker Desktop)

> No es necesario tener Go ni MySQL instalados localmente, todo corre dentro de contenedores.

## Levantar la app

Desde la raíz del proyecto:

```bash
docker compose up -d
```

Esto hará lo siguiente:

```bash
Levantar un contenedor de MySQL (backend_mysql) con la base de datos backend.
Levantar un contenedor de la API (backend_app) corriendo en http://localhost:8080.
Exponer el puerto 8080 para interactuar con la API desde Postman, curl o tu navegador.
```

## Endpoints disponibles:

```
GET /health → Verifica que la API esté levantada.

POST /users → Crear un nuevo usuario en la base de datos.
```

Puedes usar Postman o curl para probar los endpoints.

### Ejemplo con curl:

```bash
curl --location 'http://localhost:8080/health'
```

### Variables de entorno

Las variables de configuración de la base de datos se encuentran en .env.example. La app automáticamente las carga desde ahí, por lo que no es necesario hacer export manualmente.
Variables principales:

```bash
MYSQL_ROOT_PASSWORD=rootpw
MYSQL_USER=backend
MYSQL_PASSWORD=backendpw
MYSQL_DATABASE=backend
MYSQL_HOST_PORT=3306
APP_PORT=8080
```

### Detalles técnicos

Go 1.24
MySQL 8.0
Gin como framework HTTP
Arquitectura hexagonal: separando core, usecases y repositorios
Docker para contenedores y portabilidad
docker-compose para orquestar app y DB
Cómo parar los contenedores
docker compose down
Esto detiene y elimina los contenedores sin borrar los datos gracias al volumen persistente db_data.
