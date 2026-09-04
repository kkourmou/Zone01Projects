package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestApiAsciiArtHandlerSuccess(t *testing.T) {
	rec := postJSON(t, `{"text":"Hi","banner":"standard"}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var res ApiResponse
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if res.Error != "" {
		t.Fatalf("unexpected error: %s", res.Error)
	}
	if res.Result == "" {
		t.Fatal("expected generated ASCII result")
	}
}

func TestApiAsciiArtHandlerRejectsUnsafeRequests(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		contentType string
		wantStatus  int
	}{
		{
			name:        "wrong content type",
			body:        `{"text":"Hi","banner":"standard"}`,
			contentType: "text/plain",
			wantStatus:  http.StatusUnsupportedMediaType,
		},
		{
			name:        "unknown field",
			body:        `{"text":"Hi","banner":"standard","extra":true}`,
			contentType: "application/json",
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "unsupported character",
			body:        `{"text":"¿","banner":"standard"}`,
			contentType: "application/json",
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "invalid banner",
			body:        `{"text":"Hi","banner":"../../standard"}`,
			contentType: "application/json",
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "text too long",
			body:        `{"text":"` + strings.Repeat("a", maxTextBytes+1) + `","banner":"standard"}`,
			contentType: "application/json",
			wantStatus:  http.StatusRequestEntityTooLarge,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/ascii-art", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", tt.contentType)
			rec := httptest.NewRecorder()

			ApiAsciiArtHandler(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body: %s", rec.Code, tt.wantStatus, rec.Body.String())
			}
		})
	}
}

func TestApiAsciiArtHandlerRejectsOversizedBody(t *testing.T) {
	body := `{"text":"` + strings.Repeat("a", maxRequestBytes) + `","banner":"standard"}`
	rec := postJSON(t, body)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want %d; body: %s", rec.Code, http.StatusRequestEntityTooLarge, rec.Body.String())
	}
}

func TestApiAsciiArtHandlerMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/ascii-art", nil)
	rec := httptest.NewRecorder()

	ApiAsciiArtHandler(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
	if rec.Header().Get("Allow") != http.MethodPost {
		t.Fatalf("Allow = %q, want %q", rec.Header().Get("Allow"), http.MethodPost)
	}
}

func TestHomeAndStaticHandlers(t *testing.T) {
	t.Run("home renders", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		HomeHandler(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		if !strings.Contains(rec.Body.String(), "Banner Generator") {
			t.Fatal("expected home page content")
		}
	})

	t.Run("unknown page is not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/missing", nil)
		rec := httptest.NewRecorder()

		HomeHandler(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
		}
	})

	t.Run("static directory is not listed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/static/", nil)
		rec := httptest.NewRecorder()

		StaticHandler(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
		}
	})
}

func TestSecurityHeaders(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler := SecurityHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	handler.ServeHTTP(rec, req)

	requiredHeaders := []string{
		"Content-Security-Policy",
		"Permissions-Policy",
		"Referrer-Policy",
		"X-Content-Type-Options",
		"X-Frame-Options",
	}
	for _, header := range requiredHeaders {
		if rec.Header().Get(header) == "" {
			t.Fatalf("expected %s header", header)
		}
	}
}

func postJSON(t *testing.T, body string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost, "/api/ascii-art", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	ApiAsciiArtHandler(rec, req)
	return rec
}
