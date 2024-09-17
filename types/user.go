package types

import "github.com/google/uuid"

// Define a custom type for the context key
type ContextKey string

// Create a specific key for the user context
const UserContextKey ContextKey = "user"

type AuthenticatedUser struct {
	ID       uuid.UUID
	Email    string
	LoggedIn bool
	Account
}
