//go:generate task genfig
//go:generate task mocks
package main

import (
	"dev.beta.audi/gorepo/gopher_skeleton/internal/generated/config"

	"dev.beta.audi/gorepo/lib-go-common/common"

	"dev.beta.audi/gorepo/gopher_skeleton/internal/configuration"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			common.HandleGlobalPanic(r, true)
		}
	}()

	service, err := configuration.NewService(config.Current)

	if err != nil {
		panic(err)
	}
	panic(service.Start())
}
