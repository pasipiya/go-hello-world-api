package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pasipiya/go-hello-world-api/internal/repository"
	"github.com/pasipiya/go-hello-world-api/internal/service"
)

// newTestServer wires the real layers for end-to-end endpoint tests.
func newTestServer() http.Handler {
	repo := repository.NewInMemoryGreetingRepository()
	svc := service.NewGreetingService(repo)
	mux := http.NewServeMux()
	NewHelloController(svc).RegisterRoutes(mux)
	return mux
}

func TestHelloWorld(t *testing.T) {
	tests := []struct {
		name       string
		target     string
		wantStatus int
		wantBody   map[string]string
	}{
		{"valid uppercase", "/hello-world?name=Alice", http.StatusOK, map[string]string{"message": "Hello Alice"}},
		{"valid lowercase", "/hello-world?name=alice", http.StatusOK, map[string]string{"message": "Hello alice"}},
		{"valid boundary M", "/hello-world?name=Mike", http.StatusOK, map[string]string{"message": "Hello Mike"}},
		{"url-encoded spaces trimmed", "/hello-world?name=%20Bob%20", http.StatusOK, map[string]string{"message": "Hello Bob"}},
		{"multiple values uses first", "/hello-world?name=Alice&name=Zane", http.StatusOK, map[string]string{"message": "Hello Alice"}},

		{"invalid second half", "/hello-world?name=Zane", http.StatusBadRequest, map[string]string{"error": "Invalid Input"}},
		{"invalid boundary N", "/hello-world?name=nina", http.StatusBadRequest, map[string]string{"error": "Invalid Input"}},
		{"missing name", "/hello-world", http.StatusBadRequest, map[string]string{"error": "Invalid Input"}},
		{"empty name", "/hello-world?name=", http.StatusBadRequest, map[string]string{"error": "Invalid Input"}},
		{"whitespace name", "/hello-world?name=%20%20", http.StatusBadRequest, map[string]string{"error": "Invalid Input"}},
		{"digit first", "/hello-world?name=9lives", http.StatusBadRequest, map[string]string{"error": "Invalid Input"}},
		{"wrong param name", "/hello-world?username=Alice", http.StatusBadRequest, map[string]string{"error": "Invalid Input"}},
	}

	handler := newTestServer()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.target, nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tc.wantStatus)
			}
			if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
				t.Errorf("Content-Type = %q, want application/json", ct)
			}

			var got map[string]string
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Fatalf("invalid JSON body %q: %v", rec.Body.String(), err)
			}
			if len(got) != len(tc.wantBody) {
				t.Fatalf("body = %v, want %v", got, tc.wantBody)
			}
			for k, v := range tc.wantBody {
				if got[k] != v {
					t.Errorf("body[%q] = %q, want %q", k, got[k], v)
				}
			}
		})
	}
}

func TestHelloWorld_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/hello-world?name=Alice", nil)
	rec := httptest.NewRecorder()

	newTestServer().ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

// failingService simulates an unexpected service failure.
type failingService struct{}

func (failingService) Greet(string) (string, error) { return "", errors.New("boom") }

func TestHelloWorld_UnexpectedError(t *testing.T) {
	mux := http.NewServeMux()
	NewHelloController(failingService{}).RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/hello-world?name=Alice", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}
