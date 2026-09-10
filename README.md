# Task Manager API

An API for managing tasks, built with Go, PostgreSQL, and Gin. The project uses
**Hexagonal Architecture (Ports and Adapters)** and **Domain-Driven Design**
principles.

## Features

- Create, retrieve, update, and delete tasks
- Cursor-based pagination for listing tasks
- Filter the task list by `status`
- Sort by `id` in descending order
- Versioned migrations with `golang-migrate`
- Redis cache-aside for task reads with update/delete invalidation
- Prometheus metrics for request totals, latency, and current task count
- Run PostgreSQL, migrations, and the API with Docker Compose

## Prerequisites

For local execution:

- Go 1.27+
- PostgreSQL

For Docker execution:

- Docker
- Docker Compose plugin

## Local Execution

Start PostgreSQL and configure the environment variables:

```bash
export PORT=8080
export REDIS_ADDR='localhost:6379'
export PG_DSN='postgres://postgres:postgres@localhost:5432/taskmanager?sslmode=disable'
```

Then run the migrations with the `migrate` tool:

```bash
migrate \
  -path setup/migrations \
  -database "$PG_DSN" \
  up
```

Run the API:

```bash
/usr/local/go/bin/go run ./cmd/api
```

Or, if Go is available in your `PATH`:

```bash
go run ./cmd/api
```

Run the tests:

```bash
/usr/local/go/bin/go test ./...
```

## Docker Compose

Run this command from the project root:

```bash
docker compose -f setup/docker/docker-compose.yml up -d
```

Compose starts the following services:

1. `db`: PostgreSQL
2. `migrate`: Runs `migrate up` on the migrations in `setup/migrations`
3. `api`: Builds and runs the API multi-stage image

The API service starts only after the migration service exits with code zero.
If the migration fails, the API does not start.

Check the status:

```bash
docker compose -f setup/docker/docker-compose.yml ps -a
```

Expected status:

```text
db       Up (healthy)
migrate  Exited (0)
api      Up
```

View the logs:

```bash
docker compose -f setup/docker/docker-compose.yml logs -f db
docker compose -f setup/docker/docker-compose.yml logs migrate
docker compose -f setup/docker/docker-compose.yml logs -f api
```

Stop the services:

```bash
docker compose -f setup/docker/docker-compose.yml down
```

To also remove the database volume:

```bash
docker compose -f setup/docker/docker-compose.yml down -v
```

### Changing Ports

Default values:

- API: `8080`
- PostgreSQL: `5432`

If the ports are already in use:

```bash
POSTGRES_PORT=5544 API_PORT=8081 \
docker compose -f setup/docker/docker-compose.yml up -d
```

In this case, the API is available at:

```text
http://localhost:8081
```

### Docker Compose Variables

| Variable | Default | Purpose |
|---|---:|---|
| `POSTGRES_USER` | `postgres` | PostgreSQL username |
| `POSTGRES_PASSWORD` | `postgres` | PostgreSQL password |
| `POSTGRES_DB` | `taskmanager` | Database name |
| `POSTGRES_PORT` | `5432` | PostgreSQL host port |
| `REDIS_PORT` | `6379` | Redis host port |
| `API_PORT` | `8080` | API host port |

## API

Base URL:

```text
http://localhost:8080/api/v1/tasks
```

### Prometheus metrics

Prometheus metrics are exposed at:

```bash
curl http://localhost:8080/metrics
```

The API publishes:

- `requests_total{method,route,status}`: total HTTP requests
- `request_latency_histogram{method,route}`: request latency in seconds
- `tasks_count`: current number of rows in the `tasks` table

The HTTP metrics use stable route templates such as `/api/v1/tasks/:id`
instead of raw URLs, which prevents a separate time series for every task ID.
The `tasks_count` gauge is initialized from PostgreSQL at startup, incremented
after a successful task creation, and decremented after a successful deletion.
It therefore reflects the task mutations handled by the running API process.

To scrape the API with Prometheus, add this job to `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: taskmanager
    static_configs:
      - targets: ["host.docker.internal:8080"]
```

When Prometheus runs directly on the host, use `localhost:8080` instead.
After updating the configuration, start Prometheus and open its targets page
to verify that the `taskmanager` target is `UP`.

In the examples below, replace `8080` with your API port if Compose is running
with a custom port.

### Create a task

```bash
curl -i -X POST http://localhost:8080/api/v1/tasks \
  -H 'Content-Type: application/json' \
  -d '{
    "title": "Read documentation",
    "description": "Review the project README"
  }'
```

Example response:

```json
{
  "id": "0198f6c0-8c6c-7abc-9d2d-123456789abc",
  "title": "Read documentation",
  "description": "Review the project README",
  "status": "TODO",
  "priority": "MEDIUM",
  "createdAtUnixSec": 1788950000
}
```

### Get a task

```bash
curl -i http://localhost:8080/api/v1/tasks/0198f6c0-8c6c-7abc-9d2d-123456789abc
```

Example response:

