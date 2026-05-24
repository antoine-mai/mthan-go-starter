package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"mthan-go-starter/app/services"
)

// Response represents a standard JSON response wrapper.
type Response struct {
	Success bool         `json:"success"`
	Data    interface{}  `json:"data,omitempty"`
	Error   *ErrorDetail `json:"error,omitempty"`
}

// ErrorDetail represents standardized JSON error format.
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// SendJSON sends a structured successful JSON response.
func SendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    data,
	})
}

// SendError sends a structured error JSON response.
func SendError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{
		Success: false,
		Error: &ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}

// LoggingMiddleware logs HTTP request details using the injected LoggerService.
func LoggingMiddleware(logger *services.LoggerService, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		wrappedWriter := &statusResponseWriter{ResponseWriter: w, status: http.StatusOK}
		
		next(wrappedWriter, r)
		
		logger.Info("HTTP Request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrappedWriter.status,
			"latency", time.Since(start).String(),
			"ip", r.RemoteAddr,
		)
	}
}

// RecoveryMiddleware recovers from panics in handler execution using the injected LoggerService.
func RecoveryMiddleware(logger *services.LoggerService, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("HTTP Panic Recovered", "error", err)
				SendError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "An unexpected error occurred")
			}
		}()
		next(w, r)
	}
}

// CORSMiddleware handles CORS options preflight and sets global cross-origin headers.
func CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

// BasicAuth protects routes with basic authentication.
func BasicAuth(username, password string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != username || pass != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Admin Panel"`)
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Unauthorized")
			return
		}
		next(w, r)
	}
}

// statusResponseWriter is a helper to capture status codes sent to the client.
type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
