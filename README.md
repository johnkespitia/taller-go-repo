**Project**: Taller Go Repo

- **Description**: HTTP application in Go that exposes a simple API for managing events and uses PostgreSQL as storage. It comes with Docker files for development and production.

**Requirements**:
- **Docker / Docker Compose**: required to start the database and app in containers.
- **Go 1.21**: if you want to compile or run the app locally.

**Important files**:
- **`docker-compose.yml`**: definition of `app` and `db` services.
- **`docker/go.dockerfile`**: Dockerfile for the application (has `builder` and `dev` targets).
- **`docker/posgresql.dockerfile`**: Dockerfile for the DB image (uses `postgres`).
- **`docker-entrypoint-initdb.d/init.sql`**: database initialization script (creates `events` table).
- **`app/`**: Go code of the application (entrypoint `app/cmd/main.go`).

**How to run (with Docker Compose)**
- **Build and start (production)**:
```bash
cd /Users/macbookproa2141/taller-go-repo
docker compose build --progress=plain
docker compose up -d
docker compose logs -f app
```

- **Stop and remove containers**:
```bash
docker compose down
```

- **Recreate the database from `init.sql` (loses data).** If you want to apply the new `init.sql`, you must delete the data volume so that the initialization script runs again:
```bash
docker compose down
# delete volume (WARNING: deletes data)
docker volume rm taller-go-repo_postgres_data
docker compose up -d db
```

**Development mode (live-reload with `air`)**
- The `Dockerfile` contains a `target: dev` that installs development tools (`air` or `CompileDaemon`). To use it with `docker compose`, add the `target: dev` and mount the code volume:

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

Then build and start:
```bash
docker compose build --progress=plain --build-arg NODE_ENV=development
docker compose up app
```

Or alternatively build the `dev` image and run it manually:
```bash
docker build -f docker/go.dockerfile --target dev -t taller-go-dev .
docker run --rm -it -v "$(pwd)/app:/app/app" -p 8080:8080 taller-go-dev
```

**Run locally (without Docker)**
- Install dependencies and compile:
```bash
cd app
go mod tidy
go build ./cmd
./cmd/main
```

**Relevant environment variables**
- `DATABASE_URL`: full connection string (e.g. `postgres://postgres:postgres@db:5432/myapp?sslmode=disable`).
- `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, `POSTGRES_HOST`, `POSTGRES_PORT` — used if `DATABASE_URL` does not exist.

**API - Endpoints and examples**
- **Health check**
- `GET /api/health`
  - Response: `200 OK` with `{"status":"healthy"}`

- **List events**
- `GET /api/events`
  - Example:
    ```bash
    curl http://localhost:8080/api/events
    ```
  - Response: `200 OK` with JSON array of events.

- **Create event**
- `POST /api/events`
  - JSON body (minimal example):
    ```json
    {
      "title": "Team meeting",
      "description": "Weekly review",
      "start_time": "2025-12-20T10:00:00Z",
      "end_time": "2025-12-20T11:00:00Z"
    }
    ```
  - Curl example:
    ```bash
    curl -X POST http://localhost:8080/api/events \
      -H "Content-Type: application/json" \
      -d '{"title":"UUID test","description":"desc","start_time":"2025-12-20T10:00:00Z","end_time":"2025-12-20T11:00:00Z"}'
    ```
  - Expected response: `201 Created` with the created event (includes `id` UUID and `created_at`).

- **Update event**
- `PUT /api/events` (or `PUT /api/events?id=<uuid>`)
  - JSON body (with `id` or `id` in query param):
    ```json
    {
      "id": "a3f27f2b-...",
      "title": "Updated meeting",
      "description": "Time changed",
      "start_time": "2025-12-20T11:00:00Z",
      "end_time": "2025-12-20T12:00:00Z"
    }
    ```
  - Curl example (id in body):
    ```bash
    curl -X PUT http://localhost:8080/api/events \
      -H "Content-Type: application/json" \
      -d '{"id":"<UUID>","title":"Updated","description":"...","start_time":"2025-12-20T11:00:00Z","end_time":"2025-12-20T12:00:00Z"}'
    ```

**Date format**
- Use timestamps in ISO 8601 / RFC3339 format (e.g.: `2025-12-20T10:00:00Z`). The server maps these values to `time.Time`.

**Notes on schema and UUID**
- The database uses UUID for `id` and the `init.sql` script enables `pgcrypto` and defines `id UUID PRIMARY KEY DEFAULT gen_random_uuid()`. This allows the DB to generate the UUID automatically on insert.
- If you already have data in the Postgres volume and want to apply the new table/column, perform a migration or recreate the volume (see section above).

**Concurrency tests**
- To test multiple concurrent inserts (verify mutex and DB):
```bash
for i in {1..10}; do
  curl -s -X POST http://localhost:8080/api/events \
    -H "Content-Type: application/json" \
    -d '{"title":"Concurrency test","description":"desc","start_time":"2025-12-20T10:00:00Z","end_time":"2025-12-20T11:00:00Z"}' &
done
wait
```
