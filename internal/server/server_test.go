package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	s := New(":0")
	req := httptest.NewRequest("GET", "/api/health", nil)
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestSimulate(t *testing.T) {
	s := New(":0")
	body := SimulateRequest{
		Basin:    BasinParams{Area: 100, Wm: 120, B: 0.3, C: 0.15, N: 3, K: 2},
		Rainfall: []float64{10, 20, 30, 15, 5},
		ET:       []float64{2, 2, 2, 2, 2},
	}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/simulate", bytes.NewReader(data))
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp SimulateResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if len(resp.DirectRunoff) == 0 {
		t.Fatal("expected non-empty direct runoff")
	}
}

func TestUH(t *testing.T) {
	s := New(":0")
	body := UHRequest{N: 3, K: 2, Dt: 1}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/uh", bytes.NewReader(data))
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestRoute(t *testing.T) {
	s := New(":0")
	body := RouteRequest{Inflow: []float64{0, 10, 20, 15, 5}, K: 2, X: 0.2, Dt: 1}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/route", bytes.NewReader(data))
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestSimulateMethodNotAllowed(t *testing.T) {
	s := New(":0")
	req := httptest.NewRequest("GET", "/api/simulate", nil)
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", w.Code)
	}
}
