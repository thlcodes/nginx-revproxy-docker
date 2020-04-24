package repositories_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"

	"dev.beta.audi/gorepo/gopher_skeleton/internal/models"
	"google.golang.org/grpc"

	errorslib "dev.beta.audi/gorepo/lib-go-common/errors"
	storagepb "dev.beta.audi/gorepo/lib_proto_models/golib/storage"

	"dev.beta.audi/gorepo/gopher_skeleton/internal/repositories"

	"github.com/stretchr/testify/require"
)

var testCtx = context.TODO()

func Test_StorageALRepo_AddGetKnownUsers(t *testing.T) {
	ctx := models.NewContext(testCtx)
	client, _ := createMockClient()
	repo := repositories.NewPersistentStore(cfg, client)
	want := []string{"Peter", "Jacky"}
	for _, user := range want {
		require.NoError(t, repo.AddKnownUser(ctx, user))
	}
	got, err := repo.GetKnownUsers(ctx)
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func Test_StorageALRepo_ExceedKnownUsers(t *testing.T) {
	ctx := models.NewContext(testCtx)
	client, _ := createMockClient()
	repo := repositories.NewPersistentStore(cfg, client)
	for i := int64(0); i <= cfg.Greeter.MaxKnownUsers; i++ {
		err := repo.AddKnownUser(ctx, fmt.Sprintf("User%d", i))
		if i < cfg.Greeter.MaxKnownUsers {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
			require.IsType(t, &errorslib.BadRequestError{}, err)
			require.Equal(t, repositories.ErrKnownUsersCountExceeded.Error(), err.(errorslib.Error).Unwrap().Error())
		}
	}
	got, err := repo.GetKnownUsers(ctx)
	require.NoError(t, err)
	require.Equal(t, cfg.Greeter.MaxKnownUsers, int64(len(got)))
}

func Test_StorageALRepo_Reset(t *testing.T) {
	ctx := models.NewContext(testCtx)
	client, _ := createMockClient()
	repo := repositories.NewPersistentStore(cfg, client)
	send := []string{"Peter", "Jacky", "Peter"}
	want := []string{"Peter", "Jacky"}
	for _, user := range send {
		require.NoError(t, repo.AddKnownUser(ctx, user))
	}
	got, err := repo.GetKnownUsers(ctx)
	require.NoError(t, err)
	require.Equal(t, want, got)

	// Reset
	err = repo.ResetKnownUsers(ctx)
	require.NoError(t, err)
	got, err = repo.GetKnownUsers(ctx)
	require.NoError(t, err)
	require.Empty(t, got)
}

func createMockClient() (*storagepb.MockStorageServiceClient, error) {
	var err error
	client := &storagepb.MockStorageServiceClient{}
	tmpList := []string{}
	client.On(
		"Get",
		mock.Anything,
		mock.AnythingOfType(fmt.Sprintf("%T", &storagepb.StorageRequest{})),
	).Return(func(_ context.Context, req *storagepb.StorageRequest, _ ...grpc.CallOption) *storagepb.KeyValuePair {
		raw, _ := json.Marshal(tmpList)
		return &storagepb.KeyValuePair{Key: "", Value: &storagepb.Document{Value: string(raw)}}
	}, err)
	client.On(
		"Set",
		mock.Anything,
		mock.AnythingOfType(fmt.Sprintf("%T", &storagepb.StorageRequest{})),
	).Return(nil, func(_ context.Context, req *storagepb.StorageRequest, _ ...grpc.CallOption) error {
		list := []string{}
		data := (req.Body.(*storagepb.StorageRequest_KvpWithTtl)).KvpWithTtl.Kvp.Value.Value
		if err := json.Unmarshal([]byte(data), &list); err != nil {
			return err
		}
		tmpList = list
		return nil
	})
	return client, err
}
