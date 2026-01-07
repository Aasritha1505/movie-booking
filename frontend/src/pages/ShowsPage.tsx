import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { apiService } from '../services/api';
import { Show } from '../types';
import './ShowsPage.css';

const ShowsPage: React.FC = () => {
  const { movieId } = useParams<{ movieId: string }>();
  const navigate = useNavigate();
  const [shows, setShows] = useState<Show[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    if (movieId) {
      loadShows(parseInt(movieId));
    }
  }, [movieId]);

  const loadShows = async (id: number) => {
    try {
      setLoading(true);
      const response = await apiService.getShowsByMovie(id);
      if (response.success && response.values) {
        setShows(response.values);
      } else {
        setError(response.message || 'Failed to load shows');
      }
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to load shows');
    } finally {
      setLoading(false);
    }
  };

  const handleShowClick = (showId: number) => {
    navigate(`/shows/${showId}/seats`);
  };

  const formatDateTime = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleString('en-US', {
      weekday: 'short',
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  return (
    <div className="shows-page">
      <header className="header">
        <button onClick={() => navigate('/movies')} className="back-button">
          ← Back to Movies
        </button>
        <h1>Select Showtime</h1>
      </header>

      <div className="content">
        {loading && <div className="loading">Loading shows...</div>}
        {error && <div className="error-message">{error}</div>}
        {!loading && !error && shows.length === 0 && (
          <div className="empty-state">No shows available for this movie</div>
        )}
        <div className="shows-list">
          {shows.map((show) => (
            <div
              key={show.id}
              className="show-card"
              onClick={() => handleShowClick(show.id)}
            >
              <div className="show-info">
                <div className="theatre-name">{show.theatre?.name || 'Theatre'}</div>
                <div className="show-time">{formatDateTime(show.start_time)}</div>
                <div className="theatre-location">{show.theatre?.location || ''}</div>
              </div>
              <button className="select-button">Select Seats</button>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default ShowsPage;
