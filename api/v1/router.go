package v1

import (
	"net/http"

	"movie-booking/api/v1/controllers"
	"movie-booking/api/v1/interceptors"
	"movie-booking/config"
	"github.com/gorilla/mux"
)

// Route represents a route definition
type Route struct {
	Path         string
	RequestMethod string
	Handler      http.HandlerFunc
	SkipAuth     bool
	DoNotLog     bool
}

// AddRoutesToRouter registers all routes with the router
func AddRoutesToRouter(router *mux.Router, ctrl *controllers.Controller) error {
	// Create route definitions
	var routes = []Route{
		{
			Path:         "/api/v1/login",
			RequestMethod: http.MethodPost,
			Handler:      controllers.ResponseHandler(ctrl.LoginHandler),
			SkipAuth:     true,
			DoNotLog:     false,
		},
		{
			Path:         "/api/v1/movies",
			RequestMethod: http.MethodGet,
			Handler:      controllers.ResponseHandler(ctrl.GetMoviesHandler),
			SkipAuth:     true,
			DoNotLog:     false,
		},
		{
			Path:         "/api/v1/movies/{id}/shows",
			RequestMethod: http.MethodGet,
			Handler:      controllers.ResponseHandler(ctrl.GetShowsByMovieHandler),
			SkipAuth:     true,
			DoNotLog:     false,
		},
		{
			Path:         "/api/v1/shows/{id}/seats",
			RequestMethod: http.MethodGet,
			Handler:      controllers.ResponseHandler(ctrl.GetSeatsByShowHandler),
			SkipAuth:     true,
			DoNotLog:     false,
		},
		{
			Path:         "/api/v1/seats/{id}/lock",
			RequestMethod: http.MethodPatch,
			Handler:      controllers.ResponseHandler(ctrl.LockSeatHandler),
			SkipAuth:     false, // Requires auth
			DoNotLog:     false,
		},
		{
			Path:         "/api/v1/bookings",
			RequestMethod: http.MethodPost,
			Handler:      controllers.ResponseHandler(ctrl.CreateBookingHandler),
			SkipAuth:     false, // Requires auth
			DoNotLog:     false,
		},
	}

	// Register each route
	for _, route := range routes {
		// Build interceptor chain
		interceptorChain := []interceptors.Interceptor{
			interceptors.PanicRecoveryInterceptor(controllers.ErrorHandler),
			interceptors.LoggingInterceptor(route.DoNotLog),
		}

		// Add auth interceptor if required
		if !route.SkipAuth {
			interceptorChain = append(interceptorChain,
				interceptors.AuthInterceptor(controllers.ErrorHandler))
		}

		// Apply interceptors
		handler := interceptors.Intercept(route.Handler, interceptorChain...)

		// Add timeout handler
		timeout := config.GetHandlerTimeout()
		timeoutHandler := http.TimeoutHandler(handler, timeout, `{"success":false,"message":"Request timeout"}`)

		// Register route
		router.Handle(route.Path, timeoutHandler).Methods(route.RequestMethod)
	}
	
	return nil
}
