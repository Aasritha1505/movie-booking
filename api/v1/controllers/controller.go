package controllers

import (
	"encoding/json"
	"net/http"

	"movie-booking/api/v1/helpers"
	"movie-booking/api/v1/types"
	"movie-booking/core/services"
	appcontext "movie-booking/util/context"
	"movie-booking/util/errors"
	"github.com/sirupsen/logrus"
)

// Controller handles HTTP requests
type Controller struct {
	authService    services.AuthServiceInterface
	movieService   services.MovieServiceInterface
	showService    services.ShowServiceInterface
	seatService    services.SeatServiceInterface
	bookingService services.BookingServiceInterface
}

// NewController creates a new controller instance
func NewController(
	authService services.AuthServiceInterface,
	movieService services.MovieServiceInterface,
	showService services.ShowServiceInterface,
	seatService services.SeatServiceInterface,
	bookingService services.BookingServiceInterface,
) *Controller {
	return &Controller{
		authService:    authService,
		movieService:   movieService,
		showService:    showService,
		seatService:    seatService,
		bookingService: bookingService,
	}
}

// ResponseHandler wraps handlers for consistent response formatting
func ResponseHandler(f types.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := f(w, r)
		if err != nil {
			ErrorHandler(err, w, r)
			return
		}
		writeGenericResponse(res, w, r)
	}
}

// ErrorHandler formats error responses
func ErrorHandler(err error, w http.ResponseWriter, r *http.Request) {
	response := &types.GenericAPIResponse{
		Success:    false,
		Message:   "Internal server error",
		StatusCode: http.StatusInternalServerError,
	}

	// Check for HTTP errors with status codes
	if httpErr, ok := errors.IsHTTPError(err); ok {
		response.StatusCode = httpErr.StatusCode
		response.Message = httpErr.Message
	} else {
		// Check for common error types
		switch {
		case err.Error() == "invalid credentials":
			response.StatusCode = http.StatusUnauthorized
			response.Message = "Invalid credentials"
		case err.Error() == "authorization header missing" || err.Error() == "invalid authorization header format" || err.Error() == "invalid token":
			response.StatusCode = http.StatusUnauthorized
			response.Message = "Unauthorized"
		default:
			// Log unexpected errors
			logrus.WithError(err).Error("Unexpected error in handler")
		}
	}

	writeGenericResponse(response, w, r)
}

func writeGenericResponse(res *types.GenericAPIResponse, w http.ResponseWriter, r *http.Request) {
	// Ensure CORS headers are set (backup in case middleware didn't set them)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Idempotency-Key")
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.StatusCode)
	json.NewEncoder(w).Encode(res)
}

// LoginHandler handles POST /api/v1/login
func (c *Controller) LoginHandler(w http.ResponseWriter, r *http.Request) (*types.GenericAPIResponse, error) {
	TAG := "[Login]"
	ctx := r.Context()
	logger := logrus.WithContext(ctx)

	logger.Info(TAG, "Login request received")

	// Parse and validate request
	req, err := helpers.ValidateAndParseLoginRequest(r)
	if err != nil {
		return nil, errors.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Call service layer
	result, err := c.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		logger.WithError(err).Error(TAG, "Login failed")
		return nil, err
	}

	logger.Info(TAG, "Login successful")

	return &types.GenericAPIResponse{
		Success:    true,
		StatusCode: http.StatusOK,
		Message:    "Login successful",
		Values:     result,
	}, nil
}

// GetMoviesHandler handles GET /api/v1/movies
func (c *Controller) GetMoviesHandler(w http.ResponseWriter, r *http.Request) (*types.GenericAPIResponse, error) {
	TAG := "[GetMovies]"
	ctx := r.Context()
	logger := logrus.WithContext(ctx)

	logger.Info(TAG, "Get movies request")

	movies, err := c.movieService.GetAllMovies(ctx)
	if err != nil {
		logger.WithError(err).Error(TAG, "Failed to get movies")
		return nil, errors.Wrap(err, "failed to get movies")
	}

	return &types.GenericAPIResponse{
		Success:    true,
		StatusCode: http.StatusOK,
		Message:    "Movies retrieved successfully",
		Values:     movies,
	}, nil
}

