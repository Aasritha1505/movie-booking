# Quick Run Commands 🚀

## Prerequisites
```bash
# Make sure Go and MySQL are in PATH
export PATH="/opt/homebrew/bin:$PATH"
```

## 1. Start MySQL Database
```bash
# Check if MySQL is running
brew services list | grep mysql

# Start MySQL (if not running)
brew services start mysql

# Stop MySQL (if needed)
brew services stop mysql
```

## 2. Run Database Migrations (First Time Only)
```bash
cd /Users/aasrithayadav/movie-booking
export PATH="/opt/homebrew/bin:$PATH"
go run cmd/main.go --migrate --migration-command=up
```

## 3. Seed Test Data (First Time Only)
```bash
cd /Users/aasrithayadav/movie-booking
export PATH="/opt/homebrew/bin:$PATH"
go run tools/seed_data.go
```

## 4. Start Backend API Server
```bash
cd /Users/aasrithayadav/movie-booking
export PATH="/opt/homebrew/bin:$PATH"
go run cmd/main.go --api
```

**Backend runs on:** `http://localhost:8080`

## 5. Start Frontend (in a new terminal)
```bash
cd /Users/aasrithayadav/movie-booking/frontend
npm start
```

**Frontend runs on:** `http://localhost:3000`

## 6. View Database Data
```bash
cd /Users/aasrithayadav/movie-booking
export PATH="/opt/homebrew/bin:$PATH"
go run tools/view_data.go
```

## Quick Test Commands

### Test Backend Health
```bash
curl http://localhost:8080/health
```

### Test Login
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

### Test Get Movies
```bash
curl http://localhost:8080/api/v1/movies
```

## Complete Startup Sequence

**Terminal 1 - Backend:**
```bash
cd /Users/aasrithayadav/movie-booking
export PATH="/opt/homebrew/bin:$PATH"

# Start MySQL (if not running)
brew services start mysql

# Start API server
go run cmd/main.go --api
```

**Terminal 2 - Frontend:**
```bash
cd /Users/aasrithayadav/movie-booking/frontend
npm start
```

## Stop Services

### Stop Backend
Press `Ctrl+C` in the backend terminal

### Stop Frontend
Press `Ctrl+C` in the frontend terminal

### Stop MySQL
```bash
brew services stop mysql
```

## Troubleshooting

### Port 8080 Already in Use
```bash
# Find and kill process on port 8080
lsof -ti:8080 | xargs kill -9
```

### Port 3000 Already in Use
```bash
# Find and kill process on port 3000
lsof -ti:3000 | xargs kill -9

# Or use a different port
PORT=3001 npm start
```

### Database Connection Error
```bash
# Check MySQL is running
brew services list | grep mysql

# Test MySQL connection
/opt/homebrew/bin/mysql -u root -ppassword movie_booking
```

### Frontend Dependencies Not Installed
```bash
cd /Users/aasrithayadav/movie-booking/frontend
npm install
```

## All-in-One Startup Script

Create a file `start.sh`:
```bash
#!/bin/bash
export PATH="/opt/homebrew/bin:$PATH"

# Start MySQL
brew services start mysql
sleep 2

# Start backend in background
cd /Users/aasrithayadav/movie-booking
go run cmd/main.go --api > /tmp/backend.log 2>&1 &
echo "Backend started on http://localhost:8080"
echo "Logs: tail -f /tmp/backend.log"

# Start frontend
cd /Users/aasrithayadav/movie-booking/frontend
npm start
```

Make it executable:
```bash
chmod +x start.sh
./start.sh
```
