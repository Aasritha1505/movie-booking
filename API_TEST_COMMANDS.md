# API Test Commands 🧪

Complete curl commands to test all endpoints in the Movie Booking API.

## Prerequisites

Set the base URL:
```bash
export API_URL="http://localhost:8080"
```

---

## 1. Health Check

```bash
curl $API_URL/health
```

**Expected Response:**
```json
{"status":"ok"}
```

---

## 2. Login (Get JWT Token)

```bash
curl -X POST $API_URL/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

**Save the token for protected endpoints:**
```bash
# Save token to variable
TOKEN=$(curl -s -X POST $API_URL/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}' \
  | grep -o '"token":"[^"]*' | cut -d'"' -f4)

echo "Token: $TOKEN"
```

**Expected Response:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Login successful",
  "values": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "name": "Test User",
      "email": "test@example.com"
    }
  }
}
```

---

## 3. Get All Movies

```bash
curl $API_URL/api/v1/movies
```

**With pretty JSON:**
```bash
curl -s $API_URL/api/v1/movies | python3 -m json.tool
```

**Expected Response:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Movies retrieved successfully",
  "values": [
    {
      "id": 1,
      "title": "The Matrix",
      "description": "...",
      "duration_mins": 136,
      "content_rating": "R"
    }
  ]
}
```

---

## 4. Get Shows for a Movie

```bash
# Replace {movie_id} with actual movie ID (e.g., 1)
curl $API_URL/api/v1/movies/1/shows
```

**With pretty JSON:**
```bash
curl -s $API_URL/api/v1/movies/1/shows | python3 -m json.tool
```

**Expected Response:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Shows retrieved successfully",
  "values": [
    {
      "id": 8,
      "movie_id": 1,
      "theatre_id": 3,
      "start_time": "2026-01-08T03:05:24+05:30",
      "movie": {
        "id": 1,
        "title": "The Matrix"
      },
      "theatre": {
        "id": 3,
        "name": "PVR Cinemas",
        "location": "Downtown Mall"
      }
    }
  ]
}
```

---

## 5. Get Seats for a Show

```bash
# Replace {show_id} with actual show ID (e.g., 8)
curl $API_URL/api/v1/shows/8/seats
```

**With pretty JSON:**
```bash
curl -s $API_URL/api/v1/shows/8/seats | python3 -m json.tool
```

**Expected Response:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Seats retrieved successfully",
  "values": [
    {
      "id": 1,
      "show_id": 8,
      "seat_name": "A1",
      "status": "AVAILABLE",
      "locked_at": null,
      "user_id": null
    }
  ]
}
```

---

## 6. Lock a Seat (Requires Authentication)

**First, get a token:**
```bash
TOKEN=$(curl -s -X POST $API_URL/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}' \
  | grep -o '"token":"[^"]*' | cut -d'"' -f4)
```

**Then lock a seat:**
```bash
# Replace {seat_id} with actual seat ID (e.g., 1)
curl -X PATCH $API_URL/api/v1/seats/1/lock \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN"
```

**With pretty JSON:**
```bash
curl -s -X PATCH $API_URL/api/v1/seats/1/lock \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  | python3 -m json.tool
```

**Expected Response:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Seat locked successfully",
  "values": {
    "message": "Seat locked successfully",
    "expires_at": "2026-01-08T12:20:00+05:30"
  }
}
```

---

## 7. Create Booking (Requires Authentication + Locked Seat)

**First, get a token and lock a seat:**
```bash
# Get token
TOKEN=$(curl -s -X POST $API_URL/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}' \
  | grep -o '"token":"[^"]*' | cut -d'"' -f4)

# Lock a seat (replace seat_id and show_id)
curl -X PATCH $API_URL/api/v1/seats/1/lock \
  -H "Authorization: Bearer $TOKEN"
```

**Then create booking:**
```bash
# Generate unique idempotency key
IDEMPOTENCY_KEY="booking-$(date +%s)"

curl -X POST $API_URL/api/v1/bookings \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Idempotency-Key: $IDEMPOTENCY_KEY" \
  -d '{
    "show_id": 8,
    "seat_id": 1
  }'
```

