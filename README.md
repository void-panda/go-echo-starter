# Go Echo Starter

A production-ready Go starter template featuring a robust, layered architecture (Clean Architecture principles), designed for scalability and maintainability. Powered by the high-performance [Echo](https://echo.labstack.com/) framework.

## 🚀 Features

- **Layered Architecture**: Clean separation of concerns (Handler, Service, Repository, Domain).
- **Database**: PostgreSQL integration using `sqlx` and `pgx` driver.
- **Migrations**: Automated database migrations using `golang-migrate`.
- **Identity**: Secure identification using **UUID v4** instead of auto-incrementing integers.
- **Authentication**: JWT-based authentication with high-security configuration.
- **Validation**: Comprehensive request validation using `go-playground/validator`.
- **Logging**: Structured, high-performance logging with `zerolog`.
- **API Documentation**: Interactive API docs with **Swagger UI**.
- **Development**: Hot-reloading enabled via [Air](https://github.com/cosmtrek/air).
- **Containerization**: Fully Dockerized development and production setups.
- **Testing**: Comprehensive Unit and Integration testing with `testify` and mocks.

## 📂 Project Structure

```text
├── cmd/api/            # Application entry point
├── internal/
│   ├── config/         # Configuration logic (Environment variables)
│   ├── database/       # DB connection and embedded migrations
│   ├── domain/         # Models, Request/Response DTOs
│   ├── repository/     # Data access layer
│   ├── service/        # Business logic layer
│   ├── handler/        # HTTP handlers / Controllers
│   └── middleware/     # Custom HTTP middlewares (JWT, CORS, etc.)
├── pkg/
│   ├── jwt/            # JWT Helper utilities
│   ├── logger/         # Structured logger wrapper
│   ├── response/       # Unified API response format
│   └── validator/      # Request validation logic
├── migrations/         # SQL migration files
├── docs/               # Generated Swagger documentation
├── docker-compose.yml  # Infrastructure setup (PostgreSQL)
└── Makefile            # Common development commands
```

## 🛠 Getting Started

### Prerequisites

- [Go 1.25+](https://go.dev/doc/install)
- [Docker & Docker Compose](https://docs.docker.com/get-docker/)
- [Make](https://www.gnu.org/software/make/) (recommended)

### Setup

1. **Clone the repository**
2. **Setup environment variables**
   ```bash
   cp .env.example .env
   ```
3. **Start Infrastructure (PostgreSQL)**
   ```bash
   make docker-up
   ```
4. **Run Migrations**
   ```bash
   make migrate
   ```
5. **Install Dependencies**
   ```bash
   make deps
   ```
6. **Run Application**
   ```bash
   make run
   # Or with hot-reload (if Air is installed):
   air
   ```

## 📖 API Documentation

The project includes built-in Swagger documentation. Once the server is running, access it at:

👉 [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

To update the documentation after making changes to handlers:
```bash
make swagger
```

## 🧪 Testing

Run all tests including unit and integration tests with mocks:

```bash
make test
```

To run a full health check (lint, tidy, and test):
```bash
make check
```

## 📜 Makefile Commands

| Command           | Description                                  |
|-------------------|----------------------------------------------|
| `make run`         | Start the application                       |
| `make build`       | Build the binary                            |
| `make test`        | Run all tests                               |
| `make check`       | Run lint, go mod tidy, and tests            |
| `make swagger`     | Regenerate Swagger documentation            |
| `make docker-up`   | Spin up PostgreSQL container                |
| `make docker-down` | Stop and remove containers                  |
| `make migrate`     | Run SQL migrations                          |
| `make clean`       | Remove build artifacts and temp files       |

## 🛡 License

Distributed under the MIT License. See `LICENSE` for more information.
