import axios, { AxiosInstance } from 'axios';
import {
  ApiResponse,
  LoginRequest,
  LoginResponse,
  Movie,
  Show,
  ShowSeat,
  LockSeatResponse,
  CreateBookingRequest,
  BookingResponse,
} from '../types';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

class ApiService {
  private client: AxiosInstance;

  constructor() {
    this.client = axios.create({
      baseURL: API_BASE_URL,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // Add request interceptor to include auth token
    this.client.interceptors.request.use(
      (config) => {
        const token = localStorage.getItem('authToken');
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
      },
      (error) => {
        return Promise.reject(error);
      }
    );

    // Add response interceptor for error handling
    this.client.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response?.status === 401) {
          // Unauthorized - clear token and redirect to login
          localStorage.removeItem('authToken');
          localStorage.removeItem('user');
          window.location.href = '/login';
        }
        return Promise.reject(error);
      }
    );
  }

  // Auth endpoints
  async login(credentials: LoginRequest): Promise<ApiResponse<LoginResponse>> {
    const response = await this.client.post<ApiResponse<LoginResponse>>(
      '/api/v1/login',
      credentials
    );
    return response.data;
  }

  // Movie endpoints
  async getMovies(): Promise<ApiResponse<Movie[]>> {
    const response = await this.client.get<ApiResponse<Movie[]>>('/api/v1/movies');
    return response.data;
  }

  async getShowsByMovie(movieId: number): Promise<ApiResponse<Show[]>> {
    const response = await this.client.get<ApiResponse<Show[]>>(
      `/api/v1/movies/${movieId}/shows`
    );
    return response.data;
  }

  // Seat endpoints
  async getSeatsByShow(showId: number): Promise<ApiResponse<ShowSeat[]>> {
    const response = await this.client.get<ApiResponse<ShowSeat[]>>(
      `/api/v1/shows/${showId}/seats`
    );
    return response.data;
  }

  async lockSeat(seatId: number): Promise<ApiResponse<LockSeatResponse>> {
    const response = await this.client.patch<ApiResponse<LockSeatResponse>>(
      `/api/v1/seats/${seatId}/lock`
    );
    return response.data;
  }

  // Booking endpoints
  async createBooking(
    booking: CreateBookingRequest,
    idempotencyKey?: string
  ): Promise<ApiResponse<BookingResponse>> {
    const headers = idempotencyKey ? { 'Idempotency-Key': idempotencyKey } : {};
    const response = await this.client.post<ApiResponse<BookingResponse>>(
      '/api/v1/bookings',
      booking,
      { headers }
    );
    return response.data;
  }
}

export const apiService = new ApiService();
