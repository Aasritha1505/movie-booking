# Movie Booking Frontend

React TypeScript frontend for the Movie Booking System.

## Features

- 🔐 **Authentication**: JWT-based login system
- 🎬 **Movie Browsing**: Browse available movies
- 🎭 **Show Selection**: Select showtimes for movies
- 🪑 **Interactive Seat Map**: Visual seat selection with real-time status
- 🎫 **Booking Confirmation**: Complete booking flow with confirmation

## Tech Stack

- **React 18** with TypeScript
- **React Router** for navigation
- **Axios** for API calls
- **Context API** for state management

## Setup

1. **Install dependencies:**
   ```bash
   cd frontend
   npm install
   ```

2. **Start development server:**
   ```bash
   npm start
   ```

   The app will open at `http://localhost:3000`

## Configuration

The frontend is configured to proxy API requests to `http://localhost:8080` (see `package.json` proxy setting).

To change the API URL, create a `.env` file:
```
REACT_APP_API_URL=http://localhost:8080
```

## Project Structure

```
frontend/
├── src/
│   ├── components/       # Reusable components
│   ├── pages/           # Page components
│   ├── services/        # API service layer
│   ├── types/           # TypeScript type definitions
│   ├── context/         # React Context providers
│   ├── utils/           # Utility functions
│   ├── App.tsx          # Main app component
│   └── index.tsx        # Entry point
├── public/              # Static assets
└── package.json
```

## User Flow

1. **Login** → Enter credentials
2. **Browse Movies** → View available movies
3. **Select Show** → Choose a showtime
4. **Select Seat** → Click on available seat (locks for 10 minutes)
5. **Confirm Booking** → Complete the booking
6. **Confirmation** → View booking confirmation

## Test Credentials

- Email: `test@example.com`
- Password: `password123`

## Development

- Hot reload is enabled
- TypeScript strict mode is enabled
- ESLint is configured

## Build for Production

```bash
npm run build
```

This creates an optimized production build in the `build/` directory.
