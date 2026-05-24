package mods

import (
	"fmt"
	"net/http"

	"mthan-go-starter/app/routes"
	"mthan-go-starter/app/services"
)

// Register mounts the default fallback landing page at "/" when the client is disabled.
func Register(mux *http.ServeMux, svc *services.ProcessService, logger *services.LoggerService) {
	mux.HandleFunc("/", routes.LoggingMiddleware(logger, routes.RecoveryMiddleware(logger, routes.CORSMiddleware(func(w http.ResponseWriter, r *http.Request) {
		// Only handle direct root requests. Return 404 for unhandled sub-paths.
		if r.URL.Path != "/" {
			routes.SendError(w, http.StatusNotFound, "NOT_FOUND", "Path not found")
			return
		}

		if r.Method != http.MethodGet {
			routes.SendError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `
<!DOCTYPE html>
<html>
<head>
    <title>Mthan Go Starter</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; background: #0b0f19; color: #f8fafc; padding: 40px; display: flex; align-items: center; justify-content: center; min-height: 85vh; margin: 0; }
        .card { background: rgba(17, 24, 30, 0.6); backdrop-filter: blur(16px); padding: 40px; border-radius: 20px; max-width: 500px; text-align: center; box-shadow: 0 20px 40px rgba(0,0,0,0.3); border: 1px solid rgba(255,255,255,0.08); }
        h1 { background: linear-gradient(135deg, #38bdf8 0%, #818cf8 100%); -webkit-background-clip: text; -webkit-text-fill-color: transparent; margin-top: 0; font-size: 2.2rem; font-weight: 700; letter-spacing: -0.5px; }
        p { color: #94a3b8; line-height: 1.6; font-size: 1rem; }
        .badge { display: inline-block; background: rgba(56, 189, 248, 0.1); color: #38bdf8; padding: 6px 16px; border-radius: 9999px; font-size: 0.85em; font-weight: 600; margin-bottom: 24px; border: 1px solid rgba(56, 189, 248, 0.2); }
    </style>
</head>
<body>
    <div class="card">
        <div class="badge">Go Backend Active</div>
        <h1>Mthan Go Starter</h1>
        <p>The backend service is running successfully. The React client integration is currently disabled (CLIENT=false).</p>
    </div>
</body>
</html>
		`)
	}))))
}
