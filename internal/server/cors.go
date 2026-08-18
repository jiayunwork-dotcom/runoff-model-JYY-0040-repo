package server

import "net/http"

// CORSMiddleware adds CORS headers for the frontend.
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// NewWithCORS creates a server with CORS support.
func NewWithCORS(addr string) *Server {
	s := New(addr)
	// Wrap the mux with CORS.
	s.mux = http.NewServeMux()
	inner := New(addr)
	s.mux.Handle("/", CORSMiddleware(inner.mux))
	return s
}
