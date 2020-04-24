package grpc

import (
	"context"

	"dev.beta.audi/gorepo/lib-go-common/tracing"

	"google.golang.org/grpc/metadata"

	"github.com/golang/protobuf/ptypes/empty"

	"dev.beta.audi/gorepo/gopher_skeleton/internal/delivery"

	greeterpb "dev.beta.audi/gorepo/lib_proto_models/golib/greeter"
	"google.golang.org/grpc"
)

// GreeterClient wrapper
type GreeterClient struct {
	Server delivery.GreeterServer
}

var _ delivery.GreeterClient = (*GreeterClient)(nil)

func (c *GreeterClient) Hello(inctx context.Context, req *greeterpb.HelloRequest, _ ...grpc.CallOption) (*greeterpb.Message, error) {

	span, ctx := tracing.StartSpan(inctx, tracing.FuncName(tracing.FuncNameWithParentType))
	defer span.Finish()
	ctx = outgoingToIncomingCtx(ctx)
	return c.Server.Hello(ctx, req)
}

func (c *GreeterClient) SetFallbackName(inctx context.Context, in *greeterpb.SetFallbackNameRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	span, ctx := tracing.StartSpan(inctx, tracing.FuncName(tracing.FuncNameWithParentType))
	defer span.Finish()
	ctx = outgoingToIncomingCtx(ctx)
	return c.Server.SetFallbackName(ctx, in)
}

func (c *GreeterClient) GetFallbackName(inctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*greeterpb.Name, error) {
	span, ctx := tracing.StartSpan(inctx, tracing.FuncName(tracing.FuncNameWithParentType))
	defer span.Finish()
	ctx = outgoingToIncomingCtx(ctx)
	return c.Server.GetFallbackName(ctx, in)
}

func (c *GreeterClient) Reset(inctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*empty.Empty, error) {
	span, ctx := tracing.StartSpan(inctx, tracing.FuncName(tracing.FuncNameWithParentType))
	defer span.Finish()
	ctx = outgoingToIncomingCtx(ctx)
	return c.Server.Reset(ctx, in)
}

// outgoingToIncomingCtx transforms the outgoing context of a grpc client to
// an incoming context for a grpc server.
// this is needed since the client is not actually a grpc client but a wrapper for the server
func outgoingToIncomingCtx(ctx context.Context) context.Context {
	md, _ := metadata.FromOutgoingContext(ctx)
	return metadata.NewIncomingContext(ctx, md)
}
