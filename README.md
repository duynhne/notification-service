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

```bash
go mod download
go test ./...
go run cmd/main.go
```

## License

MIT
