// Package configuration configures the service and the client connections to other services it depends on.
package configuration

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	storagepb "dev.beta.audi/gorepo/lib_proto_models/golib/storage"

	"dev.beta.audi/gorepo/gopher_skeleton/internal/repositories"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/opentracing/opentracing-go"

	"dev.beta.audi/gorepo/lib_proto_models/golib/greeter"

	"google.golang.org/grpc"

	grpcdelivery "dev.beta.audi/gorepo/gopher_skeleton/internal/delivery/grpc"
	"dev.beta.audi/gorepo/gopher_skeleton/internal/generated/config"
	"dev.beta.audi/gorepo/gopher_skeleton/internal/usecases"

	"dev.beta.audi/gorepo/lib-go-common/bootstrap"
	grpclib "dev.beta.audi/gorepo/lib-go-common/grpc"
	"dev.beta.audi/gorepo/lib-go-common/tracing"
)

// Service configuration
type Service struct {
	*bootstrap.Service
}

/*
NewService instantiates a new service.

	Parameters:
		- config: The service configuration.

	Returns:
		- service: The service object. Nil if configuring the service fails.
		- error:   The error describes why configuring the service failed.
*/
func NewService(config *config.Config) (*Service, error) {
	log.Printf("NewService with config: %#v", config)

	service := bootstrap.NewService()

	_, err := setupTracing(config)
	if err != nil {
		log.Printf("WARNING: tracer could not be set up: %v", err)
	}

	/*regConn, regClient, err := connectToRegistry(config.Registry)
	if err != nil {
		return nil, err
	}
	// close connection to registry when this method is finished,
	// since we do not need connection past this method
	defer regConn.Close()
	sd := discovery.NewGRPCDiscovery(regClient)
	*/

	storageConn, storageClient, err := connectToStorageAL( /*sd, */ config.Services.Storage.Addr)
	if err != nil {
		return nil, err
	}

	greeterServer := grpcdelivery.NewGreeterServer(
		config,
		usecases.NewGreeter(
			config,
			repositories.NewPersistentStore(
				config,
				storageClient,
			),
		),
	)

	service.On(bootstrap.StoppedEvent, func() {
		storageConn.Close()
	})

	if err := service.Configure(
		bootstrap.WithNameAndVersion(config.Service.Name, config.Service.Version),
		bootstrap.WithServiceConfig(
			bootstrap.ServiceConfig{
				Host:         config.Host,
				HTTPPort:     fmt.Sprintf("%d", config.Port),
				GRPCPort:     fmt.Sprintf("%d", config.GrpcPort),
				RegistryAddr: config.Registry,
			},
		),
		bootstrap.WithGRPCServer(func(s *grpc.Server) {
			greeter.RegisterGreeterServiceServer(s, greeterServer)
		}),
		bootstrap.WithGRPCGateway(func(ctx context.Context, mux *runtime.ServeMux) {
			cc, err := grpc.Dial(":"+strconv.FormatInt(config.GrpcPort, 10), grpc.WithInsecure())
			if err != nil {
				panic(err)
			}
			_ = greeter.RegisterGreeterServiceHandler(ctx, mux, cc)
			//_ = greeter.RegisterGreeterServiceHandlerClient(ctx, mux, &grpcdelivery.GreeterClient{Server: greeterServer})
		}),
		bootstrap.WithAPIDocs(greeter.GetServiceSpecBytes),
		bootstrap.WithDefaultHTTPServer(),
		bootstrap.WithDefaultRegistration(),
		bootstrap.WithDefaultGraceful(),
	); err != nil {
		return nil, err
	}

	return &Service{
		Service: service,
	}, nil
}

/*func connectToRegistry(addr string) (conn *grpc.ClientConn, client registrypb.RegistryServiceClient, err error) {
	conn, err = grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return
	}
	client = registrypb.NewRegistryServiceClient(conn)
	return
}*/

func connectToStorageAL( /*sd discovery.ServiceDiscovery, */ addr string) (conn *grpc.ClientConn, client storagepb.StorageServiceClient, err error) {
	/*var uri *url.URL
	uri, _, err = sd.ResolveURI(addr)
	if err != nil {
		return
	}*/
	conn, err = grpclib.DefaultConnectTo(addr)
	if err != nil {
		return
	}
	client = storagepb.NewStorageServiceClient(conn)
	return
}

func setupTracing(config *config.Config) (tracer opentracing.Tracer, err error) {
	uri := config.Tracing.Uri
	secret := config.Tracing.Secret
	if config.Tracing.Service != "" && config.Trace {
		var ok bool
		if uri, secret, ok = tracing.GetAPMConfigFromCFService(config.Tracing.Service); !ok {
			tracer, _ = tracing.SetupNoopTracer()
			err = errors.New("could not read tracing config from CF service")
			return
		}
	}
	env := os.Getenv("ENV")
	if env == "" {
		env = "unknown"
	}
	tracer, err = tracing.SetupDefaultTracer(config.Trace, config.Service.Name, config.Service.Version, env, uri, secret)
	if err != nil {
		tracer, _ = tracing.SetupNoopTracer()
	}
	return
}
