// API Response Types
export interface ApiResponse<T> {
  success: boolean;
  statusCode: number;
  message?: string;
  values?: T;
  error?: FieldError[];
}

export interface FieldError {
  field: string;
  message: string;
}

// User Types
export interface User {
  id: number;
  name: string;
  email: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: User;
}

// Movie Types
export interface Movie {
  id: number;
  title: string;
  description: string;
  duration_mins: number;
  rating: string;
  created_at: string;
  updated_at: string;
}

// Show Types
export interface Theatre {
  id: number;
  name: string;
  location: string;
}

export interface Show {
  id: number;
  movie_id: number;
  theatre_id: number;
  start_time: string;
  movie?: Movie;
  theatre?: Theatre;
}

// Seat Types
export interface ShowSeat {
  id: number;
  show_id: number;
  seat_name: string;
  status: 'AVAILABLE' | 'LOCKED' | 'SOLD';
  locked_at?: string;
  user_id?: number;
}

export interface LockSeatResponse {
  message: string;
  expires_at: string;
}

// Booking Types
export interface CreateBookingRequest {
  show_id: number;
  seat_id: number;
}

export interface BookingResponse {
  booking_id: number;
  status: string;
  message: string;
}
