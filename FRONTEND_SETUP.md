# Frontend Setup Guide

## Quick Start

1. **Navigate to frontend directory:**
   ```bash
   cd frontend
   ```

2. **Install dependencies:**
   ```bash
   npm install
   ```

3. **Start the development server:**
   ```bash
   npm start
   ```

   The app will automatically open at `http://localhost:3000`

## Prerequisites

- Node.js 16+ and npm installed
- Backend API running on `http://localhost:8080`

## Features Implemented

### ✅ Authentication
- Login page with email/password
- JWT token storage in localStorage
- Protected routes (requires authentication)
- Auto-logout on token expiration

### ✅ Movie Browsing
- List all available movies
- Movie cards with details (title, rating, duration, description)
- Click to view shows for a movie

### ✅ Show Selection
- List all shows for a selected movie
- Display theatre name, location, and showtime
- Click to view seat map

### ✅ Seat Selection
- **Visual seat grid** organized by rows (A, B, C, etc.)
- **Real-time status**:
  - 🟢 Available (green)
  - 🔵 Selected/Your Lock (blue)
  - 🟠 Locked by others (orange)
  - ⚪ Sold (gray)
- **Auto-refresh** every 5 seconds to show updated seat status
- **10-minute lock** when seat is selected
- **Lock expiration** shown in UI

### ✅ Booking Flow
- Lock seat → Confirm booking → View confirmation
- Idempotency support (prevents duplicate bookings)
- Automatic redirect after booking

## Project Structure

```
frontend/
├── src/
│   ├── components/
│   │   └── ProtectedRoute.tsx    # Route protection component
│   ├── pages/
│   │   ├── LoginPage.tsx         # Login page
│   │   ├── MoviesPage.tsx        # Movie listing
│   │   ├── ShowsPage.tsx         # Show selection
│   │   ├── SeatsPage.tsx         # Seat selection with grid
│   │   └── BookingConfirmationPage.tsx
│   ├── services/
│   │   └── api.ts                # API service layer
│   ├── types/
│   │   └── index.ts              # TypeScript types
│   ├── context/
│   │   └── AuthContext.tsx      # Authentication context
│   ├── App.tsx                   # Main app with routing
│   └── index.tsx                 # Entry point
├── public/
│   └── index.html
├── package.json
└── tsconfig.json
```

## API Integration

The frontend communicates with the backend API:

- **Base URL**: `http://localhost:8080` (configurable via `.env`)
- **Authentication**: JWT token in `Authorization: Bearer <token>` header
- **Error Handling**: Automatic token cleanup on 401 errors

## Styling

- Modern gradient design
- Responsive layout
- Interactive seat grid with hover effects
- Color-coded seat status indicators

## Testing the Flow

1. **Start backend** (in another terminal):
   ```bash
   cd /Users/aasrithayadav/movie-booking
   export PATH="/opt/homebrew/bin:$PATH"
   go run cmd/main.go --api
   ```

2. **Start frontend**:
   ```bash
   cd frontend
   npm start
   ```

3. **Test the flow**:
   - Login with `test@example.com` / `password123`
   - Browse movies
   - Select a movie → Select a show → Select a seat
   - Lock seat (10-minute hold)
   - Confirm booking
   - View confirmation

## Troubleshooting

### Port 3000 Already in Use
Change the port:
```bash
PORT=3001 npm start
```

### API Connection Issues
1. Verify backend is running: `curl http://localhost:8080/health`
2. Check `.env` file has correct API URL
3. Check browser console for CORS errors

### TypeScript Errors
Run type checking:
```bash
npx tsc --noEmit
```

## Build for Production

```bash
npm run build
```

This creates an optimized build in `build/` directory that can be served by any static file server.
