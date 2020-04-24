package models

import (
	"context"

	"google.golang.org/grpc/metadata"
)

const (
	// RequestIDKey const
	RequestIDKey = "x-request-id"
)

// Context is the general purpose context
type Context struct {
	context.Context
	requestID *string // required
}

// NewContext creates a new typed context
func NewContext(c context.Context) *Context {
	ctx := Context{
		Context: c,
	}
	return &ctx
}

// RequestID retruns the request id if existsing
func (c *Context) RequestID() (string, bool) {
	if c.requestID != nil {
		return *c.requestID, true
	}
	if c.Context == nil {
		return "", false
	}
	if md, ok := metadata.FromIncomingContext(c.Context); ok {
		if v := md.Get(RequestIDKey); len(v) > 0 {
			if len(v[0]) == 0 {
				return "", false
			}
			id := v[0] // copy
			c.requestID = &id
			return id, true
		}
	}
	return "", false
}