```json
{
  "id": "0198f6c0-8c6c-7abc-9d2d-123456789abc",
  "title": "Read documentation",
  "description": "Review the project README",
  "status": "TODO",
  "priority": "MEDIUM",
  "createdAtUnixSec": 1788950000
}
```

### List tasks

Query parameters:

- `status`: Filter by `TODO`, `IN_PROGRESS`, or `DONE`
- `limit`: Number of items, between 1 and 100
- `last_id`: ID of the last task on the previous page for cursor pagination

```bash
curl -i \
  'http://localhost:8080/api/v1/tasks?status=TODO&limit=20'
```

Example response:

```json
{
  "tasks": [
    {
      "id": "0198f6c0-8c6c-7abc-9d2d-123456789abc",
      "title": "Read documentation",
      "description": "Review the project README",
      "status": "TODO",
      "priority": "MEDIUM",
      "createdAtUnixSec": 1788950000
    }
  ]
}
```

For the next page, send the last item's `id` as `last_id`:

```bash
curl -i \
  'http://localhost:8080/api/v1/tasks?status=TODO&limit=20&last_id=0198f6c0-8c6c-7abc-9d2d-123456789abc'
```

### Update a task

```bash
curl -i -X PUT \
  http://localhost:8080/api/v1/tasks/0198f6c0-8c6c-7abc-9d2d-123456789abc \
  -H 'Content-Type: application/json' \
  -d '{
    "title": "Read updated documentation",
    "status": "DONE",
    "priority": "HIGH"
  }'
```

Example response:

```json
{
  "id": "0198f6c0-8c6c-7abc-9d2d-123456789abc",
  "title": "Read updated documentation",
  "description": "Review the project README",
  "status": "DONE",
  "priority": "HIGH",
  "createdAtUnixSec": 1788950000,
  "updatedAtUnixSec": 1788950100
}
```

### Delete a task

```bash
curl -i -X DELETE \
  http://localhost:8080/api/v1/tasks/0198f6c0-8c6c-7abc-9d2d-123456789abc
```

Example response:

```json
{
  "id": "0198f6c0-8c6c-7abc-9d2d-123456789abc",
  "deleted": true
}
```

## Architecture

Project structure:

```text
cmd/api/                         Application entry point
config/                          Environment configuration
internal/
  domain/task/                   Task entity and business rules
  app/
    dto/                         Request and response DTOs
    port/                        Input and output interfaces
    service/task/                Application services
  infra/
    http/                        Routes, binding, and middleware
    persistence/postgres/        PostgreSQL adapter
pkg/
  metrics/                       Prometheus collectors and HTTP middleware
setup/
  migrations/                    Versioned migrations
  docker/                        Dockerfile, Compose, and migration script
```

### Hexagonal Architecture

The application logic is separated from the core and its adapters:

- **Domain**: Located in `internal/domain/task`; it contains business rules such as
  permissions, status, and priority.
- **Application**: Located in `internal/app/service`; it executes use cases and
  uses the interfaces defined in `internal/app/port`.
- **Ports**: `TaskService` and `TaskRepo` define the contracts between the core
  and the infrastructure.
- **Adapters**:
  - HTTP in `internal/infra/http`
  - PostgreSQL in `internal/infra/persistence/postgres`

This keeps application services independent of Gin and PostgreSQL. The adapters
are connected to the application through dependency injection.

Prometheus metrics follow the same composition-root approach. The metrics
component is constructed in `cmd/api/main.go` with Fx, receives the PostgreSQL
executor through its port, and is injected into the HTTP router. Business
services do not import Prometheus directly.

### Domain-Driven Design

- `task.Entity` is the main domain model.
- `Status` and `Priority` are domain types.
- Task mutation rules are enforced inside the entity, not in the HTTP handler.
- DTOs are separate from the domain model and are used for the API contract.
- The repository interface is defined in the application layer, while its
  implementation is located in the infrastructure layer.

## Migrations

Migrations are stored in the following directory with versioned names:

```text
setup/migrations/
  000001_tasks_table.up.sql
  000001_tasks_table.down.sql
```

In Docker, migrations are not run manually. The `migrate/migrate` service runs
them with the `migrate up` command. Migration state is stored in the internal
`schema_migrations` table, so successful migrations are not run again.

## Cache

The API uses Redis with a cache-aside strategy for `GET /api/v1/tasks/:id`:

1. The application checks Redis using a key scoped by both `user_id` and task ID.
2. On a cache miss, it reads the task from PostgreSQL and stores the response in
   Redis for five minutes.
3. After a successful task update or delete, the corresponding Redis entry is
   invalidated.

Redis is started automatically by Docker Compose. Runtime cache failures are
logged and do not prevent the request from falling back to PostgreSQL.

## Integration Tests

The HTTP tests use a real PostgreSQL instance in Docker. Run them with:

```bash
/usr/local/go/bin/go test ./internal/infra/http -count=1
```

Before running, the test starts PostgreSQL and applies the migrations. After the
test completes, it removes the container.

The metrics package also includes a unit test for request counters, latency
histograms, and the task-count collector. Run all tests with:

```bash
/usr/local/go/bin/go test ./... -count=1
```
