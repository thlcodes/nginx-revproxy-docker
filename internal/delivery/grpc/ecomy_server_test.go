package grpc_test

import (
	"context"
	"testing"

	commonpb "dev.beta.audi/gorepo/lib_proto_models/golib/common"
	"github.com/golang/protobuf/ptypes/empty"

	"github.com/stretchr/testify/require"

	grpcdelivery "dev.beta.audi/gorepo/gopher-user-ecomy/internal/delivery/grpc"
	"dev.beta.audi/gorepo/gopher-user-ecomy/internal/generated/config"
)

var (
	cfg = &config.Envs.Test
)

func Test_GetStatus(t *testing.T) {
	server := grpcdelivery.NewEcomyServiceServer(cfg)
	require.NotNil(t, server)

	want := &commonpb.StatusResponse{Status: "OK"}
	got, err := server.GetStatus(context.Background(), &empty.Empty{})
	require.NoError(t, err)
	require.Equal(t, want, got)
}
