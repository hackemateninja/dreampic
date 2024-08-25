package types

// Define a custom type for the context key
type ContextKey string

// Create a specific key for the user context
const UserContextKey ContextKey = "user"

type AuthenticatedUser struct {
	Email    string
	Avatar   string
	LoggedIn bool
}
