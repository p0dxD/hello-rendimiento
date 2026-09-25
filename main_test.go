package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func get(t *testing.T, path string) (int, string) {
	t.Helper()
	rec := httptest.NewRecorder()
	handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	return rec.Code, rec.Body.String()
}

func TestHealthz(t *testing.T) {
	if code, body := get(t, "/healthz"); code != 200 || strings.TrimSpace(body) != "ok" {
		t.Fatalf("healthz = %d %q", code, body)
	}
}

func TestHomeUsesGreeting(t *testing.T) {
	t.Setenv("GREETING", "Hola")
	if code, body := get(t, "/"); code != 200 || !strings.Contains(body, "Hola") {
		t.Fatalf("home = %d", code)
	}
}

func TestUnknownPath(t *testing.T) {
	if code, _ := get(t, "/nope"); code != 404 {
		t.Fatalf("unknown path = %d", code)
	}
}
