SHELL := /bin/bash

.PHONY: help backend-run frontend-serve dev dev-stop open status

TMPDIR := .tmp
# Use absolute paths so commands running from subdirectories can write logs/pids reliably
TMPDIR_ABS := $(CURDIR)/$(TMPDIR)
BACKEND_LOG := $(TMPDIR_ABS)/backend.log
FRONTEND_LOG := $(TMPDIR_ABS)/frontend.log

help: ## Show this help
	@echo "Usage: make <target>"
	@grep -E '^[a-zA-Z0-9._-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

backend-run: ## Run backend in foreground (use this for direct dev debugging)
	@echo "==> Running backend (foreground)"
	@cd backend && make run

frontend-run: ## Serve frontend on port 3000 (foreground)
	@echo "==> Serving frontend on http://localhost:3000 (foreground)"
	@cd frontend && python3 -m http.server 3000

dev: ## Start backend and frontend in background (dev). Logs -> $(TMPDIR)/
	@mkdir -p $(TMPDIR_ABS)
	@echo "Starting backend (logs -> $(BACKEND_LOG))..."
	@cd backend && nohup make run > $(BACKEND_LOG) 2>&1 &
	@sleep 0.5
	@echo "Starting frontend (logs -> $(FRONTEND_LOG))..."
	@cd frontend && nohup python3 -m http.server 3000 > $(FRONTEND_LOG) 2>&1 &
	@echo "Started. Backend on port 8080, Frontend on port 3000."
	@echo "Open frontend: http://localhost:3000"

dev-stop: ## Stop background dev servers by killing processes on ports 8080 & 3000
	@echo "==> Stopping dev servers by port"
	@BACKEND_PID=$$(lsof -t -i:8080) || true; \
	if [ -n "$$BACKEND_PID" ]; then \
		echo "  -> Killing backend process on port 8080 (PID: $$BACKEND_PID)"; \
		kill -9 $$BACKEND_PID; \
	else \
		echo "  -> No process found on port 8080"; \
	fi
	@FRONTEND_PID=$$(lsof -t -i:3000) || true; \
	if [ -n "$$FRONTEND_PID" ]; then \
		echo "  -> Killing frontend process on port 3000 (PID: $$FRONTEND_PID)"; \
		kill -9 $$FRONTEND_PID; \
	else \
		echo "  -> No process found on port 3000"; \
	fi
	@echo "==> Done."

status: ## Show running PIDs (if any)
	@printf "Backend PID: %s\n" "$$(cat $(BACKEND_PID) 2>/dev/null || echo -)"
	@printf "Frontend PID: %s\n" "$$(cat $(FRONTEND_PID) 2>/dev/null || echo -)"

open: ## Open frontend in default browser (macOS)
	@open http://localhost:3000

open-docs:
	@open http://localhost:8080/swagger/index.html

docs: ## Generate swagger docs (requires swag installed)
	@echo "==> Generating Swagger docs"
	swag init -g ./backend/cmd/main.go -o ./backend/docs