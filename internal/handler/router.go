package handler

import (
	"log/slog"
	"net/http"
	"time"
)

// NewRouter собирает маршруты net/http (требование ТЗ).
func NewRouter(dh *DepartmentHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /departments/", dh.Create)
	mux.HandleFunc("GET /departments/{id}", dh.Get)
	mux.HandleFunc("PATCH /departments/{id}", dh.Patch)
	mux.HandleFunc("DELETE /departments/{id}", dh.Delete)
	mux.HandleFunc("POST /departments/{id}/employees/", dh.CreateEmployee)

	return withLogging(mux)
}

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)
		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.status,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}
