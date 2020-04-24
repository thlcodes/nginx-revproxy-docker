package repositories_test

import (
	"context"
	"fmt"
	"testing"

	"dev.beta.audi/gorepo/gopher_skeleton/internal/models"
	"dev.beta.audi/gorepo/gopher_skeleton/internal/repositories"

	errorslib "dev.beta.audi/gorepo/lib-go-common/errors"

	"github.com/stretchr/testify/require"
)

func Test_MemoryRepo_AddGetKnownUsers(t *testing.T) {
	ctx := models.NewContext(context.TODO())
	repo := repositories.NewInMemoryStore(cfg)
	want := []string{"Peter", "Jacky"}
	for _, user := range want {
		require.NoError(t, repo.AddKnownUser(ctx, user))
	}
	got, err := repo.GetKnownUsers(ctx)
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func Test_MemoryRepo_ExceedKnownUsers(t *testing.T) {
	ctx := models.NewContext(context.TODO())
	repo := repositories.NewInMemoryStore(cfg)
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

func Test_MemoryRepo_Reset(t *testing.T) {
	ctx := models.NewContext(context.TODO())
	repo := repositories.NewInMemoryStore(cfg)
	want := []string{"Peter", "Jacky"}
	for _, user := range want {
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
