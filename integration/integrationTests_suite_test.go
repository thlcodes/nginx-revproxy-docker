// +build integration

// Package integrationtests_test implements the bootstrapping for the integration tests.
package integration_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"dev.beta.audi/gorepo/lib-go-common/http"
	ecomypb "dev.beta.audi/gorepo/lib_proto_models/golib/ecomy"
	registrypb "dev.beta.audi/gorepo/lib_proto_models/golib/registry"

	integrationtests "dev.beta.audi/gorepo/gopher-user-ecomy/integration"
	"dev.beta.audi/gorepo/gopher-user-ecomy/internal/configuration"
	"dev.beta.audi/gorepo/gopher-user-ecomy/internal/generated/config"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// The variables will be captured by Ginkgo's block closures, e.g.,
// for later graceful shutdown and for setting mock expectations in the tests.
var (
	cfg = &config.Envs.Test
	err error

	service          *configuration.Service
	statusClient     ecomypb.EcomyServiceClient
	clientConnection *grpc.ClientConn

	registryMock       = &registrypb.MockRegistryServiceServer{}
	registryMockServer *grpc.Server

	// Add your mock server variables here
	//storageMock       = &storage.MockStorageServiceServer{}
	//storageMockServer *grpc.Server
)

/*
Configure free ports for the test environment.
Start the mock servers on which the service depends.
Start the service under test.
Create a client connection with which to test the service.
*/
var _ = BeforeSuite(func() {

	cfg.Port = int64(http.GetFreePort())
	cfg.GrpcPort = int64(http.GetFreePort())
	cfg.Services.Registry = cfg.Host + ":" + strconv.Itoa(http.GetFreePort())
	//cfg.Services.Storage.Addr = cfg.Host + ":" + strconv.Itoa(http.GetFreePort())

	registryMockServer = integrationtests.CreateRegistryServerMock(registryMock, cfg.Services.Registry)
	registryMock.On("AddServiceInstance", mock.Anything, mock.Anything).Return(&empty.Empty{}, nil)
	registryMock.On("RemoveServiceInstance", mock.Anything, mock.Anything).Return(&empty.Empty{}, nil)

	// Add and start your mock servers here
	//storageMockServer = integrationtests.CreateStorageServerMock(storageMock, cfg.Services.Storage.Addr)

	// Start the service as it is done in main()
	go func() {
		service, err = configuration.NewService(cfg)
		Expect(err).NotTo(HaveOccurred(), "failed to instantiate service: %v", err)

		err = service.Start()
		Expect(err).NotTo(HaveOccurred(), "failed to start service: %v", err)
	}()

	// Wait for the service to start
	time.Sleep(1 * time.Second)

	var address = fmt.Sprintf(cfg.Host+":%v", cfg.GrpcPort)
	clientConnection, err = grpc.Dial(address, grpc.WithInsecure())
	statusClient = ecomypb.NewEcomyServiceClient(clientConnection)
	Expect(err).NotTo(HaveOccurred(), "did not connect: %v", err)
})

// Close the client connection to the service under test and terminate the goroutines of the mock servers and of the service.
var _ = AfterSuite(func() {

	err = clientConnection.Close()
	Expect(err).NotTo(HaveOccurred(), "failed to close connection")

	// Stop your mock servers here
	//storageMockServer.GracefulStop()

	err = service.Stop()
	Expect(err).NotTo(HaveOccurred(), "failed to stop service: %v", err)

	registryMockServer.GracefulStop()
	registryMock.AssertExpectations(GinkgoT())
})

func TestIntegrationTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "IntegrationTests Suite")
}
