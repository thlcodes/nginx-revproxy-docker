package grpc_test

import (
	"context"
	"testing"

	"dev.beta.audi/gorepo/gopher_skeleton/internal/generated/mocks"

	errorslib "dev.beta.audi/gorepo/lib-go-common/errors"

	"github.com/stretchr/testify/mock"

	"dev.beta.audi/gorepo/gopher_skeleton/internal/models"
	greeterpb "dev.beta.audi/gorepo/lib_proto_models/golib/greeter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	grpcdelivery "dev.beta.audi/gorepo/gopher_skeleton/internal/delivery/grpc"
	"dev.beta.audi/gorepo/gopher_skeleton/internal/generated/config"
)

var (
	cfg = &config.Envs.Test
)

// Test_Hello
func Test_Hello(t *testing.T) {
	salut := cfg.Greeter.Hello + " "
	again := " again"
	defName := cfg.Greeter.DefaultName
	name := "Peter"
	type args struct {
		name string
	}
	type mocked struct {
		known bool
		err   error
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
		{"w/ known name", args{name}, mocked{true, nil}, wants{salut + name, nil}},
		{"inappropriate name", args{"nasty guy"}, mocked{false, &errorslib.InvalidArgumentError{}}, wants{"", &errorslib.InvalidArgumentError{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suffix := ""
			if tt.mocked.known {
				suffix = again
			}
			greeter := &mocks.GreeterUsecases{}
			greeter.On("SayHello", mock.Anything, mock.Anything).Return(&greeterpb.Message{Message: salut + suffix}, tt.mocked.err)

			server := grpcdelivery.NewGreeterServer(cfg, greeter)
			require.NotNil(t, server)
			client := grpcdelivery.GreeterClient{Server: server}

			ctx := context.TODO()
			got, err := client.Hello(ctx, &greeterpb.HelloRequest{Name: tt.args.name})
			if tt.wants.err != nil {
				require.Error(t, err)
				assert.IsType(t, tt.wants.err, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
			greeter.AssertExpectations(t)
		})
	}
}

func Test_FallbackName(t *testing.T) {
	name := cfg.Greeter.DefaultName
	greeter := &mocks.GreeterUsecases{}
	greeter.On("SetFallbackName", mock.Anything, mock.Anything).Return(nil)
	greeter.On("GetFallbackName", mock.Anything, mock.Anything).Return(func(*models.Context) string {
		return name
	}, nil)

	server := grpcdelivery.NewGreeterServer(cfg, greeter)
	require.NotNil(t, server)
	client := grpcdelivery.GreeterClient{Server: server}
	ctx := context.TODO()

	got, err := client.GetFallbackName(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, cfg.Greeter.DefaultName, got.Name)

	name = "Peter"
	_, err = client.SetFallbackName(ctx, &greeterpb.SetFallbackNameRequest{Name: name})
	require.NoError(t, err)

	got, err = client.GetFallbackName(ctx, nil)
	require.NoError(t, err)
	assert.Equal(t, name, got.Name)
}

func Test_Reset(t *testing.T) {
	greeter := &mocks.GreeterUsecases{}
	greeter.On("Reset", mock.Anything).Return(nil)
	server := grpcdelivery.NewGreeterServer(cfg, greeter)
	require.NotNil(t, server)
	client := grpcdelivery.GreeterClient{Server: server}
	ctx := context.TODO()
	_, err := client.Reset(ctx, nil)
	require.NoError(t, err)
}
