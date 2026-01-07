package interceptors

import (
	"fmt"
	"net/http"
	"time"

	"movie-booking/api/v1/helpers"
	appcontext "movie-booking/util/context"
	"github.com/sirupsen/logrus"
)

// Interceptor is a function that wraps an HTTP handler
type Interceptor func(http.HandlerFunc) http.HandlerFunc

// Intercept chains multiple interceptors
func Intercept(handler http.HandlerFunc, interceptors ...Interceptor) http.HandlerFunc {
	for i := len(interceptors) - 1; i >= 0; i-- {
		handler = interceptors[i](handler)
	}
	return handler
}

// LoggingInterceptor logs request details
func LoggingInterceptor(doNotLog bool) Interceptor {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if doNotLog {
				next(w, r)
				return
			}

			start := time.Now()
			logger := logrus.WithFields(logrus.Fields{
				"method": r.Method,
				"path":   r.URL.Path,
				"ip":     r.RemoteAddr,
			})

			logger.Info("Request started")

			// Create response writer wrapper to capture status code
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next(rw, r)

			duration := time.Since(start)
			logger.WithFields(logrus.Fields{
				"status_code": rw.statusCode,
				"duration_ms": duration.Milliseconds(),
			}).Info("Request completed")
		}
	}
}

// AuthInterceptor validates JWT token and sets user ID in context
func AuthInterceptor(errorHandler func(error, http.ResponseWriter, *http.Request)) Interceptor {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Extract token
			token, err := helpers.ExtractBearerToken(r)
			if err != nil {
				errorHandler(err, w, r)
				return
			}

			// Validate token
			claims, err := helpers.ValidateJWT(token)
			if err != nil {
				errorHandler(err, w, r)
				return
			}

			// Extract user ID from claims
			userIDFloat, ok := claims["user_id"].(float64)
			if !ok {
				errorHandler(fmt.Errorf("invalid user_id in token"), w, r)
				return
			}
			userID := uint(userIDFloat)

			// Set user ID in context
			ctx := appcontext.SetUserID(r.Context(), userID)
			r = r.WithContext(ctx)

			next(w, r)
		}
	}
}

// PanicRecoveryInterceptor recovers from panics
func PanicRecoveryInterceptor(errorHandler func(error, http.ResponseWriter, *http.Request)) Interceptor {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					err := fmt.Errorf("panic recovered: %v", rec)
					logrus.WithError(err).Error("Panic in handler")
					errorHandler(err, w, r)
				}
			}()
			next(w, r)
		}
	}
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
