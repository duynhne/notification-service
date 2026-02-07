# notification-service

Notification microservice for email, SMS, and in-app notifications.

## Features

- Email notifications
- SMS notifications
- In-app notifications
- Mark as read

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/v1/notifications` | Get all |
| `GET` | `/api/v1/notifications/count` | Get unread count |
| `GET` | `/api/v1/notifications/:id` | Get by ID |
| `PATCH` | `/api/v1/notifications/:id` | Mark as read |
| `POST` | `/api/v1/notify/email` | Send email |
| `POST` | `/api/v1/notify/sms` | Send SMS |

## Tech Stack

- Go + Gin framework
- PostgreSQL 16 (supporting-db cluster, cross-namespace)
- PgBouncer connection pooling
- OpenTelemetry tracing

## Development

### Prerequisites

- Go 1.25+
- [golangci-lint](https://golangci-lint.run/welcome/install/) v2+

### Local Development

```bash
# Install dependencies
go mod tidy
go mod download

# Build
go build ./...

# Test
go test ./...

# Lint (must pass before PR merge)
golangci-lint run --timeout=10m

# Run locally (requires .env or env vars)
go run cmd/main.go
```

### Pre-push Checklist

```bash
go build ./... && go test ./... && golangci-lint run --timeout=10m
```
