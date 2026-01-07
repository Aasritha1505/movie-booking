import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import ProtectedRoute from './components/ProtectedRoute';
import LoginPage from './pages/LoginPage';
import MoviesPage from './pages/MoviesPage';
import ShowsPage from './pages/ShowsPage';
import SeatsPage from './pages/SeatsPage';
import BookingConfirmationPage from './pages/BookingConfirmationPage';
import './App.css';

const App: React.FC = () => {
  return (
    <AuthProvider>
      <Router>
        <div className="App">
          <Routes>
            <Route path="/login" element={<LoginPage />} />
            <Route
              path="/movies"
              element={
                <ProtectedRoute>
                  <MoviesPage />
                </ProtectedRoute>
              }
            />
            <Route
              path="/movies/:movieId/shows"
              element={
                <ProtectedRoute>
                  <ShowsPage />
                </ProtectedRoute>
              }
            />
            <Route
              path="/shows/:showId/seats"
              element={
                <ProtectedRoute>
                  <SeatsPage />
                </ProtectedRoute>
              }
            />
            <Route
              path="/booking-confirmation/:bookingId"
              element={
                <ProtectedRoute>
                  <BookingConfirmationPage />
                </ProtectedRoute>
              }
            />
            <Route path="/" element={<Navigate to="/movies" replace />} />
          </Routes>
        </div>
      </Router>
    </AuthProvider>
  );
};

export default App;
