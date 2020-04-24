package repositories

import (
	"encoding/json"
	"strings"

	storagepb "dev.beta.audi/gorepo/lib_proto_models/golib/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errorslib "dev.beta.audi/gorepo/lib-go-common/errors"
	"dev.beta.audi/gorepo/lib-go-common/tracing"

	"dev.beta.audi/gorepo/gopher_skeleton/internal/generated/config"
	"dev.beta.audi/gorepo/gopher_skeleton/internal/models"
)

const collection = "greeter"
const knownUsersKey = "knownUsers"

// PersistentStore provides a user repository
type PersistentStore struct {
	config *config.Config
	client storagepb.StorageServiceClient
}

var _ UsersStore = (*PersistentStore)(nil)

/*
NewPersistentStore constructs a greeter repository.

	Parameters:
		- config: The service configuration.
		- client: An object that implements the API of the storage service client.

	Returns:
		- The initialized repository.
*/
func NewPersistentStore(config *config.Config, client storagepb.StorageServiceClient) UsersStore {
	return &PersistentStore{
		config: config,
		client: client,
	}
}

// GetKnownUsers returns a list of known users
func (r *PersistentStore) GetKnownUsers(ctx *models.Context) (list []string, err error) {
	span, sctx := tracing.StartSpan(ctx.Context, tracing.FuncName(tracing.FuncNameWithParentType))
	defer span.Finish()
	ctx.Context = sctx
	if r.client == nil {
		return nil, tracing.LogErrorSpan(span, errorslib.NewInternalError("storage client is nil", nil))
	}
	req := &storagepb.StorageRequest{
		Context:    storagepb.StorageContext__Platform,
		Collection: collection,
		Body:       &storagepb.StorageRequest_Selector{Selector: &storagepb.Selector{Selector: knownUsersKey}},
	}
	item, err := r.client.Get(ctx.Context, req)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return []string{}, nil
		}
		return nil, tracing.LogErrorSpan(span, errorslib.NewClientError("could not retreive known users", err))
	}
	if item == nil || item.Value == nil || item.Value.Value == "" {
		return nil, tracing.LogErrorSpan(span, errorslib.NewClientError("got nil pr empty document from", nil))
	}
	if err := json.Unmarshal([]byte(item.Value.Value), &list); err != nil {
		return nil, tracing.LogErrorSpan(span, errorslib.NewTypeError("could not unmarshal known users list", err))
	}
	return
}

// AddKnownUser adds a user to the 'known' lis
func (r *PersistentStore) AddKnownUser(ctx *models.Context, name string) error {
	var span tracing.Span
	span, ctx.Context = tracing.StartSpan(ctx.Context, tracing.FuncName(tracing.FuncNameWithParentType))
	defer span.Finish()
	if r.client == nil {
		return tracing.LogErrorSpan(span, errorslib.NewInternalError("storage client is nil", nil))
	}
	knownUsers, err := r.GetKnownUsers(ctx)
	if err != nil {
		return err
	}
	if int64(len(knownUsers)) >= r.config.Greeter.MaxKnownUsers {
		return tracing.LogErrorSpan(span, errorslib.WithCause(errorslib.NewBadRequestError("could not add user"), ErrKnownUsersCountExceeded))
	}
	found := false
	lowerName := strings.ToLower(name)
	for _, user := range knownUsers {
		if strings.ToLower(user) == lowerName {
			found = true
			break
		}
	}
	if found { // do nothing, name is already kwown
		return nil
	}
	knownUsers = append(knownUsers, name)
	raw, _ := json.Marshal(knownUsers)
	req := &storagepb.StorageRequest{
		Context:    storagepb.StorageContext__Platform,
		Collection: collection,
		Body: &storagepb.StorageRequest_KvpWithTtl{KvpWithTtl: &storagepb.KeyValuePairWithTTL{
			Kvp: &storagepb.KeyValuePair{
				Key:   knownUsersKey,
				Value: &storagepb.Document{Value: string(raw)},
			},
		}},
	}
	_, err = r.client.Set(ctx.Context, req)
	if err != nil {
		return tracing.LogErrorSpan(span, errorslib.NewClientError("could not update known users list", err))
	}
	return nil
}

// ResetKnownUsers resets the known users
func (r *PersistentStore) ResetKnownUsers(ctx *models.Context) error {
	span, sctx := tracing.StartSpan(ctx.Context, tracing.FuncName(tracing.FuncNameWithParentType))
	defer span.Finish()
	ctx.Context = sctx
	if r.client == nil {
		return tracing.LogErrorSpan(span, errorslib.NewInternalError("storage client is nil", nil))
	}
	req := &storagepb.StorageRequest{
		Context:    storagepb.StorageContext__Platform,
		Collection: collection,
		Body: &storagepb.StorageRequest_KvpWithTtl{KvpWithTtl: &storagepb.KeyValuePairWithTTL{
			Kvp: &storagepb.KeyValuePair{
				Key:   knownUsersKey,
				Value: &storagepb.Document{Value: "[]"},
			},
		}},
	}
	_, err := r.client.Set(ctx.Context, req)
	if err != nil {
		return tracing.LogErrorSpan(span, errorslib.NewClientError("could not reset known users list", err))
	}
	return nil
}
