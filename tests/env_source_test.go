package tests

import (
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/dig"

	flam "github.com/happyhippyhippo/flam"
	config "github.com/happyhippyhippo/flam-config"
	filesystem "github.com/happyhippyhippo/flam-filesystem"
	time "github.com/happyhippyhippo/flam-time"
)

func Test_envSource(t *testing.T) {
	t.Run("should return env source file loading error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(config.PathSources, flam.Bag{
			"defaults": flam.Bag{
				"driver":   config.SourceDriverEnv,
				"priority": 123,
				"files":    []string{"./testdata/invalid"},
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		assert.ErrorContains(
			t,
			config.NewProvider().(flam.BootableProvider).Boot(container),
			"no such file or directory")
	})

	t.Run("should boot a valid env source with mappings and loaded env files", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		require.NoError(t, os.Setenv("ENV_SYSTEM_FIELD", "system_value"))
		defer func() { _ = os.Unsetenv("ENV_SYSTEM_FIELD") }()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(config.PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   config.SourceDriverEnv,
				"priority": 123,
				"files":    []string{"./testdata/env"},
				"mappings": flam.Bag{
					"ENV_FILE_FIELD":     "env.file_field",
					"ENV_SYSTEM_FIELD":   "env.system_field",
					"ENV_SYSTEM_INVALID": "env.invalid",
				},
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			got, e := facade.GetSource("my_source")
			require.NotNil(t, got)
			require.NoError(t, e)

			assert.Equal(t, 123, got.GetPriority())
			assert.Equal(t, "file_value", got.Get("env.file_field"))
			assert.Equal(t, "system_value", got.Get("env.system_field"))
			assert.Equal(t, nil, got.Get("env.invalid"))
		}))
	})
}
