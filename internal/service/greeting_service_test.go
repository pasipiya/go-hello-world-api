package service

import (
	"errors"
	"testing"
)

// stubRepo is a test double for repository.GreetingRepository.
type stubRepo struct{ prefix string }

func (s stubRepo) GreetingPrefix() string { return s.prefix }

func TestGreetingService_Greet(t *testing.T) {
	svc := NewGreetingService(stubRepo{prefix: "Hello"})

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr error
	}{
		// Valid: first half of the alphabet.
		{name: "uppercase A", input: "Alice", want: "Hello Alice"},
		{name: "lowercase a", input: "alice", want: "Hello alice"},
		{name: "upper boundary M", input: "Mike", want: "Hello Mike"},
		{name: "lower boundary m", input: "mike", want: "Hello mike"},
		{name: "single letter", input: "b", want: "Hello b"},
		{name: "surrounding whitespace trimmed", input: "  Alice  ", want: "Hello Alice"},
		{name: "rest of name may be anything", input: "J0hn-Ü", want: "Hello J0hn-Ü"},

		// Invalid: second half of the alphabet.
		{name: "boundary N", input: "Nina", wantErr: ErrInvalidInput},
		{name: "boundary n", input: "nina", wantErr: ErrInvalidInput},
		{name: "uppercase Z", input: "Zane", wantErr: ErrInvalidInput},
		{name: "lowercase z", input: "zane", wantErr: ErrInvalidInput},

		// Invalid: edge cases.
		{name: "empty", input: "", wantErr: ErrInvalidInput},
		{name: "whitespace only", input: "   \t\n", wantErr: ErrInvalidInput},
		{name: "starts with digit", input: "1alice", wantErr: ErrInvalidInput},
		{name: "starts with symbol", input: "@alice", wantErr: ErrInvalidInput},
		{name: "non-ASCII letter", input: "Élodie", wantErr: ErrInvalidInput},
		{name: "emoji", input: "😀", wantErr: ErrInvalidInput},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := svc.Greet(tc.input)

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Greet(%q) error = %v, want %v", tc.input, err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("Greet(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestGreetingService_UsesRepositoryPrefix(t *testing.T) {
	svc := NewGreetingService(stubRepo{prefix: "Hi"})

	got, err := svc.Greet("Alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "Hi Alice"; got != want {
		t.Errorf("Greet() = %q, want %q", got, want)
	}
}
