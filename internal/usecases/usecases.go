// Package usecases defines the API of the use cases layer.
package usecases

import (
	"errors"

	"dev.beta.audi/gorepo/gopher_skeleton/internal/models"
)

// List of possible errors for this layer
var (
	ErrInappropriateName = errors.New("the given name is inappropriate")
)

// GreeterUsecases service
type GreeterUsecases interface {
	// Says hell to given name
	SayHello(ctx *models.Context, name string) (*models.Message, error)
	SetFallbackName(ctx *models.Context, name string) error
	GetFallbackName(ctx *models.Context) (string, error)
	Reset(ctx *models.Context) error
}
