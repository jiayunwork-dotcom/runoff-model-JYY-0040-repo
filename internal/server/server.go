// Package server provides an HTTP API for the runoff model, allowing the web
// frontend to submit rainfall data and receive computed hydrographs.
package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"runoff-model/internal/hydro"
	"runoff-model/internal/source"
	"runoff-model/internal/basin"
)

// Server is the HTTP API server.
type Server struct {
	mux  *http.ServeMux
	addr string
}

// New creates a server listening on addr.
func New(addr string) *Server {
	s := &Server{mux: http.NewServeMux(), addr: addr}
	s.mux.HandleFunc("/api/simulate", s.handleSimulate)
	s.mux.HandleFunc("/api/uh", s.handleUH)
	s.mux.HandleFunc("/api/route", s.handleRoute)
	s.mux.HandleFunc("/api/health", s.handleHealth)
	// Serve static frontend files.
	s.mux.Handle("/", http.FileServer(http.Dir("frontend")))
	return s
}

// Handler returns the HTTP handler (for testing).
func (s *Server) Handler() http.Handler { return s.mux }

// ListenAndServe starts the server.
func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(s.addr, s.mux)
}

// SimulateRequest is the input for /api/simulate.
type SimulateRequest struct {
	Basin    BasinParams `json:"basin"`
	Rainfall []float64   `json:"rainfall"`
	ET       []float64   `json:"et"`
}

// BasinParams matches basin.Params fields.
type BasinParams struct {
	Area float64 `json:"area"`
	Wm   float64 `json:"wm"`
	B    float64 `json:"b"`
	C    float64 `json:"c"`
	N    int     `json:"n"`
	K    float64 `json:"k"`
}

// SimulateResponse is the output of /api/simulate.
type SimulateResponse struct {
	EffectiveRain []float64 `json:"effective_rain"`
	UH            []float64 `json:"uh"`
	DirectRunoff  []float64 `json:"direct_runoff"`
}

func (s *Server) handleSimulate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req SimulateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Build records from rainfall/ET arrays.
	records := make([]basin.Record, len(req.Rainfall))
	for i := range records {
		records[i].Rain = req.Rainfall[i]
		if i < len(req.ET) {
			records[i].ET = req.ET[i]
		}
	}

	effRain := source.Production(records, req.Basin.Wm, req.Basin.B, req.Basin.C)

	n := req.Basin.N
	if n <= 0 {
		n = 3
	}
	k := req.Basin.K
	if k <= 0 {
		k = 2.0
	}
	uh := hydro.NashUH(n, k, 1.0)
	directRunoff := hydro.Convolve(effRain, uh)

	resp := SimulateResponse{
		EffectiveRain: effRain,
		UH:            uh,
		DirectRunoff:  directRunoff,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// UHRequest is the input for /api/uh.
type UHRequest struct {
	N  int     `json:"n"`
	K  float64 `json:"k"`
	Dt float64 `json:"dt"`
}

func (s *Server) handleUH(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req UHRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.Dt <= 0 {
		req.Dt = 1.0
	}
	uh := hydro.NashUH(req.N, req.K, req.Dt)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"uh": uh})
}

// RouteRequest is the input for /api/route.
type RouteRequest struct {
	Inflow []float64 `json:"inflow"`
	K      float64   `json:"k"`
	X      float64   `json:"x"`
	Dt     float64   `json:"dt"`
}

func (s *Server) handleRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req RouteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	outflow := hydro.Muskingum(req.Inflow, req.K, req.X, req.Dt)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"outflow": outflow})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"ok"}`)
}
