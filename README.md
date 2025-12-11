**Proyecto**: Taller Go Repo

- **Descripción**: Aplicación HTTP en Go que expone una API simple para gestionar eventos y usa PostgreSQL como almacenamiento. Se suministra con archivos Docker para desarrollo y producción.

**Requisitos**:
- **Docker / Docker Compose**: necesarios para levantar la base de datos y la app en contenedores.
- **Go 1.21**: si quieres compilar o ejecutar la app localmente.

**Archivos importantes**:
- **`docker-compose.yml`**: definición de servicios `app` y `db`.
- **`docker/go.dockerfile`**: Dockerfile para la aplicación (tiene target `builder` y `dev`).
- **`docker/posgresql.dockerfile`**: Dockerfile para la imagen de DB (usa `postgres`).
- **`docker-entrypoint-initdb.d/init.sql`**: script de inicialización de la base de datos (crea tabla `events`).
- **`app/`**: código go de la aplicación (entrypoint `app/cmd/main.go`).

**Cómo ejecutar (con Docker Compose)**
- **Construir y levantar (producción)**:
```bash
cd /Users/macbookproa2141/taller-go-repo
docker compose build --progress=plain
docker compose up -d
docker compose logs -f app
```

- **Parar y eliminar contenedores**:
```bash
docker compose down
```

- **Recrear la base de datos desde `init.sql` (pierde datos).** Si quieres aplicar el `init.sql` nuevo debes eliminar el volumen de datos para que se vuelva a ejecutar el script de inicialización:
```bash
docker compose down
# eliminar volumen (ATENCIÓN: borra datos)
docker volume rm taller-go-repo_postgres_data
docker compose up -d db
```

**Modo desarrollo (live-reload con `air`)**
- El `Dockerfile` contiene un `target: dev` que instala herramientas de desarrollo (`air` o `CompileDaemon`). Para usarlo con `docker compose` añade el `target: dev` y monta el volumen del código:

```yaml
services:
  app:
    build:
      context: .
      dockerfile: ./docker/go.dockerfile
      target: dev
    volumes:
      - ./app:/app/app
    ports:
      - "8080:8080"
```

Luego construir y levantar:
```bash
docker compose build --progress=plain --build-arg NODE_ENV=development
docker compose up app
```

O alternativamente construir la imagen `dev` y ejecutarla manualmente:
```bash
docker build -f docker/go.dockerfile --target dev -t taller-go-dev .
docker run --rm -it -v "$(pwd)/app:/app/app" -p 8080:8080 taller-go-dev
```

**Ejecutar localmente (sin Docker)**
- Instala dependencias y compila:
```bash
cd app
go mod tidy
go build ./cmd
./cmd/main
```

**Variables de entorno relevantes**
- `DATABASE_URL`: cadena de conexión completa (ej. `postgres://postgres:postgres@db:5432/myapp?sslmode=disable`).
- `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, `POSTGRES_HOST`, `POSTGRES_PORT` — usadas si no existe `DATABASE_URL`.

**API - Endpoints y ejemplos**
- **Health check**
- `GET /api/health`
  - Respuesta: `200 OK` con `{"status":"healthy"}`

- **Listar eventos**
- `GET /api/events`
  - Ejemplo:
    ```bash
    curl http://localhost:8080/api/events
    ```
  - Respuesta: `200 OK` con JSON array de eventos.

- **Crear evento**
- `POST /api/events`
  - Body JSON (ejemplo mínimo):
    ```json
    {
      "title": "Reunión de equipo",
      "description": "Repaso semanal",
      "start_time": "2025-12-20T10:00:00Z",
      "end_time": "2025-12-20T11:00:00Z"
    }
    ```
  - Curl ejemplo:
    ```bash
    curl -X POST http://localhost:8080/api/events \
      -H "Content-Type: application/json" \
      -d '{"title":"Prueba UUID","description":"desc","start_time":"2025-12-20T10:00:00Z","end_time":"2025-12-20T11:00:00Z"}'
    ```
  - Respuesta esperada: `201 Created` con el evento creado (incluye `id` UUID y `created_at`).

- **Actualizar evento**
- `PUT /api/events` (o `PUT /api/events?id=<uuid>`)
  - Body JSON (con `id` o `id` en query param):
    ```json
    {
      "id": "a3f27f2b-...",
      "title": "Reunión actualizada",
      "description": "Hora cambiada",
      "start_time": "2025-12-20T11:00:00Z",
      "end_time": "2025-12-20T12:00:00Z"
    }
    ```
  - Curl ejemplo (id en body):
    ```bash
    curl -X PUT http://localhost:8080/api/events \
      -H "Content-Type: application/json" \
      -d '{"id":"<UUID>","title":"Actualizado","description":"...","start_time":"2025-12-20T11:00:00Z","end_time":"2025-12-20T12:00:00Z"}'
    ```

- **Eliminar evento**
- `DELETE /api/events?id=<uuid>`
  - Ejemplo:
    ```bash
    curl -X DELETE "http://localhost:8080/api/events?id=<UUID>"
    ```

**Formato de fecha**
- Use timestamps en formato ISO 8601 / RFC3339 (ej.: `2025-12-20T10:00:00Z`). El servidor mapea estos valores a `time.Time`.

**Notas sobre el esquema y UUID**
- La base usa UUID para `id` y el script `init.sql` habilita `pgcrypto` y define `id UUID PRIMARY KEY DEFAULT gen_random_uuid()`. Esto permite que la BD genere el UUID automáticamente al insertar.
- Si ya tienes datos en el volumen de Postgres y quieres aplicar la nueva tabla/columna, realiza una migración o recrea el volumen (ver sección arriba).

**Pruebas de concurrencia**
- Para probar múltiples inserts concurrentes (verificar mutex y DB):
```bash
for i in {1..10}; do
  curl -s -X POST http://localhost:8080/api/events \
    -H "Content-Type: application/json" \
    -d '{"title":"Prueba concurrencia","description":"desc","start_time":"2025-12-20T10:00:00Z","end_time":"2025-12-20T11:00:00Z"}' &
done
wait
```

**Problemas comunes**
- `pq: SSL is not enabled on the server` — asegúrate de que la conexión use `sslmode=disable` en `DATABASE_URL` o que el código lo añada automáticamente.
- `invalid input syntax for type uuid` — ocurre cuando insertas un valor que no es UUID en una columna UUID. Si usas UUID generados por la BD, no envíes un `id` numérico.

