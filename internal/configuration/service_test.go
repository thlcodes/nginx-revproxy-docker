package configuration_test

import (
	"os"
	"testing"
	"time"

	httplib "dev.beta.audi/gorepo/lib-go-common/http"

	"dev.beta.audi/gorepo/gopher_skeleton/internal/configuration"
	"dev.beta.audi/gorepo/gopher_skeleton/internal/generated/config"

	"github.com/stretchr/testify/require"
)

var (
	cfg = &config.Envs.Test
)

func TestMain(m *testing.M) {
	cfg.Port = int64(httplib.GetFreePort())
	os.Exit(m.Run())
}

func Test_Service(t *testing.T) {
	service, err := configuration.NewService(cfg)
	require.NoError(t, err)
	require.NotNil(t, service)

	go func() {
		time.Sleep(3 * time.Second)
		require.NoError(t, service.Stop())
	}()
	require.NoError(t, service.Start())
}
