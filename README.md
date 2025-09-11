## radius-server

Backend service with PostgreSQL using Fiber, GORM, and zerolog. Includes an integration-style test suite under `tests/` that initializes the database and required fixtures before running.

### Requirements
- Go (as in `go.mod`): 1.24.x
- Docker and Docker Compose
- PostgreSQL (provided via Docker Compose)

### Project layout
- `src/`: application code
- `src/database/migrations`: SQL migrations (golang-migrate compatible)
- `tests/`: test suite (ordered suite supported)
- `docker/`: local infra
- `Makefile`: helper targets

### Setup
1) Copy and fill environment variables

Create `.env` in the project root. You can use `template.env` as a starting point.

Key variables:
- `SECURITY_X_API_KEY` (required)
- `DB_DNS` (defaults to local compose DB)
- `DB_AUTO_RUN_MIGRATION` (defaults to true)

2) Start dependencies (PostgreSQL)

From the project root:
```bash
make docker-up
```

Useful alternatives:
```bash
make docker-up-clean
make docker-down
make docker-prune
```

### Run the application
```bash
go run main.go
```

On startup the app will:
- load `.env`
- connect to the database via `DB_DNS`
- run migrations if `DB_AUTO_RUN_MIGRATION=true`

### Running tests
Tests live in `./tests`. The test bootstrap (`TestMain`) ensures:
- working directory is set to the module root so `./.env` is used
- logger is initialized
- database connection is established
- a temporary NAS test fixture is created and cleaned up

Run all tests:
```bash
make run-tests
```

#### Ordered test suite
There is a single orchestrator test (`TestOrderedSuite`) that runs subtests in a specific order. To run only the ordered suite:

Direct command:
```bash
go test -v ./tests -run "^TestOrderedSuite$"
```

Makefile target:
```bash
make run-tests
```

#### Running a specific test or subtest
```bash
go test -v ./tests -run "^TestOrderedSuite$/Auth"
```

### Troubleshooting
- `.env` not picked up in tests? The `TestMain` changes the working directory to the module root by locating `go.mod`. Ensure `.env` is at the project root and readable.
- DB connection errors: confirm Docker is up and `DB_DNS` points to the compose database; defaults typically look like `postgres://postgres:postgres@localhost:5532/postgres?sslmode=disable`.

### Make targets
```make
docker-up          # Start local infra
docker-up-clean    # Start infra removing orphans
docker-down        # Stop infra
docker-prune       # Prune docker system (dangerous)
run-tests          # Run ordered test suite
```

### Migrations
Migrations are applied automatically on app startup when `DB_AUTO_RUN_MIGRATION=true`. During tests, the same flag controls whether migrations run before executing the suite.
