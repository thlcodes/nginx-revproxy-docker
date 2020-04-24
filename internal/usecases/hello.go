package usecases

import (
	"strings"
	"sync"

	errorslib "dev.beta.audi/gorepo/lib-go-common/errors"
	"dev.beta.audi/gorepo/lib-go-common/tracing"

	"dev.beta.audi/gorepo/gopher_skeleton/internal/models"

	"dev.beta.audi/gorepo/gopher_skeleton/internal/generated/config"
	"dev.beta.audi/gorepo/gopher_skeleton/internal/repositories"
)

// Greeter implements the API of the use cases layer.
type Greeter struct {
	config       *config.Config
	usersStore   repositories.UsersStore
	lock         sync.RWMutex
	fallbackName string
}

var _ GreeterUsecases = (*Greeter)(nil)

/*
NewGreeter constructs a greeter usecases object.

	Parameters:
		- config:  The service configuration.
		- repo:    An object that implements the repositories API.

	Returns:
		- The initialized use cases object.
*/
func NewGreeter(config *config.Config, repo repositories.UsersStore) GreeterUsecases {
	return &Greeter{
		config:       config,
		usersStore:   repo,
		fallbackName: config.Greeter.DefaultName,
	}
}

/*
SayHello says hello to the given name.
If the name is known, " again" is appended.

	Parameters:
		- ctx: The context object.
		- name: The name to greet.

	Returns:
		- message: The greeting message.
		- error:
			Nil if the greeting succeeds.
			InvalidArgument error if the name is within a list of inappropriate names.
			If the operation fails for other reasons, refer to the errors returned by the repositories layer (GetKnownUsers, AddKnownUser).
*/
func (u *Greeter) SayHello(ctx *models.Context, name string) (*models.Message, error) {
	span, sctx := tracing.StartSpan(ctx.Context, tracing.FuncName(tracing.FuncNameWithParentType))
	defer span.Finish()
	ctx.Context = sctx
	if name != "" {
		lowerName := strings.ToLower(name)
		for _, n := range u.config.Greeter.InappropriateNames {
			if strings.ToLower(n) == lowerName {
				err := errorslib.NewInvalidArgumentError("could not say hello to '" + name + "'")
				err.WithCause(ErrInappropriateName)
				return nil, err
			}
		}
	} else {
		u.lock.RLock()
		name = u.fallbackName
		u.lock.RUnlock()
	}
	knownUsers, err := u.usersStore.GetKnownUsers(ctx)
	if err != nil {
		return nil, mapError(err)
	}
	known := false
	for _, user := range knownUsers {
		if user == name {
			known = true
			break
		}
	}
	suffix := ""
	if known {
		suffix = " again"
	} else {
		if err := u.usersStore.AddKnownUser(ctx, name); err != nil {
			return nil, mapError(err)
		}
	}
	return &models.Message{Message: u.config.Greeter.Hello + " " + name + suffix}, nil
}

// SetFallbackName set the default name to greet. The operation always succeeds.
func (u *Greeter) SetFallbackName(ctx *models.Context, name string) error {
	span, sctx := tracing.StartSpan(ctx.Context, tracing.FuncName(tracing.FuncNameWithParentType))
	defer span.Finish()
	ctx.Context = sctx
	u.lock.Lock()
	defer u.lock.Unlock()
	u.fallbackName = name
	return nil
}

// GetFallbackName returns the default name to greet. The operation always succeeds.
func (u *Greeter) GetFallbackName(ctx *models.Context) (string, error) {
	span, sctx := tracing.StartSpan(ctx.Context, tracing.FuncName(tracing.FuncNameWithParentType))
	defer span.Finish()
	ctx.Context = sctx
	u.lock.RLock()
	defer u.lock.RUnlock()
	return u.fallbackName, nil
}

// Reset clears the list of known users. If the operation fails, refer to the errors returned by the repositories layer (ResetKnownUsers).
func (u *Greeter) Reset(ctx *models.Context) error {
	span, sctx := tracing.StartSpan(ctx.Context, tracing.FuncName(tracing.FuncNameWithParentType))
	defer span.Finish()
	ctx.Context = sctx
	u.fallbackName = u.config.Greeter.DefaultName
	return mapError(u.usersStore.ResetKnownUsers(ctx))
}

func mapError(err error) error {
	if err == nil {
		return nil
	}
	switch err := err.(type) {
	case errorslib.Error:
		return err
	}
	return errorslib.NewUnknownError(err)
}
