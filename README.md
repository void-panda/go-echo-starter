# Go Echo Starter

A production-ready Go starter template using Echo framework with layered architecture.

## Tech Stack

- **Framework**: Echo v4
- **Database**: PostgreSQL (pgx driver)
- **Logger**: Zerolog
- **Validation**: go-playground/validator
- **Documentation**: Swagger

## Project Structure

```
├── cmd/api/          # Application entry point
├── internal/
│   ├── config/       # Configuration
│   ├── database/     # Database connection
│   ├── domain/       # Domain entities
│   ├── repository/   # Data access layer
│   ├── service/      # Business logic layer
│   ├── handler/      # HTTP handlers
│   └── middleware/   # Custom middlewares
├── pkg/
│   ├── logger/       # Structured logger
│   ├── response/     # Response helpers
│   └── validator/    # Custom validator
└── docs/             # Swagger documentation
```

## Getting Started

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Make (optional)

### Setup

1. Clone the repository
2. Copy environment file:
   ```bash
   cp .env.example .env
   ```
3. Start PostgreSQL:
   ```bash
   docker-compose up -d
   ```
4. Run migrations:
   ```bash
   make migrate
   ```
5. Run the application:
   ```bash
   make run
   ```

## API Endpoints

| Method | Endpoint           | Description      |
|--------|-------------------|------------------|
| GET    | /health           | Health check     |
| GET    | /swagger/*        | API Documentation|
| GET    | /api/v1/users     | Get all users    |
| GET    | /api/v1/users/:id | Get user by ID   |
| POST   | /api/v1/users     | Create user      |
| PUT    | /api/v1/users/:id | Update user      |
| DELETE | /api/v1/users/:id | Delete user      |

## Available Commands

```bash
make run          # Run the application
make build        # Build binary
make test         # Run tests
make swagger      # Generate swagger docs
make docker-up    # Start PostgreSQL
make docker-down  # Stop PostgreSQL
make migrate      # Run migrations
```

## Environment Variables

| Variable    | Description           | Default     |
|-------------|-----------------------|-------------|
| APP_PORT    | Application port      | 8080        |
| APP_ENV     | Environment           | development |
| DB_HOST     | Database host         | localhost   |
| DB_PORT     | Database port         | 5432        |
| DB_USER     | Database user         | postgres    |
| DB_PASSWORD | Database password     | postgres    |
| DB_NAME     | Database name         | go_echo_db  |
| LOG_LEVEL   | Log level             | debug       |

## License

MIT
