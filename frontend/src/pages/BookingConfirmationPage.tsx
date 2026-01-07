import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import './BookingConfirmationPage.css';

const BookingConfirmationPage: React.FC = () => {
  const { bookingId } = useParams<{ bookingId: string }>();
  const navigate = useNavigate();
  const { user } = useAuth();
  const [countdown, setCountdown] = useState(5);

  useEffect(() => {
    const timer = setInterval(() => {
      setCountdown((prev) => {
        if (prev <= 1) {
          clearInterval(timer);
          navigate('/movies');
          return 0;
        }
        return prev - 1;
      });
    }, 1000);

    return () => clearInterval(timer);
  }, [navigate]);

  return (
    <div className="confirmation-page">
      <div className="confirmation-card">
        <div className="success-icon">✓</div>
        <h1>Booking Confirmed!</h1>
        <p className="booking-id">Booking ID: #{bookingId}</p>
        <p className="confirmation-message">
          Your ticket has been sent to your email.
        </p>
        <div className="user-info">
          <p>User: {user?.name}</p>
          <p>Email: {user?.email}</p>
        </div>
        <div className="actions">
          <button onClick={() => navigate('/movies')} className="primary-button">
            Book Another Movie
          </button>
          <p className="redirect-message">
            Redirecting to movies in {countdown} seconds...
          </p>
        </div>
      </div>
    </div>
  );
};

export default BookingConfirmationPage;
