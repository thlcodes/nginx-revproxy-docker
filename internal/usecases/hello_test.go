package usecases_test

import (
	"context"
	"fmt"
	"testing"

	"dev.beta.audi/gorepo/gopher_skeleton/internal/generated/mocks"

	errorslib "dev.beta.audi/gorepo/lib-go-common/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"dev.beta.audi/gorepo/gopher_skeleton/internal/usecases"

	"github.com/stretchr/testify/mock"

	"dev.beta.audi/gorepo/gopher_skeleton/internal/models"

	"dev.beta.audi/gorepo/gopher_skeleton/internal/generated/config"
)

var (
	cfg = &config.Envs.Test
)

// Test_SayHello
func Test_SayHello(t *testing.T) {
	salut := cfg.Greeter.Hello + " "
	defName := cfg.Greeter.DefaultName
	name := "Peter"
	type args struct {
		name string
	}
	type mocked struct {
		knownUsers []string
		getErr     error
		addErr     error
	}
	type wants struct {
		msg string
		err error
	}
	tests := []struct {
		name   string
		args   args
		mocked mocked
		wants  wants
	}{
		{"w/o name", args{}, mocked{}, wants{salut + defName, nil}},
		{"w/ name", args{name}, mocked{}, wants{salut + name, nil}},
		{"w/ known name, assureing repo.AddKnownUser is not called", args{name}, mocked{[]string{name, "peter"}, nil, nil}, wants{salut + name + " again", nil}},
		{"inappropriate name", args{cfg.Greeter.InappropriateNames[0]}, mocked{}, wants{"", &errorslib.InvalidArgumentError{}}},
		{"repo.GetKnownUsers failed", args{name}, mocked{nil, errorslib.NewClientError("", nil), nil}, wants{"", &errorslib.ClientError{}}},
		{"repo.AddKnownUser failed", args{name}, mocked{nil, nil, errorslib.NewClientError("", nil)}, wants{"", &errorslib.ClientError{}}},
		{"repo.AddKnownUser failed with unknown error", args{name}, mocked{nil, nil, fmt.Errorf("woops")}, wants{"", &errorslib.UnknownError{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			known_users := tt.mocked.knownUsers
			repo := &mocks.UsersRepo{}
			repo.On("GetKnownUsers", mock.Anything).Return(func(_ *models.Context) []string { return known_users }, tt.mocked.getErr)
			repo.On("AddKnownUser", mock.Anything, mock.AnythingOfType("string")).Maybe().Return(func(_ *models.Context, name string) error {
				known_users = append(known_users, name)
				return tt.mocked.addErr
			})
			greeter := usecases.NewGreeter(cfg, repo)
			require.NotNil(t, greeter)
			ctx := models.NewContext(context.TODO())
			got, err := greeter.SayHello(ctx, tt.args.name)
			if tt.wants.err != nil {
				require.Error(t, err)
				//assert.IsType(t, tt.wants.err, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, got.Message, tt.wants.msg)
			repo.AssertExpectations(t)
		})
	}
}

func Test_FallbackName(t *testing.T) {
	repo := &mocks.UsersRepo{}
	greeter := usecases.NewGreeter(cfg, repo)
	ctx := models.NewContext(context.TODO())
	got, err := greeter.GetFallbackName(ctx)
	require.NoError(t, err)
	assert.Equal(t, cfg.Greeter.DefaultName, got)

	name := "Peter"
	err = greeter.SetFallbackName(ctx, name)
	require.NoError(t, err)

	got, err = greeter.GetFallbackName(ctx)
	require.NoError(t, err)
	assert.Equal(t, name, got)
}

func Test_Reset(t *testing.T) {
	repo := &mocks.UsersRepo{}
	repo.On("ResetKnownUsers", mock.Anything).Return(nil)
	greeter := usecases.NewGreeter(cfg, repo)
	ctx := models.NewContext(context.TODO())
	_ = greeter.SetFallbackName(ctx, "Peter")
	require.NoError(t, greeter.Reset(ctx))
	name, _ := greeter.GetFallbackName(ctx)
	assert.Equal(t, cfg.Greeter.DefaultName, name)
}
