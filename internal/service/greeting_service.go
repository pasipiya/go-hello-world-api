// Package service holds the business rules, independent of HTTP.
package service

import (
	"errors"
	"strings"

	"github.com/pasipiya/go-hello-world-api/internal/repository"
)

// ErrInvalidInput means the name is empty or does not start with A–M.
var ErrInvalidInput = errors.New("invalid input")

// GreetingService builds greetings for names.
type GreetingService interface {
	Greet(name string) (string, error)
}

type greetingService struct {
	repo repository.GreetingRepository
}

// NewGreetingService returns a GreetingService that uses repo for its data.
func NewGreetingService(repo repository.GreetingRepository) GreetingService {
	return &greetingService{repo: repo}
}

// Greet returns "<prefix> <name>" for trimmed names starting with A–M, else ErrInvalidInput.
func (s *greetingService) Greet(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" || !isFirstHalfOfAlphabet(name[0]) {
		return "", ErrInvalidInput
	}
	return s.repo.GreetingPrefix() + " " + name, nil
}

// isFirstHalfOfAlphabet reports whether b is in A–M/a–m; multi-byte UTF-8 never matches.
func isFirstHalfOfAlphabet(b byte) bool {
	return (b >= 'a' && b <= 'm') || (b >= 'A' && b <= 'M')
}
