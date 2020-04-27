package grpc

import (
	"context"

	errorslib "dev.beta.audi/gorepo/lib-go-common/errors"
	commonpb "dev.beta.audi/gorepo/lib_proto_models/golib/common"

	emptypb "github.com/golang/protobuf/ptypes/empty"

	"dev.beta.audi/gorepo/gopher-user-ecomy/internal/delivery"
	"dev.beta.audi/gorepo/gopher-user-ecomy/internal/generated/config"
)

// EcomyServiceServer implements the EcomyService
type EcomyServiceServer struct {
	config *config.Config
}

var _ delivery.EcomyService = (*EcomyServiceServer)(nil)

/*
NewEcomyServiceServer constructs a status service server.

	Parameters:
		- config:  The service configuration.

	Returns:
		- The initialized greeter server.
*/
func NewEcomyServiceServer(config *config.Config) delivery.EcomyService {
	return &EcomyServiceServer{
		config: config,
	}
}

func (*EcomyServiceServer) GetStatus(context.Context, *emptypb.Empty) (*commonpb.StatusResponse, error) {
	return &commonpb.StatusResponse{Status: "OK"}, mapError(nil)
}

func mapError(err error) error {
	if err == nil {
		return nil
	}
	switch err.(type) {
	case errorslib.Error:
		return err
	}
	return errorslib.NewUnknownError(err)
}
