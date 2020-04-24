package grpc

import (
	"context"
	"log"

	"google.golang.org/grpc/metadata"

	"github.com/opentracing/opentracing-go"

	errorslib "dev.beta.audi/gorepo/lib-go-common/errors"
	"dev.beta.audi/gorepo/lib-go-common/tracing"
	greeterpb "dev.beta.audi/gorepo/lib_proto_models/golib/greeter"

	"dev.beta.audi/gorepo/gopher_skeleton/internal/models"

	"github.com/golang/protobuf/ptypes/empty"

	"dev.beta.audi/gorepo/gopher_skeleton/internal/delivery"
	"dev.beta.audi/gorepo/gopher_skeleton/internal/generated/config"
	"dev.beta.audi/gorepo/gopher_skeleton/internal/usecases"
)

// GreeterServer provides a gRPC delivery.
type GreeterServer struct {
	config  *config.Config
	greeter usecases.GreeterUsecases
}

var _ delivery.GreeterServer = (*GreeterServer)(nil)

/*
NewGreeterServer constructs a greeter server.

	Parameters:
		- config:  The service configuration.
		- greeter: An object that implements the use cases API.

	Returns:
		- The initialized greeter server.
*/
func NewGreeterServer(config *config.Config, greeter usecases.GreeterUsecases) delivery.GreeterServer {
	return &GreeterServer{
		config:  config,
		greeter: greeter,
	}
}

// Hello will greet someone. Logic details see usecases.
func (d *GreeterServer) Hello(inctx context.Context, req *greeterpb.HelloRequest) (*greeterpb.Message, error) {
	md, _ := metadata.FromIncomingContext(inctx)
	log.Printf("parent: %#v context-metadata: %v", opentracing.SpanFromContext(inctx), md)
	span, ctx := tracing.StartSpan(inctx, tracing.FuncName(tracing.FuncNameWithParentType))
	defer span.Finish()
	if err := req.Validate(); err != nil {
		return nil, errorslib.NewInvalidArgumentError(err.Error())
	}
	msg, err := d.greeter.SayHello(models.NewContext(ctx), req.GetName())
	if err != nil {
		return nil, mapError(err)
	}
	return msg, nil
}

// SetFallbackName will set the name to be greeted if no name is given in a hello request
func (d *GreeterServer) SetFallbackName(inctx context.Context, in *greeterpb.SetFallbackNameRequest) (*empty.Empty, error) {
	span, ctx := tracing.StartSpan(inctx, tracing.FuncName(tracing.FuncNameWithParentType))
	defer span.Finish()
	if err := in.Validate(); err != nil {
		return nil, errorslib.NewInvalidArgumentError(err.Error())
	}
	if err := d.greeter.SetFallbackName(models.NewContext(ctx), in.Name); err != nil {
		return nil, mapError(err)
	}
	return &empty.Empty{}, nil
}

// GetFallbackName will set the name to be greeted if no name is given in a hello request
func (d *GreeterServer) GetFallbackName(inctx context.Context, _ *empty.Empty) (*greeterpb.Name, error) {
	span, ctx := tracing.StartSpan(inctx, tracing.FuncName(tracing.FuncNameWithParentType))
	defer span.Finish()
	name, err := d.greeter.GetFallbackName(models.NewContext(ctx))
	if err != nil {
		return nil, mapError(err)
	}
	return &greeterpb.Name{Name: name}, nil
}

// Reset ill reset the name to be greeted if no name is given in a hello request to the services' default
func (d *GreeterServer) Reset(inctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	span, ctx := tracing.StartSpan(inctx, tracing.FuncName(tracing.FuncNameWithParentType))
	defer span.Finish()
	if err := d.greeter.Reset(models.NewContext(ctx)); err != nil {
		return nil, mapError(err)
	}
	return &empty.Empty{}, nil
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
