HR Payroll - Simple Frontend

This is a minimal static frontend for the `hr-payroll` backend. It lists employees and lets you create a new employee.

Files:
- index.html: UI with form and list
- app.js: JS to call backend API
- styles.css: Minimal styling

Usage

1. Run the backend (default: http://localhost:8080)
   - Make sure your PostgreSQL is configured and the backend is running: `go run ./cmd` or `make run` depending on your project setup.

2. Serve the `frontend/` folder as static files. Examples:

   - Using Python 3 built-in server (from the `frontend/` directory):

     python3 -m http.server 3000

     Then open: http://localhost:3000

   - Using VS Code Live Server extension: open `frontend/index.html` and start Live Server.

CORS note

Because the backend listens on `http://localhost:8080` and the frontend will usually be served from another origin (e.g. `http://localhost:3000`), the backend must allow CORS for the frontend to call it from the browser.

Quick backend CORS enable (add to `SetupRouter` in `backend/internal/delivery/http/router.go`):

```go
import "github.com/gin-contrib/cors"

func SetupRouter(cfg RouterConfig) *gin.Engine {
  router := gin.Default()
  router.Use(cors.Default())
  // ...
}
```

Then `go mod tidy` to add the dependency and restart the backend.

If you prefer not to change the backend, you can test the frontend by using a proxy or by using the browser with relaxed CORS (not recommended).

Notes

- Do NOT include `id` in the create form — the backend generates it automatically.
- The frontend expects the API at `http://localhost:8080/api/v1`.
