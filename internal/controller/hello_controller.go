// Package controller is the HTTP layer: parses requests and maps service results to responses.
package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/pasipiya/go-hello-world-api/internal/service"
)

// MessageResponse is the body returned on success.
type MessageResponse struct {
	Message string `json:"message"`
}

// ErrorResponse is the body returned on failure.
type ErrorResponse struct {
	Error string `json:"error"`
}

// HelloController handles the /hello-world endpoint.
type HelloController struct {
	svc service.GreetingService
}

// NewHelloController returns a HelloController backed by svc.
func NewHelloController(svc service.GreetingService) *HelloController {
	return &HelloController{svc: svc}
}

// RegisterRoutes binds GET /hello-world; other methods get 405 automatically.
func (c *HelloController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /hello-world", c.HelloWorld)
}

// HelloWorld handles GET /hello-world?name=<name>.
func (c *HelloController) HelloWorld(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	msg, err := c.svc.Greet(name)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid Input"})
			return
		}
		log.Printf("hello-world: unexpected error: %v", err)
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Internal Server Error"})
		return
	}

	writeJSON(w, http.StatusOK, MessageResponse{Message: msg})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("writeJSON: failed to encode response: %v", err)
	}
}
