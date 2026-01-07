# Local MySQL Setup - Complete! ✅

## What Was Done

1. ✅ **MySQL Installed** - MySQL 9.5.0 installed via Homebrew
2. ✅ **MySQL Started** - Service running and auto-starts on login
3. ✅ **Database Created** - `movie_booking` database created
4. ✅ **Password Set** - Root password set to `password` (matches `.env`)
5. ✅ **Migrations Run** - All tables created successfully:
   - `users`
   - `movies`
   - `theatres`
   - `shows`
   - `show_seats`
   - `bookings`

## Quick Start Commands

### Start the API Server
```bash
export PATH="/opt/homebrew/bin:$PATH"
go run cmd/main.go --api
```

The server will start on `http://localhost:8080`

### Check MySQL Status
```bash
brew services list | grep mysql
```

### Stop MySQL (if needed)
```bash
brew services stop mysql
```

### Start MySQL (if stopped)
```bash
brew services start mysql
```

### Connect to MySQL
```bash
/opt/homebrew/bin/mysql -u root -ppassword movie_booking
```

### Run Migrations Again (if needed)
```bash
export PATH="/opt/homebrew/bin:$PATH"
go run cmd/main.go --migrate --migration-command=up
```

### Check Migration Status
```bash
export PATH="/opt/homebrew/bin:$PATH"
go run cmd/main.go --migrate --migration-command=status
```

## Test the API

### 1. Health Check
```bash
curl http://localhost:8080/health
```

### 2. Create a Test User (Required for login)

**Option A: Using the helper script (Recommended)**
```bash
export PATH="/opt/homebrew/bin:$PATH"
go run tools/create_user.go test@example.com password123 "Test User"
```

**Option B: Manual SQL (if needed)**
```bash
# Generate a password hash first
export PATH="/opt/homebrew/bin:$PATH"
HASH=$(go run tools/generate_hash.go | grep "^Hash:" | cut -d' ' -f2)

# Then insert into database
/opt/homebrew/bin/mysql -u root -ppassword movie_booking -e "INSERT INTO users (email, password_hash, name, created_at, updated_at) VALUES ('test@example.com', '$HASH', 'Test User', NOW(), NOW());"
```

**Note:** The test user `test@example.com` with password `password123` has already been created and should work for testing.

### 3. Login
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

### 4. Get Movies
```bash
curl http://localhost:8080/api/v1/movies
```

## Configuration

Your `.env` file is configured for local MySQL:
- **Host**: `localhost`
- **Port**: `3306`
- **User**: `root`
- **Password**: `password`
- **Database**: `movie_booking`

## Troubleshooting

### MySQL Not Running
```bash
brew services start mysql
```

### Can't Connect to Database
1. Check MySQL is running: `brew services list | grep mysql`
2. Verify password: Try connecting manually: `/opt/homebrew/bin/mysql -u root -ppassword`
3. Check `.env` file has correct credentials

### Port 8080 Already in Use
Edit `.env` and change `SERVER_PORT=8080` to a different port (e.g., `8081`)

### Reset Database (DANGER: Deletes all data!)
```bash
/opt/homebrew/bin/mysql -u root -ppassword -e "DROP DATABASE movie_booking; CREATE DATABASE movie_booking;"
export PATH="/opt/homebrew/bin:$PATH"
go run cmd/main.go --migrate --migration-command=up
```

## Next Steps

1. **Start the server**: `go run cmd/main.go --api`
2. **Create test data**: Add some movies, theatres, and shows via SQL or API
3. **Test the booking flow**: Login → Browse movies → Lock seat → Book seat

Enjoy your movie booking system! 🎬
