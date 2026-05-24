package api

import (
	"net/http"

	"mthan-go-starter/app/routes"
	"mthan-go-starter/app/services"
)

// Register mounts the /api base endpoint to the multiplexer.
func Register(mux *http.ServeMux, svc *services.ProcessService, logger *services.LoggerService) {
	mux.HandleFunc("/api", routes.LoggingMiddleware(logger, routes.RecoveryMiddleware(logger, routes.CORSMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			routes.SendError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
			return
		}
		routes.SendJSON(w, http.StatusOK, map[string]string{
			"message": "Welcome to the public client API. Access sub-routes like /api/hello.",
		})
	}))))
}
