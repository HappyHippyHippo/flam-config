package tests

import (
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/dig"

	flam "github.com/happyhippyhippo/flam"
	config "github.com/happyhippyhippo/flam-config"
	mocks "github.com/happyhippyhippo/flam-config/tests/mocks"
	filesystem "github.com/happyhippyhippo/flam-filesystem"
	time "github.com/happyhippyhippo/flam-time"
)

func Test_restSource(t *testing.T) {
	t.Run("should ignore config without uri field", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(config.PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   config.SourceDriverRest,
				"parser":   "my_parser",
				"path":     flam.Bag{"config": "config"},
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		assert.ErrorIs(
			t,
			config.NewProvider().(flam.BootableProvider).Boot(container),
			flam.ErrInvalidResourceConfig)
	})

	t.Run("should ignore config without path.config field", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(config.PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   config.SourceDriverRest,
				"parser":   "my_parser",
				"uri":      "http://path/",
				"path":     flam.Bag{},
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		assert.ErrorIs(
			t,
			config.NewProvider().(flam.BootableProvider).Boot(container),
			flam.ErrInvalidResourceConfig)
	})

	t.Run("should return requester generation error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(config.PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   config.SourceDriverRest,
				"parser":   "my_parser",
				"uri":      "http://path/",
				"path":     flam.Bag{"config": "config"},
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		expectedErr := errors.New("requester error")
		requestGenerator := mocks.NewRestRequesterGenerator(ctrl)
		requestGenerator.EXPECT().Create().Return(nil, expectedErr).Times(1)
		require.NoError(t, container.Decorate(func(generator config.RestRequesterGenerator) config.RestRequesterGenerator {
			return requestGenerator
		}))

		assert.ErrorIs(
			t,
			config.NewProvider().(flam.BootableProvider).Boot(container),
			expectedErr)
	})

	t.Run("should return parser generation error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(config.PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   config.SourceDriverRest,
				"parser":   "my_parser",
				"uri":      "http://path/",
				"path":     flam.Bag{"config": "config"},
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		assert.ErrorIs(
			t,
			config.NewProvider().(flam.BootableProvider).Boot(container),
			flam.ErrUnknownResource)
	})

	t.Run("should return request generation error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(config.PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": config.ParserDriverJson,
			}})
		_ = config.Defaults.Set(config.PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   config.SourceDriverRest,
				"parser":   "my_parser",
				"uri":      ":/uri",
				"path":     flam.Bag{"config": "config"},
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		assert.ErrorContains(
			t,
			config.NewProvider().(flam.BootableProvider).Boot(container),
			"missing protocol scheme")
	})

	t.Run("should return requester error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(config.PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": config.ParserDriverJson,
			}})
		_ = config.Defaults.Set(config.PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   config.SourceDriverRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config"},
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		expectedErr := errors.New("requester error")
		requester := mocks.NewRestRequester(ctrl)
		requester.EXPECT().Do(gomock.Any()).Return(nil, expectedErr).Times(1)

		requestGenerator := mocks.NewRestRequesterGenerator(ctrl)
		requestGenerator.EXPECT().Create().Return(requester, nil).Times(1)
		require.NoError(t, container.Decorate(func(generator config.RestRequesterGenerator) config.RestRequesterGenerator {
			return requestGenerator
		}))

		assert.ErrorIs(
			t,
			config.NewProvider().(flam.BootableProvider).Boot(container),
			expectedErr)
	})

	t.Run("should return response read error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(config.PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": config.ParserDriverJson,
			}})
		_ = config.Defaults.Set(config.PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   config.SourceDriverRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config"},
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		expectedErr := errors.New("requester error")
		body := mocks.NewReadCloser(ctrl)
		body.EXPECT().Read(gomock.Any()).Return(0, expectedErr).Times(1)

		response := &http.Response{Body: body}

		requester := mocks.NewRestRequester(ctrl)
		requester.EXPECT().Do(gomock.Any()).Return(response, nil).Times(1)

		requestGenerator := mocks.NewRestRequesterGenerator(ctrl)
		requestGenerator.EXPECT().Create().Return(requester, nil).Times(1)
		require.NoError(t, container.Decorate(func(generator config.RestRequesterGenerator) config.RestRequesterGenerator {
			return requestGenerator
		}))

		assert.ErrorIs(
			t,
			config.NewProvider().(flam.BootableProvider).Boot(container),
			expectedErr)
	})

	t.Run("should return response parsing error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(config.PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": config.ParserDriverJson,
			}})
		_ = config.Defaults.Set(config.PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   config.SourceDriverRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config"},
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		data := "{"
		reader := func(b []byte) (int, error) {
			copy(b, data)
			return len(data), io.EOF
		}

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		body := mocks.NewReadCloser(ctrl)
		body.EXPECT().Read(gomock.Any()).DoAndReturn(reader).Times(1)

		response := &http.Response{Body: body}

		requester := mocks.NewRestRequester(ctrl)
		requester.EXPECT().Do(gomock.Any()).Return(response, nil).Times(1)

		requestGenerator := mocks.NewRestRequesterGenerator(ctrl)
		requestGenerator.EXPECT().Create().Return(requester, nil).Times(1)
		require.NoError(t, container.Decorate(func(generator config.RestRequesterGenerator) config.RestRequesterGenerator {
			return requestGenerator
		}))

		assert.ErrorContains(
			t,
			config.NewProvider().(flam.BootableProvider).Boot(container),
			"unexpected end of JSON input")
	})

	t.Run("should return config not found in response error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(config.PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": config.ParserDriverJson,
			}})
		_ = config.Defaults.Set(config.PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   config.SourceDriverRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config"},
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		data := "{}"
		reader := func(b []byte) (int, error) {
			copy(b, data)
			return len(data), io.EOF
		}

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		body := mocks.NewReadCloser(ctrl)
		body.EXPECT().Read(gomock.Any()).DoAndReturn(reader).Times(1)

		response := &http.Response{Body: body}

		requester := mocks.NewRestRequester(ctrl)
		requester.EXPECT().Do(gomock.Any()).Return(response, nil).Times(1)

		requestGenerator := mocks.NewRestRequesterGenerator(ctrl)
		requestGenerator.EXPECT().Create().Return(requester, nil).Times(1)
		require.NoError(t, container.Decorate(func(generator config.RestRequesterGenerator) config.RestRequesterGenerator {
			return requestGenerator
		}))

		assert.ErrorIs(
			t,
			config.NewProvider().(flam.BootableProvider).Boot(container),
			config.ErrRestConfigNotFound)
	})

	t.Run("should return invalid config in response error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(config.PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": config.ParserDriverJson,
			}})
		_ = config.Defaults.Set(config.PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   config.SourceDriverRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config"},
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		data := "{\"config\": \"invalid\"}"
		reader := func(b []byte) (int, error) {
			copy(b, data)
			return len(data), io.EOF
		}

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		body := mocks.NewReadCloser(ctrl)
		body.EXPECT().Read(gomock.Any()).DoAndReturn(reader).Times(1)

		response := &http.Response{Body: body}

		requester := mocks.NewRestRequester(ctrl)
		requester.EXPECT().Do(gomock.Any()).Return(response, nil).Times(1)

		requestGenerator := mocks.NewRestRequesterGenerator(ctrl)
		requestGenerator.EXPECT().Create().Return(requester, nil).Times(1)
		require.NoError(t, container.Decorate(func(generator config.RestRequesterGenerator) config.RestRequesterGenerator {
			return requestGenerator
		}))

		assert.ErrorIs(
			t,
			config.NewProvider().(flam.BootableProvider).Boot(container),
			config.ErrRestInvalidConfig)
	})

	t.Run("should correctly load the config", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(config.PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": config.ParserDriverJson,
			}})
		_ = config.Defaults.Set(config.PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   config.SourceDriverRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config"},
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		data := "{\"config\": {\"field\": \"value\"}}"
		reader := func(b []byte) (int, error) {
			copy(b, data)
			return len(data), io.EOF
		}

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		body := mocks.NewReadCloser(ctrl)
		body.EXPECT().Read(gomock.Any()).DoAndReturn(reader).Times(1)

		response := &http.Response{Body: body}

		requester := mocks.NewRestRequester(ctrl)
		requester.EXPECT().Do(gomock.Any()).Return(response, nil).Times(1)

		requestGenerator := mocks.NewRestRequesterGenerator(ctrl)
		requestGenerator.EXPECT().Create().Return(requester, nil).Times(1)
		require.NoError(t, container.Decorate(func(generator config.RestRequesterGenerator) config.RestRequesterGenerator {
			return requestGenerator
		}))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			got, e := facade.GetSource("my_source")
			require.NotNil(t, got)
			require.NoError(t, e)

			assert.Equal(t, "value", got.Get("field"))
		}))
	})
}