**With pretty JSON:**
```bash
curl -s -X POST $API_URL/api/v1/bookings \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Idempotency-Key: $IDEMPOTENCY_KEY" \
  -d '{"show_id":8,"seat_id":1}' \
  | python3 -m json.tool
```

**Expected Response:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Booking created successfully",
  "values": {
    "booking_id": 1,
    "status": "confirmed",
    "message": "Booking created successfully"
  }
}
```

---

## Complete Test Script

Save this as `test_all_apis.sh`:

```bash
#!/bin/bash

API_URL="http://localhost:8080"

echo "=== Testing Movie Booking API ==="
echo ""

# 1. Health Check
echo "1. Health Check..."
curl -s $API_URL/health
echo -e "\n"

# 2. Login
echo "2. Login..."
LOGIN_RESPONSE=$(curl -s -X POST $API_URL/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}')

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)
echo "Token obtained: ${TOKEN:0:20}..."
echo ""

# 3. Get Movies
echo "3. Get Movies..."
curl -s $API_URL/api/v1/movies | python3 -m json.tool | head -20
echo ""

# 4. Get Shows (assuming movie ID 1)
echo "4. Get Shows for Movie ID 1..."
curl -s $API_URL/api/v1/movies/1/shows | python3 -m json.tool | head -30
echo ""

# 5. Get Seats (assuming show ID 8)
echo "5. Get Seats for Show ID 8..."
SEATS_RESPONSE=$(curl -s $API_URL/api/v1/shows/8/seats)
SEAT_ID=$(echo $SEATS_RESPONSE | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
echo "Found seat ID: $SEAT_ID"
echo ""

# 6. Lock Seat
echo "6. Lock Seat ID $SEAT_ID..."
curl -s -X PATCH $API_URL/api/v1/seats/$SEAT_ID/lock \
  -H "Authorization: Bearer $TOKEN" | python3 -m json.tool
echo ""

# 7. Create Booking
echo "7. Create Booking..."
IDEMPOTENCY_KEY="test-$(date +%s)"
curl -s -X POST $API_URL/api/v1/bookings \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Idempotency-Key: $IDEMPOTENCY_KEY" \
  -d "{\"show_id\":8,\"seat_id\":$SEAT_ID}" | python3 -m json.tool
echo ""

echo "=== All tests completed ==="
```

**Make it executable and run:**
```bash
chmod +x test_all_apis.sh
./test_all_apis.sh
```

---

## Quick Reference

| Endpoint | Method | Auth Required | Description |
|----------|--------|---------------|-------------|
| `/health` | GET | No | Health check |
| `/api/v1/login` | POST | No | User login |
| `/api/v1/movies` | GET | No | List all movies |
| `/api/v1/movies/{id}/shows` | GET | No | Get shows for movie |
| `/api/v1/shows/{id}/seats` | GET | No | Get seats for show |
| `/api/v1/seats/{id}/lock` | PATCH | Yes | Lock a seat |
| `/api/v1/bookings` | POST | Yes | Create booking |

---

## Error Testing

### Test Invalid Login
```bash
curl -X POST $API_URL/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"wrong@example.com","password":"wrong"}'
```

### Test Lock Without Auth
```bash
curl -X PATCH $API_URL/api/v1/seats/1/lock
```

### Test Lock Already Sold Seat
```bash
# First, book a seat, then try to lock it again
TOKEN=$(curl -s -X POST $API_URL/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}' \
  | grep -o '"token":"[^"]*' | cut -d'"' -f4)

# Try to lock an already sold seat
curl -X PATCH $API_URL/api/v1/seats/1/lock \
  -H "Authorization: Bearer $TOKEN"
```

---

## Tips

1. **Use jq for better JSON formatting:**
   ```bash
   curl -s $API_URL/api/v1/movies | jq
   ```

2. **Save responses to files:**
   ```bash
   curl -s $API_URL/api/v1/movies > movies.json
   ```

3. **Verbose output for debugging:**
   ```bash
   curl -v $API_URL/api/v1/movies
   ```

4. **Check response headers:**
   ```bash
   curl -i $API_URL/api/v1/movies
   ```
