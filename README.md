# Movie Booking Service

A high-performance movie seat booking system built with Go, following the SAFE Security engineering guidelines.

## Features

- **Concurrent Seat Locking**: Prevents double-booking using MySQL row-level locks (`SELECT ... FOR UPDATE`)
- **10-Minute Seat Holds**: Lazy lock expiration (no background cron jobs)
- **JWT Authentication**: Stateless authentication with 15-minute token expiry
- **Idempotent Bookings**: Support for retry-safe booking operations
- **Three-Layer Architecture**: API → Service → DataStore (strict separation of concerns)

## Architecture

This service follows the three-layer architecture pattern:

- **API Layer** (`api/v1/`): HTTP handlers, request parsing, response formatting
- **Service Layer** (`core/services/`): Business logic and orchestration
- **DataStore Layer** (`datastore/`): Database operations with GORM

## Prerequisites

- Go 1.22+
- MySQL 5.7+ or MySQL 8.0+
- Docker (optional, for containerized deployment)

## Setup

1. **Clone and navigate to the project:**
   ```bash
   cd movie-booking
   ```

2. **Copy environment configuration:**
   ```bash
   cp env.sample .env
   # Edit .env with your database credentials
   ```

3. **Install dependencies:**
   ```bash
   go mod download
   ```

4. **Run database migrations:**
   ```bash
   go run cmd/main.go --migrate --migration-command=up
   ```

5. **Start the API server:**
   ```bash
   go run cmd/main.go --api
   ```

The server will start on port 8080 (configurable via `SERVER_PORT`).

## API Endpoints

### Public Endpoints

- `POST /api/v1/login` - User authentication
- `GET /api/v1/movies` - List all movies
- `GET /api/v1/movies/{id}/shows` - Get shows for a movie
- `GET /api/v1/shows/{id}/seats` - Get seat grid for a show

### Protected Endpoints (Require JWT)

- `PATCH /api/v1/seats/{id}/lock` - Lock a seat for 10 minutes
- `POST /api/v1/bookings` - Create a booking (converts lock to sale)

## Usage Examples

### 1. Login

```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

Response:
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Login successful",
  "values": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "user@example.com"
    }
  }
}
```

### 2. Lock a Seat

```bash
curl -X PATCH http://localhost:8080/api/v1/seats/1/lock \
  -H "Authorization: Bearer <token>"
```

Response:
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Locked",
  "values": {
    "message": "Locked",
    "expires_at": "2024-01-01T10:10:00Z"
  }
}
```

### 3. Create Booking

```bash
curl -X POST http://localhost:8080/api/v1/bookings \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: unique-key-123" \
  -d '{"show_id":1,"seat_id":1}'
```

Response:
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Ticket sent to your email.",
  "values": {
    "booking_id": 1,
    "status": "CONFIRMED",
    "message": "Ticket sent to your email."
  }
}
```

## Concurrency Strategy

The core concurrency challenge is handled in the `LockSeat` operation:

1. **Transaction begins** with row-level lock (`SELECT ... FOR UPDATE`)
2. **Validation**: Seat must be `AVAILABLE` or `LOCKED` with expired timestamp
3. **Update**: Set status to `LOCKED`, `locked_at = now()`, `user_id = caller`
4. **Commit**: Releases the lock

This ensures that only one user can lock a seat at a time, preventing double-booking.

## Testing

Run tests:
```bash
go test ./...
```

With coverage:
```bash
go test -cover ./...
```

## Project Structure

```
movie-booking/
├── cmd/
│   └── main.go                 # Entry point
├── api/v1/
│   ├── router.go              # Route registration
│   ├── controllers/           # HTTP handlers
│   ├── helpers/               # Request parsing, validation
│   ├── interceptors/          # Middleware (auth, logging, panic recovery)
│   └── types/                 # Request/Response DTOs
├── core/
│   ├── model/                 # GORM models
│   ├── services/              # Business logic
│   └── types/                 # Shared types
├── datastore/                 # Database access layer
├── config/                    # Viper configuration
├── constants/                 # Constants (no hardcoding!)
├── dbmigrations/
│   └── migrations/mysql/      # Goose migration files
└── util/                      # Utilities (context, errors)
```

## Configuration

All configuration is managed via environment variables (see `env.sample`):

- `DATABASE_HOST`, `DATABASE_PORT`, `DATABASE_USER`, `DATABASE_PASSWORD`, `DATABASE_NAME`
- `SERVER_PORT` (default: 8080)
- `HANDLER_TIMEOUT` (default: 30s)
- `JWT_SECRET` (change in production!)
- `JWT_EXPIRY` (default: 15m)
- `SEAT_LOCK_DURATION` (default: 10m)

## Development Guidelines

This project follows the SAFE Security Go development guidelines:

- **Three-layer architecture**: Never skip layers
- **Error handling**: Always wrap errors with context
- **Logging**: Structured logging with logrus
- **Security**: Parameterized queries only, JWT RS256 (HS256 for MVP)
- **Testing**: Unit tests with mocks, integration tests with real DB
- **Constants**: No hardcoded strings

## License

[Add your license here]
