import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { apiService } from '../services/api';
import { Movie } from '../types';
import './MoviesPage.css';

const MoviesPage: React.FC = () => {
  const [movies, setMovies] = useState<Movie[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    loadMovies();
  }, []);

  const loadMovies = async () => {
    try {
      setLoading(true);
      const response = await apiService.getMovies();
      if (response.success && response.values) {
        setMovies(response.values);
      } else {
        setError(response.message || 'Failed to load movies');
      }
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to load movies');
    } finally {
      setLoading(false);
    }
  };

  const handleMovieClick = (movieId: number) => {
    navigate(`/movies/${movieId}/shows`);
  };

  return (
    <div className="movies-page">
      <header className="header">
        <h1>Movie Booking</h1>
        <div className="user-info">
          <span>Welcome, {user?.name}</span>
          <button onClick={logout} className="logout-button">
            Logout
          </button>
        </div>
      </header>

      <div className="content">
        <h2>Available Movies</h2>
        {loading && <div className="loading">Loading movies...</div>}
        {error && <div className="error-message">{error}</div>}
        {!loading && !error && movies.length === 0 && (
          <div className="empty-state">No movies available</div>
        )}
        <div className="movies-grid">
          {movies.map((movie) => (
            <div
              key={movie.id}
              className="movie-card"
              onClick={() => handleMovieClick(movie.id)}
            >
              <div className="movie-title">{movie.title}</div>
              <div className="movie-rating">{movie.rating}</div>
              <div className="movie-duration">{movie.duration_mins} minutes</div>
              <div className="movie-description">{movie.description}</div>
              <button className="select-button">Select Movie</button>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default MoviesPage;
