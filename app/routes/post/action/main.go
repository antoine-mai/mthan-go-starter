package action

import (
	"encoding/json"
	"net/http"

	"mthan-go-starter/app/routes"
	"mthan-go-starter/app/services"
)

// Register mounts the POST /post/action route on the multiplexer.
func Register(mux *http.ServeMux, svc *services.ProcessService, logger *services.LoggerService) {
	mux.HandleFunc("/post/action", routes.LoggingMiddleware(logger, routes.RecoveryMiddleware(logger, routes.CORSMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			routes.SendError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
			return
		}
		
		var payload map[string]string
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			routes.SendError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
			return
		}
		
		result, err := svc.ProcessPost(payload)
		if err != nil {
			routes.SendError(w, http.StatusBadRequest, "PROCESS_ERROR", err.Error())
			return
		}
		
		routes.SendJSON(w, http.StatusOK, map[string]string{
			"message": result,
		})
	}))))
}
