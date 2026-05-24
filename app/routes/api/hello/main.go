package hello

import (
	"net/http"

	"mthan-go-starter/app/routes"
	"mthan-go-starter/app/services"
)

// Register mounts the GET /api/hello route on the multiplexer.
func Register(mux *http.ServeMux, svc *services.ProcessService, logger *services.LoggerService) {
	mux.HandleFunc("/api/hello", routes.LoggingMiddleware(logger, routes.RecoveryMiddleware(logger, routes.CORSMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			routes.SendError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
			return
		}
		
		name := r.URL.Query().Get("name")
		
		result, err := svc.ProcessAPI(name)
		if err != nil {
			routes.SendError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
			return
		}
		
		routes.SendJSON(w, http.StatusOK, map[string]string{
			"message": result,
		})
	}))))
}
