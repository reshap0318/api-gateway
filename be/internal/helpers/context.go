package helpers

import "context"

// ContextKey represents a type for context keys.
type ContextKey string

const (
	// KeyUserID is the context key for storing the authenticated user's ID.
	KeyUserID ContextKey = "user_id"
)

// GetCallerID extracts the user ID from the context.
// Returns 0 if the user ID is not present in the context.
func GetCallerID(ctx context.Context) uint {
	uid, ok := ctx.Value(KeyUserID).(uint)
	if !ok {
		return 0
	}
	return uid
}
