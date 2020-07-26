// Package viewer provides utility functions for passing the current user to a context.
package viewer

import "context"

type key int

const viewerKey key = 0

// Viewer represents a user performing an action on the server.
type Viewer struct {
	UserID    string
	Username  string
	FullName  string
	SessionID string
}

// FromContext returns the current user information if present.
func FromContext(ctx context.Context) (Viewer, bool) {
	v, ok := ctx.Value(viewerKey).(Viewer)
	return v, ok
}

// NewContext updates the parent context with viewer information.
func NewContext(ctx context.Context, v Viewer) context.Context {
	return context.WithValue(ctx, viewerKey, v)
}
