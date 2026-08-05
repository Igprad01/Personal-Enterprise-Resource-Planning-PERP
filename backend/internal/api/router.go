package api

import (
	"log"
	"net/http"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", s.handleHealth)

	mux.HandleFunc("GET /api/categories", s.listCategories)
	mux.HandleFunc("POST /api/categories", s.createCategory)
	mux.HandleFunc("PUT /api/categories/{id}", s.updateCategory)
	mux.HandleFunc("DELETE /api/categories/{id}", s.deleteCategory)

	mux.HandleFunc("GET /api/activities", s.listActivities)
	mux.HandleFunc("GET /api/activities/stats", s.activityStats)
	mux.HandleFunc("POST /api/activities", s.createActivity)
	mux.HandleFunc("GET /api/activities/{id}", s.getActivity)
	mux.HandleFunc("PUT /api/activities/{id}", s.updateActivity)
	mux.HandleFunc("DELETE /api/activities/{id}", s.deleteActivity)

	return corsMiddleware(loggingMiddleware(mux))
}
