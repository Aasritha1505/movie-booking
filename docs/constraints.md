# Constraints and requirements (extracted from PDFs)

This repo implements the MVP described in `docs/movie_hld.txt`, following the conventions in `docs/go-onboarding-guide.txt`.

## HLD requirements (movie booking)

- **Architecture**: containerized **Go monolith**.
- **API**: **REST**.
- **DB**: **MySQL**.
- **Concurrency core**: prevent double-booking via `SELECT ... FOR UPDATE` on `show_seats` rows.
- **Seat holds**: **10-minute** hold window using **lazy lock expiration** (validate + clear expired locks during read/write, not a cron).
- **Auth**: stateless **JWT** login; HLD calls out `Authorization` header required for **lock** and **book** operations.
- **Endpoints (MVP)**:
  - `POST /api/v1/login`
  - `GET /api/v1/movies`
  - `GET /api/v1/movies/:id/shows`
  - `GET /api/v1/shows/:id/seats`
  - `PATCH /api/v1/seats/:id/lock`
  - `POST /api/v1/bookings`
- **Data model** (as shown in HLD): `Movie`, `Theatre`, `Show`, `User`, `ShowSeat`, `Booking`.

## Go onboarding guide requirements (repo standards)

- **Go version**: **Go 1.24+** (per guide).
- **Key libs**: Gorilla Mux, GORM, Viper, Logrus, Goose, gomock/mockgen, testify, OpenTelemetry.
- **Folder layout** (service-agnostic standard):
  - `cmd/main.go` with `--api` / `--migrate` flags
  - `api/v1/` (router, controllers, helpers, types)
  - `core/` (services, model, types)
  - `datastore/` (store implementations + fakes)
  - `config/` (viper getters)
  - `constants/` (no hardcoding)
  - `dbmigrations/migrations/mysql/` (goose migrations)
  - `tests/` (integration, e2e/storm)
- **Architecture rule**: data flows down, errors flow up; don’t skip layers.
- **Error handling**: wrap errors; don’t ignore errors; no panics in production code.
- **Logging**: structured logging (logrus style); include a `TAG` string in handlers for traceability.
- **Security**: parameterized DB queries only; validate inputs; JWT guidance in the guide mentions RS256 + 15-minute expiry.

