// Package repositories defines the API of the repository layer.
package repositories

import (
	"errors"

	"dev.beta.audi/gorepo/gopher_skeleton/internal/models"
)

// List of possible errors for this layer
var (
	ErrKnownUsersCountExceeded = errors.New("too many known users")
)

// UsersStore repository
type UsersStore interface {
	GetKnownUsers(ctx *models.Context) ([]string, error)
	AddKnownUser(ctx *models.Context, name string) error
	ResetKnownUsers(ctx *models.Context) error
}
