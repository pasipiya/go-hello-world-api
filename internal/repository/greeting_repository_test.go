package repository

import "testing"

func TestInMemoryGreetingRepository_GreetingPrefix(t *testing.T) {
	repo := NewInMemoryGreetingRepository()

	if got, want := repo.GreetingPrefix(), "Hello"; got != want {
		t.Errorf("GreetingPrefix() = %q, want %q", got, want)
	}
}
