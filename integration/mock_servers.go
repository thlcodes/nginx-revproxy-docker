// Package integrationtests implements the mock servers for the integration tests.
package integration

import (
	"net"

	"dev.beta.audi/gorepo/lib_proto_models/golib/registry"
	"dev.beta.audi/gorepo/lib_proto_models/golib/storage"

	"google.golang.org/grpc"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

/*
CreateRegistryServerMock starts a goroutine with a gRPC server for the registry service.

	Parameters:
		- registryMock: a mock implementation of the registry service
		- address: the network address for the server

	Returns:
		- server: the reference to the gRPC server, e.g., for later shutdown
*/
func CreateRegistryServerMock(registryMock registry.RegistryServiceServer, address string) *grpc.Server {
	listener, err := net.Listen("tcp", address)
	gomega.Expect(err).NotTo(gomega.HaveOccurred(), "failed to listen: %v", err)

	server := grpc.NewServer()
	registry.RegisterRegistryServiceServer(server, registryMock)

	go func() {
		err := server.Serve(listener)
		gomega.Expect(err).NotTo(gomega.HaveOccurred(), "failed to serve: %v", err)
		ginkgo.GinkgoT().Log("Stopping RegistryServerMock")
	}()

	return server
}

/*
CreateStorageServerMock starts a goroutine with a gRPC server for the storage service.

	Parameters:
		- storageMock: a mock implementation of the storage service
		- address: the network address for the server

	Returns:
		- server: the reference to the gRPC server
*/
func CreateStorageServerMock(storageMock storage.StorageServiceServer, address string) *grpc.Server {
	listener, err := net.Listen("tcp", address)
	gomega.Expect(err).NotTo(gomega.HaveOccurred(), "failed to listen: %v", err)

	server := grpc.NewServer()
	storage.RegisterStorageServiceServer(server, storageMock)

	go func() {
		err := server.Serve(listener)
		gomega.Expect(err).NotTo(gomega.HaveOccurred(), "failed to serve: %v", err)
		ginkgo.GinkgoT().Log("Stopping StorageServerMock")
	}()

	return server
}
