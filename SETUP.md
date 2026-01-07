# Setup Instructions

## Step 1: Start MySQL Database

You have two options:

### Option A: Using Docker (Recommended - Easiest)

1. **Start MySQL container:**
   ```bash
   docker-compose up -d
   ```

2. **Verify MySQL is running:**
   ```bash
   docker-compose ps
   ```

3. **Check logs if needed:**
   ```bash
   docker-compose logs mysql
   ```

The `.env` file is already configured to work with the Docker MySQL container:
- Host: `localhost`
- Port: `3306`
- User: `root`
- Password: `password`
- Database: `movie_booking`

### Option B: Install MySQL Locally

If you prefer to install MySQL directly on your Mac:

```bash
brew install mysql
brew services start mysql
mysql -u root -p
# Create database: CREATE DATABASE movie_booking;
```

Then update `.env` with your MySQL credentials.

## Step 2: Run Database Migrations

Once MySQL is running, run the migrations:

```bash
export PATH="/opt/homebrew/bin:$PATH"
go run cmd/main.go --migrate --migration-command=up
```

This will create all the necessary tables:
- `users`
- `movies`
- `theatres`
- `shows`
- `show_seats`
- `bookings`

## Step 3: Start the API Server

```bash
export PATH="/opt/homebrew/bin:$PATH"
go run cmd/main.go --api
```

The server will start on `http://localhost:8080`

## Step 4: Test the API

### Health Check
```bash
curl http://localhost:8080/health
```

### Login (First, you'll need to create a user - see below)
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

## Creating a Test User

Since there's no user registration endpoint in the MVP, you'll need to create a user directly in the database:

```bash
# Connect to MySQL
docker exec -it movie-booking-mysql mysql -uroot -ppassword movie_booking

# Create a test user (password is "password123" hashed with bcrypt)
# You can generate the hash using Go or use this pre-hashed value:
INSERT INTO users (email, password_hash, name, created_at, updated_at) 
VALUES ('test@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Test User', NOW(), NOW());
```

Or use this Go one-liner to generate a bcrypt hash:
```bash
go run -c 'package main; import ("fmt"; "golang.org/x/crypto/bcrypt"); func main() { hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), 10); fmt.Println(string(hash)) }'
```

## Troubleshooting

### MySQL Connection Issues

1. **Check if MySQL is running:**
   ```bash
   docker-compose ps
   ```

2. **Check MySQL logs:**
   ```bash
   docker-compose logs mysql
   ```

3. **Verify connection:**
   ```bash
   docker exec -it movie-booking-mysql mysql -uroot -ppassword -e "SELECT 1;"
   ```

### Migration Issues

If migrations fail:
1. Make sure MySQL is running
2. Check your `.env` file has correct credentials
3. Try running migrations again: `go run cmd/main.go --migrate --migration-command=up`

### Port Already in Use

If port 8080 is already in use, change `SERVER_PORT` in `.env` to a different port (e.g., `8081`).

## Stopping Services

To stop MySQL:
```bash
docker-compose down
```

To stop MySQL and remove data:
```bash
docker-compose down -v
```
