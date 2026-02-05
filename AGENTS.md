# notification-service

> AI Agent context for understanding this repository

## 📋 Overview

Notification microservice. Handles email, SMS, and in-app notifications.

## 🏗️ Architecture

```
notification-service/
├── cmd/main.go
├── config/config.go
├── db/migrations/sql/
├── internal/
│   ├── core/
│   │   ├── database.go
│   │   └── domain/
│   ├── logic/v1/service.go
│   └── web/v1/handler.go
├── middleware/
└── Dockerfile
```

## 🔌 API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/v1/notifications` | Get all notifications |
| `GET` | `/api/v1/notifications/:id` | Get notification by ID |
| `PATCH` | `/api/v1/notifications/:id` | Mark as read |
| `POST` | `/api/v1/notify/email` | Send email (internal) |
| `POST` | `/api/v1/notify/sms` | Send SMS (internal) |

## 📐 3-Layer Architecture

| Layer | Location | Responsibility |
|-------|----------|----------------|
| **Web** | `internal/web/v1/handler.go` | HTTP, validation |
| **Logic** | `internal/logic/v1/service.go` | Business rules (❌ NO SQL) |
| **Core** | `internal/core/` | Domain models, repositories |

## 🗄️ Database

| Component | Value |
|-----------|-------|
| **Cluster** | supporting-db (shared with user, shipping) |
| **PostgreSQL** | 16 |
| **HA** | Single instance |
| **Pooler** | PgBouncer Sidecar |
| **Endpoint** | `supporting-db-pooler.user.svc.cluster.local:5432` |
| **Pool Mode** | Transaction |
| **Cross-namespace** | Yes (cluster in `user` namespace) |

**Note:** Database cluster is in `user` namespace. Zalando Operator syncs credentials via cross-namespace secret.

## 🚀 Graceful Shutdown

**VictoriaMetrics Pattern:**
1. `/ready` → 503 when shutting down
2. Drain delay (5s)
3. Sequential: HTTP → Database → Tracer

## 🔧 Tech Stack

| Component | Technology |
|-----------|------------|
| **Framework** | Gin |
| **Database** | PostgreSQL 16 via pgx/v5 |
| **Tracing** | OpenTelemetry |

## 🛠️ Development

```bash
go mod download && go test ./... && go build ./cmd/main.go
```