// GetShowsByMovieHandler handles GET /api/v1/movies/:id/shows
func (c *Controller) GetShowsByMovieHandler(w http.ResponseWriter, r *http.Request) (*types.GenericAPIResponse, error) {
	TAG := "[GetShowsByMovie]"
	ctx := r.Context()
	logger := logrus.WithContext(ctx)

	// Parse movie ID from path
	movieID, err := helpers.ParseMovieIDFromPath(r)
	if err != nil {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "invalid movie ID")
	}

	logger.WithField("movieID", movieID).Info(TAG, "Get shows for movie")

	shows, err := c.showService.GetShowsByMovieID(ctx, movieID)
	if err != nil {
		logger.WithError(err).Error(TAG, "Failed to get shows")
		return nil, errors.Wrap(err, "failed to get shows")
	}

	return &types.GenericAPIResponse{
		Success:    true,
		StatusCode: http.StatusOK,
		Message:    "Shows retrieved successfully",
		Values:     shows,
	}, nil
}

// GetSeatsByShowHandler handles GET /api/v1/shows/:id/seats
func (c *Controller) GetSeatsByShowHandler(w http.ResponseWriter, r *http.Request) (*types.GenericAPIResponse, error) {
	TAG := "[GetSeatsByShow]"
	ctx := r.Context()
	logger := logrus.WithContext(ctx)

	// Parse show ID from path
	showID, err := helpers.ParseShowIDFromPath(r)
	if err != nil {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "invalid show ID")
	}

	logger.WithField("showID", showID).Info(TAG, "Get seats for show")

	seats, err := c.seatService.GetSeatsByShowID(ctx, showID)
	if err != nil {
		logger.WithError(err).Error(TAG, "Failed to get seats")
		return nil, errors.Wrap(err, "failed to get seats")
	}

	return &types.GenericAPIResponse{
		Success:    true,
		StatusCode: http.StatusOK,
		Message:    "Seats retrieved successfully",
		Values:     seats,
	}, nil
}

// LockSeatHandler handles PATCH /api/v1/seats/:id/lock
func (c *Controller) LockSeatHandler(w http.ResponseWriter, r *http.Request) (*types.GenericAPIResponse, error) {
	TAG := "[LockSeat]"
	ctx := r.Context()
	logger := logrus.WithContext(ctx)

	// Get user ID from context (set by auth interceptor)
	userID, ok := appcontext.GetUserID(ctx)
	if !ok {
		return nil, errors.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	// Parse seat ID from path
	seatID, err := helpers.ParseSeatIDFromPath(r)
	if err != nil {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "invalid seat ID")
	}

	logger.WithFields(logrus.Fields{
		"seatID": seatID,
		"userID": userID,
	}).Info(TAG, "Lock seat request")

	// Call service layer
	result, err := c.seatService.LockSeat(ctx, seatID, userID)
	if err != nil {
		logger.WithError(err).Error(TAG, "Failed to lock seat")
		return nil, err
	}

	logger.Info(TAG, "Seat locked successfully")

	return &types.GenericAPIResponse{
		Success:    true,
		StatusCode: http.StatusOK,
		Message:    result.Message,
		Values:     result,
	}, nil
}

// CreateBookingHandler handles POST /api/v1/bookings
func (c *Controller) CreateBookingHandler(w http.ResponseWriter, r *http.Request) (*types.GenericAPIResponse, error) {
	TAG := "[CreateBooking]"
	ctx := r.Context()
	logger := logrus.WithContext(ctx)

	// Get user ID from context
	userID, ok := appcontext.GetUserID(ctx)
	if !ok {
		return nil, errors.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	// Parse and validate request
	input, err := helpers.ValidateAndParseBookingRequest(r, userID)
	if err != nil {
		return nil, errors.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	logger.WithFields(logrus.Fields{
		"showID": input.ShowID,
		"seatID": input.SeatID,
		"userID": userID,
	}).Info(TAG, "Create booking request")

	// Call service layer
	result, err := c.bookingService.CreateBooking(ctx, input)
	if err != nil {
		logger.WithError(err).Error(TAG, "Failed to create booking")
		return nil, err
	}

	logger.Info(TAG, "Booking created successfully")

	return &types.GenericAPIResponse{
		Success:    true,
		StatusCode: http.StatusOK,
		Message:    result.Message,
		Values:     result,
	}, nil
}
