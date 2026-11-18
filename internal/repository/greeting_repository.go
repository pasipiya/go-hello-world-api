// Package repository is the data layer that supplies greeting data.
package repository

// GreetingRepository provides the data needed to build a greeting.
type GreetingRepository interface {
	// GreetingPrefix returns the word used to greet a user, e.g. "Hello".
	GreetingPrefix() string
}

// InMemoryGreetingRepository returns a constant prefix; swappable for a DB/config store.
type InMemoryGreetingRepository struct {
	prefix string
}

// NewInMemoryGreetingRepository returns a repository that greets with "Hello".
func NewInMemoryGreetingRepository() *InMemoryGreetingRepository {
	return &InMemoryGreetingRepository{prefix: "Hello"}
}

// GreetingPrefix implements GreetingRepository.
func (r *InMemoryGreetingRepository) GreetingPrefix() string {
	return r.prefix
}
