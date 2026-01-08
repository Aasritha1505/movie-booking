#!/bin/bash
# Quick one-line API test commands

API="http://localhost:8080"

# 1. Health Check
curl $API/health

# 2. Login (save token)
TOKEN=$(curl -s -X POST $API/api/v1/login -H "Content-Type: application/json" -d '{"email":"test@example.com","password":"password123"}' | grep -o '"token":"[^"]*' | cut -d'"' -f4)

# 3. Get Movies
curl -s $API/api/v1/movies | python3 -m json.tool

# 4. Get Shows for Movie 1
curl -s $API/api/v1/movies/1/shows | python3 -m json.tool

# 5. Get Seats for Show 8
curl -s $API/api/v1/shows/8/seats | python3 -m json.tool

# 6. Lock Seat 1
curl -s -X PATCH $API/api/v1/seats/1/lock -H "Authorization: Bearer $TOKEN" | python3 -m json.tool

# 7. Create Booking
curl -s -X POST $API/api/v1/bookings -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -H "Idempotency-Key: test-$(date +%s)" -d '{"show_id":8,"seat_id":1}' | python3 -m json.tool
