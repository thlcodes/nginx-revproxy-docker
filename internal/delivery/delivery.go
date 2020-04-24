// Package delivery defines the API of the delivery layer.
package delivery

import (
	greeter "dev.beta.audi/gorepo/lib_proto_models/golib/greeter"
)

// GreeterServer is the server API.
type GreeterServer = greeter.GreeterServiceServer

// GreeterClient is the client API.
type GreeterClient = greeter.GreeterServiceClient
