package repositories

import (
	"sync"

	errorslib "dev.beta.audi/gorepo/lib-go-common/errors"

	"dev.beta.audi/gorepo/gopher_skeleton/internal/generated/config"
	"dev.beta.audi/gorepo/gopher_skeleton/internal/models"
)

// InMemoryStore provides a user repository
type InMemoryStore struct {
	config     *config.Config
	knownUsers []string
	lock       sync.RWMutex
}

var _ UsersStore = (*InMemoryStore)(nil)

// NewInMemoryStore ctr
func NewInMemoryStore(config *config.Config) UsersStore {
	return &InMemoryStore{
		config:     config,
		knownUsers: []string{},
	}
}

// GetKnownUsers returns a list of known users
func (r *InMemoryStore) GetKnownUsers(ctx *models.Context) (list []string, err error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	// copy curent knonw users and return that to prevent concurrency issues
	list = make([]string, len(r.knownUsers))
	copy(list, r.knownUsers)
	return list, nil
}

// AddKnownUser adds a user to the 'known' lis
func (r *InMemoryStore) AddKnownUser(ctx *models.Context, name string) error {
	if int64(len(r.knownUsers)) >= r.config.Greeter.MaxKnownUsers {
		return errorslib.WithCause(errorslib.NewBadRequestError("could not add user"), ErrKnownUsersCountExceeded)
	}
	r.lock.Lock()
	defer r.lock.Unlock()
	r.knownUsers = append(r.knownUsers, name)
	return nil
}

// ResetKnownUsers resets the known users
func (r *InMemoryStore) ResetKnownUsers(ctx *models.Context) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.knownUsers = []string{}
	return nil
}
