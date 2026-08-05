# Personal ERP — Backend API (Go)

REST API for the Daily Activity module. Uses Go 1.22+ stdlib routing, MySQL/MariaDB.

## Requirements

- Go 1.26+
- MySQL/MariaDB 8+ running on `127.0.0.1:3306`

## Setup

```bash
cp .env.example .env   # adjust credentials as needed
go run ./cmd/server
```

The server creates the `personal_erp` database schema and seeds default categories on startup.

## Endpoints

| Method | Path                      | Description                              |
|--------|---------------------------|------------------------------------------|
| GET    | /api/health               | liveness + db check                      |
| GET    | /api/categories           | list categories                          |
| POST   | /api/categories           | create category                          |
| PUT    | /api/categories/{id}      | update category                          |
| DELETE | /api/categories/{id}      | delete category                          |
| GET    | /api/activities           | list activities (filters below)          |
| GET    | /api/activities/{id}      | get one activity                         |
| POST   | /api/activities           | create activity                          |
| PUT    | /api/activities/{id}      | update activity                          |
| DELETE | /api/activities/{id}      | delete activity                          |
| GET    | /api/activities/stats     | daily/range statistics                   |

List filters: `date`, `start`, `end`, `category_id`, `status`, `priority`, `q`.

See [`docs/PRD_Daily_Activity.md`](../docs/PRD_Daily_Activity.md) for the full spec.
