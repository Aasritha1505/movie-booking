import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { apiService } from '../services/api';
import { ShowSeat } from '../types';
import './SeatsPage.css';

const SeatsPage: React.FC = () => {
  const { showId } = useParams<{ showId: string }>();
  const navigate = useNavigate();
  const [seats, setSeats] = useState<ShowSeat[]>([]);
  const [selectedSeat, setSelectedSeat] = useState<number | null>(null);
  const [lockedSeat, setLockedSeat] = useState<number | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [lockExpiresAt, setLockExpiresAt] = useState<Date | null>(null);

  useEffect(() => {
    if (showId) {
      loadSeats(parseInt(showId));
    }
  }, [showId]);

  useEffect(() => {
    // Refresh seats every 5 seconds to show updated status
    const interval = setInterval(() => {
      if (showId) {
        loadSeats(parseInt(showId));
      }
    }, 5000);
    return () => clearInterval(interval);
  }, [showId]);

  const loadSeats = async (id: number) => {
    try {
      const response = await apiService.getSeatsByShow(id);
      if (response.success && response.values) {
        setSeats(response.values);
      } else {
        setError(response.message || 'Failed to load seats');
      }
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to load seats');
    } finally {
      setLoading(false);
    }
  };

  const handleSeatClick = async (seat: ShowSeat) => {
    if (seat.status === 'SOLD') {
      setError('This seat is already sold');
      return;
    }

    if (seat.status === 'LOCKED' && seat.locked_at) {
      const lockTime = new Date(seat.locked_at);
      const now = new Date();
      const lockDuration = 10 * 60 * 1000; // 10 minutes
      if (now.getTime() - lockTime.getTime() < lockDuration) {
        setError('This seat is currently locked by another user');
        return;
      }
    }

    try {
      setError('');
      const response = await apiService.lockSeat(seat.id);
      if (response.success && response.values) {
        setLockedSeat(seat.id);
        setSelectedSeat(seat.id);
        setLockExpiresAt(new Date(response.values.expires_at));
        // Reload seats to show updated status
        if (showId) {
          loadSeats(parseInt(showId));
        }
      } else {
        setError(response.message || 'Failed to lock seat');
      }
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to lock seat');
    }
  };

  const handleBook = async () => {
    if (!selectedSeat || !showId) return;

    try {
      setError('');
      const idempotencyKey = `booking-${selectedSeat}-${Date.now()}`;
      const response = await apiService.createBooking(
        {
          show_id: parseInt(showId),
          seat_id: selectedSeat,
        },
        idempotencyKey
      );

      if (response.success && response.values) {
        navigate(`/booking-confirmation/${response.values.booking_id}`);
      } else {
        setError(response.message || 'Failed to create booking');
      }
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to create booking');
    }
  };

  const getSeatStatusClass = (seat: ShowSeat): string => {
    if (seat.status === 'SOLD') return 'seat-sold';
    if (seat.status === 'LOCKED') {
      if (seat.id === lockedSeat) return 'seat-locked-by-me';
      return 'seat-locked';
    }
    if (seat.id === selectedSeat) return 'seat-selected';
    return 'seat-available';
  };

  const getSeatStatusLabel = (seat: ShowSeat): string => {
    if (seat.status === 'SOLD') return 'Sold';
    if (seat.status === 'LOCKED') {
      if (seat.id === lockedSeat) return 'Your Lock';
      return 'Locked';
    }
    return 'Available';
  };

  // Organize seats into rows (assuming A1-A10, B1-B10, etc.)
  const organizeSeats = () => {
    const rows: { [key: string]: ShowSeat[] } = {};
    seats.forEach((seat) => {
      const row = seat.seat_name.charAt(0);
      if (!rows[row]) {
        rows[row] = [];
      }
      rows[row].push(seat);
    });

    // Sort seats within each row
    Object.keys(rows).forEach((row) => {
      rows[row].sort((a, b) => {
        const numA = parseInt(a.seat_name.substring(1));
        const numB = parseInt(b.seat_name.substring(1));
        return numA - numB;
      });
    });

    return rows;
  };

  const seatRows = organizeSeats();

  return (
    <div className="seats-page">
      <header className="header">
        <button onClick={() => navigate(-1)} className="back-button">
          ← Back
        </button>
        <h1>Select Your Seat</h1>
      </header>

      <div className="content">
        {loading && <div className="loading">Loading seats...</div>}
        {error && <div className="error-message">{error}</div>}

        {!loading && (
          <>
            <div className="screen-indicator">SCREEN</div>

            <div className="seats-container">
              {Object.keys(seatRows).sort().map((row) => (
                <div key={row} className="seat-row">
                  <div className="row-label">{row}</div>
                  <div className="seats-in-row">
                    {seatRows[row].map((seat) => (
                      <button
                        key={seat.id}
                        className={`seat ${getSeatStatusClass(seat)}`}
                        onClick={() => handleSeatClick(seat)}
                        disabled={seat.status === 'SOLD' || (seat.status === 'LOCKED' && seat.id !== lockedSeat)}
                        title={getSeatStatusLabel(seat)}
                      >
                        {seat.seat_name.substring(1)}
                      </button>
                    ))}
                  </div>
                </div>
              ))}
            </div>

            <div className="legend">
              <div className="legend-item">
                <div className="seat-legend seat-available"></div>
                <span>Available</span>
              </div>
              <div className="legend-item">
                <div className="seat-legend seat-selected"></div>
                <span>Selected</span>
              </div>
              <div className="legend-item">
                <div className="seat-legend seat-locked-by-me"></div>
                <span>Your Lock</span>
              </div>
              <div className="legend-item">
                <div className="seat-legend seat-locked"></div>
                <span>Locked</span>
              </div>
              <div className="legend-item">
                <div className="seat-legend seat-sold"></div>
                <span>Sold</span>
              </div>
            </div>

            {selectedSeat && lockedSeat && (
              <div className="booking-section">
                <div className="selected-seat-info">
                  <h3>Selected Seat: {seats.find((s) => s.id === selectedSeat)?.seat_name}</h3>
                  {lockExpiresAt && (
                    <p>Lock expires at: {lockExpiresAt.toLocaleTimeString()}</p>
                  )}
                </div>
                <button onClick={handleBook} className="book-button">
                  Confirm Booking
                </button>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
};

export default SeatsPage;
