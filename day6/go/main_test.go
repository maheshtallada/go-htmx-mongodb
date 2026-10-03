package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	rec := httptest.NewRecorder()
	Router().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))

	if rec.Code != http.StatusOK || rec.Body.String() != "ok" {
		t.Fatalf("got %d %q; want 200 %q", rec.Code, rec.Body.String(), "ok")
	}
}

func TestGreet(t *testing.T) {
	rec := httptest.NewRecorder()
	Router().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/greet/Mahesh", nil))

	if got, want := rec.Body.String(), "Hello, Mahesh!"; got != want {
		t.Fatalf("body = %q; want %q", got, want)
	}
}

func TestRequestID(t *testing.T) {
	rec := httptest.NewRecorder()
	RequestID(Router()).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))

	if got, want := rec.Header().Get("X-Request-ID"), "test-123"; got != want {
		t.Fatalf("X-Request-ID = %q; want %q", got, want)
	}
}

func TestRecovererReturnsInternalServerError(t *testing.T) {
	rec := httptest.NewRecorder()
	Router().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/panic", nil))

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d; want %d", rec.Code, http.StatusInternalServerError)
	}
}
