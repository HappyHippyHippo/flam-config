package tests

import (
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/dig"

	flam "github.com/happyhippyhippo/flam"
	config "github.com/happyhippyhippo/flam-config"
	mocks "github.com/happyhippyhippo/flam-config/tests/mocks"
	filesystem "github.com/happyhippyhippo/flam-filesystem"
	flamTime "github.com/happyhippyhippo/flam-time"
)

func Test_observableRestSource(t *testing.T) {
	t.Run("should ignore config without uri field", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(config.PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   config.SourceDriverObservableRest,
				"parser":   "my_parser",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
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
				"driver":   config.SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://path/",
				"path":     flam.Bag{"timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		assert.ErrorIs(
			t,
			config.NewProvider().(flam.BootableProvider).Boot(container),
			flam.ErrInvalidResourceConfig)
	})

	t.Run("should ignore config without path.timestamp field", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(config.PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   config.SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://path/",
				"path":     flam.Bag{"config": "config"},
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
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
				"driver":   config.SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://path/",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
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
				"driver":   config.SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://path/",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		requester := mocks.NewRestRequester(ctrl)

		requestGenerator := mocks.NewRestRequesterGenerator(ctrl)
		requestGenerator.EXPECT().Create().Return(requester, nil).Times(1)
		require.NoError(t, container.Decorate(func(generator config.RestRequesterGenerator) config.RestRequesterGenerator {
			return requestGenerator
		}))

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
				"driver":   config.SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      ":/uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		requester := mocks.NewRestRequester(ctrl)

		requestGenerator := mocks.NewRestRequesterGenerator(ctrl)
		requestGenerator.EXPECT().Create().Return(requester, nil).Times(1)
		require.NoError(t, container.Decorate(func(generator config.RestRequesterGenerator) config.RestRequesterGenerator {
			return requestGenerator
		}))

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
				"driver":   config.SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
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
				"driver":   config.SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
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
				"driver":   config.SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		data := "{"
		reader := func(b []byte) (int, error) {
			copy(b, data)
			return len(data), io.EOF
		}

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

	t.Run("should return timestamp not found in response error", func(t *testing.T) {
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
				"driver":   config.SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		data := "{}"
		reader := func(b []byte) (int, error) {
			copy(b, data)
			return len(data), io.EOF
		}

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
			config.ErrRestTimestampNotFound)
	})

	t.Run("should return invalid timestamp in response error (type)", func(t *testing.T) {
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
				"driver":   config.SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		data := "{\"timestamp\": 123}"
		reader := func(b []byte) (int, error) {
			copy(b, data)
			return len(data), io.EOF
		}

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
			config.ErrRestInvalidTimestamp)
	})

	t.Run("should return invalid timestamp in response error (string parsing)", func(t *testing.T) {
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
				"driver":   config.SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		expectedErr := errors.New("parse error")
		timeFacade := mocks.NewTimeFacade(ctrl)
		timeFacade.EXPECT().Now().Return(time.Now()).Times(1)
		timeFacade.EXPECT().Parse(time.RFC3339, "invalid").Return(time.Time{}, expectedErr).Times(1)
		require.NoError(t, container.Provide(func() flamTime.Facade { return timeFacade }))

		data := "{\"timestamp\": \"invalid\"}"
		reader := func(b []byte) (int, error) {
			copy(b, data)
			return len(data), io.EOF
		}

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
			expectedErr)
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
				"driver":   config.SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		now := time.Now()
		timeFacade := mocks.NewTimeFacade(ctrl)
		timeFacade.EXPECT().Now().Return(now).Times(1)
		timeFacade.EXPECT().
			Parse(time.RFC3339, "1234-56-78 90:12:34 +0000 UTC").
			Return(now.Add(time.Hour*24), nil).
			Times(1)
		timeFacade.EXPECT().Unix(int64(0), int64(0)).Return(time.Unix(0, 0)).Times(1)
		require.NoError(t, container.Provide(func() flamTime.Facade { return timeFacade }))

		data := "{\"timestamp\": \"1234-56-78 90:12:34 +0000 UTC\"}"
		reader := func(b []byte) (int, error) {
			copy(b, data)
			return len(data), io.EOF
		}

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
				"driver":   config.SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		now := time.Now()
		timeFacade := mocks.NewTimeFacade(ctrl)
		timeFacade.EXPECT().Now().Return(now).Times(1)
		timeFacade.EXPECT().
			Parse(time.RFC3339, "1234-56-78 90:12:34 +0000 UTC").
			Return(now.Add(time.Hour*24), nil).
			Times(1)
		timeFacade.EXPECT().Unix(int64(0), int64(0)).Return(time.Unix(0, 0)).Times(1)
		require.NoError(t, container.Provide(func() flamTime.Facade { return timeFacade }))

		data := "{\"timestamp\": \"1234-56-78 90:12:34 +0000 UTC\", \"config\": 123}"
		reader := func(b []byte) (int, error) {
			copy(b, data)
			return len(data), io.EOF
		}

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
				"driver":   config.SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		now := time.Now()
		timeFacade := mocks.NewTimeFacade(ctrl)
		timeFacade.EXPECT().Now().Return(now).Times(1)
		timeFacade.EXPECT().
			Parse(time.RFC3339, "1234-56-78 90:12:34 +0000 UTC").
			Return(now.Add(time.Hour*24), nil).
			Times(1)
		timeFacade.EXPECT().Unix(int64(0), int64(0)).Return(time.Unix(0, 0)).Times(1)
		require.NoError(t, container.Provide(func() flamTime.Facade { return timeFacade }))

		data := "{\"timestamp\": \"1234-56-78 90:12:34 +0000 UTC\", \"config\": {\"field\": \"value\"}}"
		reader := func(b []byte) (int, error) {
			copy(b, data)
			return len(data), io.EOF
		}

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
			assert.NotNil(t, got)
			assert.NoError(t, e)

			assert.Equal(t, "value", got.Get("field"))
		}))
	})
}

func Test_observableRestSource_Reload(t *testing.T) {
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
				"driver":   config.SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		now := time.Now()
		timeFacade := mocks.NewTimeFacade(ctrl)
		timeFacade.EXPECT().Now().Return(now).Times(1)
		timeFacade.EXPECT().
			Parse(time.RFC3339, "1234-56-78 90:12:34 +0000 UTC").
			Return(now.Add(time.Hour*24), nil).
			Times(1)
		timeFacade.EXPECT().Unix(int64(0), int64(0)).Return(time.Unix(0, 0)).Times(1)
		require.NoError(t, container.Provide(func() flamTime.Facade { return timeFacade }))

		data := "{\"timestamp\": \"1234-56-78 90:12:34 +0000 UTC\", \"config\": {\"field\": \"value\"}}"
		reader := func(b []byte) (int, error) {
			copy(b, data)
			return len(data), io.EOF
		}

		body := mocks.NewReadCloser(ctrl)
		body.EXPECT().Read(gomock.Any()).DoAndReturn(reader).Times(1)

		response := &http.Response{Body: body}

		expectedErr := errors.New("requester error")
		requester := mocks.NewRestRequester(ctrl)
		requester.EXPECT().Do(gomock.Any()).Return(response, nil)
		requester.EXPECT().Do(gomock.Any()).Return(nil, expectedErr)

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

			reloaded, e := got.(config.ObservableSource).Reload()
			assert.False(t, reloaded)
			assert.ErrorIs(t, e, expectedErr)
		}))
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
				"driver":   config.SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		now := time.Now()
		timeFacade := mocks.NewTimeFacade(ctrl)
		timeFacade.EXPECT().Now().Return(now).Times(1)
		timeFacade.EXPECT().
			Parse(time.RFC3339, "1234-56-78 90:12:34 +0000 UTC").
			Return(now.Add(time.Hour*24), nil).
			Times(1)
		timeFacade.EXPECT().Unix(int64(0), int64(0)).Return(time.Unix(0, 0)).Times(1)
		require.NoError(t, container.Provide(func() flamTime.Facade { return timeFacade }))

		data := "{\"timestamp\": \"1234-56-78 90:12:34 +0000 UTC\", \"config\": {\"field\": \"value\"}}"
		reader := func(b []byte) (int, error) {
			copy(b, data)
			return len(data), io.EOF
		}

		body1 := mocks.NewReadCloser(ctrl)
		body1.EXPECT().Read(gomock.Any()).DoAndReturn(reader).Times(1)

		expectedErr := errors.New("reader error")
		body2 := mocks.NewReadCloser(ctrl)
		body2.EXPECT().Read(gomock.Any()).Return(0, expectedErr).Times(1)

		response1 := &http.Response{Body: body1}

		response2 := &http.Response{Body: body2}

		requester := mocks.NewRestRequester(ctrl)
		requester.EXPECT().Do(gomock.Any()).Return(response1, nil)
		requester.EXPECT().Do(gomock.Any()).Return(response2, nil)

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

			reloaded, e := got.(config.ObservableSource).Reload()
			assert.False(t, reloaded)
			assert.ErrorIs(t, e, expectedErr)
		}))
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
				"driver":   config.SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		now := time.Now()
		timeFacade := mocks.NewTimeFacade(ctrl)
		timeFacade.EXPECT().Now().Return(now).Times(1)
		timeFacade.EXPECT().
			Parse(time.RFC3339, "1234-56-78 90:12:34 +0000 UTC").
			Return(now.Add(time.Hour*24), nil).
			Times(1)
		timeFacade.EXPECT().Unix(int64(0), int64(0)).Return(time.Unix(0, 0)).Times(1)
		require.NoError(t, container.Provide(func() flamTime.Facade { return timeFacade }))

		data1 := "{\"timestamp\": \"1234-56-78 90:12:34 +0000 UTC\", \"config\": {\"field\": \"value\"}}"
		reader1 := func(b []byte) (int, error) {
			copy(b, data1)
			return len(data1), io.EOF
		}

		body1 := mocks.NewReadCloser(ctrl)
		body1.EXPECT().Read(gomock.Any()).DoAndReturn(reader1).Times(1)

		data2 := "{"
		reader2 := func(b []byte) (int, error) {
			copy(b, data2)
			return len(data2), io.EOF
		}

		body2 := mocks.NewReadCloser(ctrl)
		body2.EXPECT().Read(gomock.Any()).DoAndReturn(reader2).Times(1)

		response1 := &http.Response{Body: body1}

		response2 := &http.Response{Body: body2}

		requester := mocks.NewRestRequester(ctrl)
		requester.EXPECT().Do(gomock.Any()).Return(response1, nil)
		requester.EXPECT().Do(gomock.Any()).Return(response2, nil)

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

			reloaded, e := got.(config.ObservableSource).Reload()
			assert.False(t, reloaded)
			assert.ErrorContains(t, e, "unexpected end of JSON input")
		}))
	})

	t.Run("should return timestamp not found in response error", func(t *testing.T) {
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
				"driver":   config.SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		now := time.Now()
		timeFacade := mocks.NewTimeFacade(ctrl)
		timeFacade.EXPECT().Now().Return(now).Times(1)
		timeFacade.EXPECT().
			Parse(time.RFC3339, "1234-56-78 90:12:34 +0000 UTC").
			Return(now.Add(time.Hour*24), nil).
			Times(1)
		timeFacade.EXPECT().Unix(int64(0), int64(0)).Return(time.Unix(0, 0)).Times(1)
		require.NoError(t, container.Provide(func() flamTime.Facade { return timeFacade }))

		data1 := "{\"timestamp\": \"1234-56-78 90:12:34 +0000 UTC\", \"config\": {\"field\": \"value\"}}"
		reader1 := func(b []byte) (int, error) {
			copy(b, data1)
			return len(data1), io.EOF
		}

		body1 := mocks.NewReadCloser(ctrl)
		body1.EXPECT().Read(gomock.Any()).DoAndReturn(reader1).Times(1)

		data2 := "{}"
		reader2 := func(b []byte) (int, error) {
			copy(b, data2)
			return len(data2), io.EOF
		}

		body2 := mocks.NewReadCloser(ctrl)
		body2.EXPECT().Read(gomock.Any()).DoAndReturn(reader2).Times(1)

		response1 := &http.Response{Body: body1}

		response2 := &http.Response{Body: body2}

		requester := mocks.NewRestRequester(ctrl)
		requester.EXPECT().Do(gomock.Any()).Return(response1, nil)
		requester.EXPECT().Do(gomock.Any()).Return(response2, nil)

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

			reloaded, e := got.(config.ObservableSource).Reload()
			assert.False(t, reloaded)
			assert.ErrorIs(t, e, config.ErrRestTimestampNotFound)
		}))
	})

	t.Run("should return invalid timestamp in response error (type)", func(t *testing.T) {
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
				"driver":   config.SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		now := time.Now()
		timeFacade := mocks.NewTimeFacade(ctrl)
		timeFacade.EXPECT().Now().Return(now).Times(1)
		timeFacade.EXPECT().
			Parse(time.RFC3339, "1234-56-78 90:12:34 +0000 UTC").
			Return(now.Add(time.Hour*24), nil).
			Times(1)
		timeFacade.EXPECT().Unix(int64(0), int64(0)).Return(time.Unix(0, 0)).Times(1)
		require.NoError(t, container.Provide(func() flamTime.Facade { return timeFacade }))

		data1 := "{\"timestamp\": \"1234-56-78 90:12:34 +0000 UTC\", \"config\": {\"field\": \"value\"}}"
		reader1 := func(b []byte) (int, error) {
			copy(b, data1)
			return len(data1), io.EOF
		}

		body1 := mocks.NewReadCloser(ctrl)
		body1.EXPECT().Read(gomock.Any()).DoAndReturn(reader1).Times(1)

		data2 := "{\"timestamp\": 1234567890}"
		reader2 := func(b []byte) (int, error) {
			copy(b, data2)
			return len(data2), io.EOF
		}

		body2 := mocks.NewReadCloser(ctrl)
		body2.EXPECT().Read(gomock.Any()).DoAndReturn(reader2).Times(1)

		response1 := &http.Response{Body: body1}

		response2 := &http.Response{Body: body2}

		requester := mocks.NewRestRequester(ctrl)
		requester.EXPECT().Do(gomock.Any()).Return(response1, nil)
		requester.EXPECT().Do(gomock.Any()).Return(response2, nil)

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

			reloaded, e := got.(config.ObservableSource).Reload()
			assert.False(t, reloaded)
			assert.ErrorIs(t, e, config.ErrRestInvalidTimestamp)
		}))
	})

	t.Run("should return invalid timestamp in response error (string parsing)", func(t *testing.T) {
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
				"driver":   config.SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		expectedErr := errors.New("invalid timestamp")
		now := time.Now()
		timeFacade := mocks.NewTimeFacade(ctrl)
		timeFacade.EXPECT().Now().Return(now).Times(1)
		timeFacade.EXPECT().
			Parse(time.RFC3339, "1234-56-78 90:12:34 +0000 UTC").
			Return(now.Add(time.Hour*24), nil)
		timeFacade.EXPECT().
			Parse(time.RFC3339, "invalid").
			Return(time.Time{}, expectedErr)
		timeFacade.EXPECT().Unix(int64(0), int64(0)).Return(time.Unix(0, 0)).Times(1)
		require.NoError(t, container.Provide(func() flamTime.Facade { return timeFacade }))

		data1 := "{\"timestamp\": \"1234-56-78 90:12:34 +0000 UTC\", \"config\": {\"field\": \"value\"}}"
		reader1 := func(b []byte) (int, error) {
			copy(b, data1)
			return len(data1), io.EOF
		}

		body1 := mocks.NewReadCloser(ctrl)
		body1.EXPECT().Read(gomock.Any()).DoAndReturn(reader1).Times(1)

		data2 := "{\"timestamp\": \"invalid\"}"
		reader2 := func(b []byte) (int, error) {
			copy(b, data2)
			return len(data2), io.EOF
		}

		body2 := mocks.NewReadCloser(ctrl)
		body2.EXPECT().Read(gomock.Any()).DoAndReturn(reader2).Times(1)

		response1 := &http.Response{Body: body1}

		response2 := &http.Response{Body: body2}

		requester := mocks.NewRestRequester(ctrl)
		requester.EXPECT().Do(gomock.Any()).Return(response1, nil)
		requester.EXPECT().Do(gomock.Any()).Return(response2, nil)

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

			reloaded, e := got.(config.ObservableSource).Reload()
			assert.False(t, reloaded)
			assert.ErrorIs(t, e, expectedErr)
		}))
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
				"driver":   config.SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		now := time.Now()
		timeFacade := mocks.NewTimeFacade(ctrl)
		timeFacade.EXPECT().Now().Return(now).Times(1)
		timeFacade.EXPECT().
			Parse(time.RFC3339, "1234-56-78 90:12:34 +0000 UTC").
			Return(now.Add(time.Hour*24), nil)
		timeFacade.EXPECT().
			Parse(time.RFC3339, "1234-56-78 90:12:35 +0000 UTC").
			Return(now.Add(time.Hour*48), nil)
		timeFacade.EXPECT().Unix(int64(0), int64(0)).Return(time.Unix(0, 0)).Times(2)
		require.NoError(t, container.Provide(func() flamTime.Facade { return timeFacade }))

		data1 := "{\"timestamp\": \"1234-56-78 90:12:34 +0000 UTC\", \"config\": {\"field\": \"value\"}}"
		reader1 := func(b []byte) (int, error) {
			copy(b, data1)
			return len(data1), io.EOF
		}

		body1 := mocks.NewReadCloser(ctrl)
		body1.EXPECT().Read(gomock.Any()).DoAndReturn(reader1).Times(1)

		data2 := "{\"timestamp\": \"1234-56-78 90:12:35 +0000 UTC\"}"
		reader2 := func(b []byte) (int, error) {
			copy(b, data2)
			return len(data2), io.EOF
		}

		body2 := mocks.NewReadCloser(ctrl)
		body2.EXPECT().Read(gomock.Any()).DoAndReturn(reader2).Times(1)

		response1 := &http.Response{Body: body1}

		response2 := &http.Response{Body: body2}

		requester := mocks.NewRestRequester(ctrl)
		requester.EXPECT().Do(gomock.Any()).Return(response1, nil)
		requester.EXPECT().Do(gomock.Any()).Return(response2, nil)

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

			reloaded, e := got.(config.ObservableSource).Reload()
			assert.False(t, reloaded)
			assert.ErrorIs(t, e, config.ErrRestConfigNotFound)
		}))
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
				"driver":   config.SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		now := time.Now()
		timeFacade := mocks.NewTimeFacade(ctrl)
		timeFacade.EXPECT().Now().Return(now).Times(1)
		timeFacade.EXPECT().
			Parse(time.RFC3339, "1234-56-78 90:12:34 +0000 UTC").
			Return(now.Add(time.Hour*24), nil)
		timeFacade.EXPECT().
			Parse(time.RFC3339, "1234-56-78 90:12:35 +0000 UTC").
			Return(now.Add(time.Hour*48), nil)
		timeFacade.EXPECT().Unix(int64(0), int64(0)).Return(time.Unix(0, 0)).Times(2)
		require.NoError(t, container.Provide(func() flamTime.Facade { return timeFacade }))

		data1 := "{\"timestamp\": \"1234-56-78 90:12:34 +0000 UTC\", \"config\": {\"field\": \"value\"}}"
		reader1 := func(b []byte) (int, error) {
			copy(b, data1)
			return len(data1), io.EOF
		}

		body1 := mocks.NewReadCloser(ctrl)
		body1.EXPECT().Read(gomock.Any()).DoAndReturn(reader1).Times(1)

		data2 := "{\"timestamp\": \"1234-56-78 90:12:35 +0000 UTC\", \"config\": 123}"
		reader2 := func(b []byte) (int, error) {
			copy(b, data2)
			return len(data2), io.EOF
		}

		body2 := mocks.NewReadCloser(ctrl)
		body2.EXPECT().Read(gomock.Any()).DoAndReturn(reader2).Times(1)

		response1 := &http.Response{Body: body1}

		response2 := &http.Response{Body: body2}

		requester := mocks.NewRestRequester(ctrl)
		requester.EXPECT().Do(gomock.Any()).Return(response1, nil)
		requester.EXPECT().Do(gomock.Any()).Return(response2, nil)

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

			reloaded, e := got.(config.ObservableSource).Reload()
			assert.False(t, reloaded)
			assert.ErrorIs(t, e, config.ErrRestInvalidConfig)
		}))
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
				"driver":   config.SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		now := time.Now()
		timeFacade := mocks.NewTimeFacade(ctrl)
		timeFacade.EXPECT().Now().Return(now).Times(1)
		timeFacade.EXPECT().
			Parse(time.RFC3339, "1234-56-78 90:12:34 +0000 UTC").
			Return(now.Add(time.Hour*24), nil)
		timeFacade.EXPECT().
			Parse(time.RFC3339, "1234-56-78 90:12:35 +0000 UTC").
			Return(now.Add(time.Hour*48), nil)
		timeFacade.EXPECT().Unix(int64(0), int64(0)).Return(time.Unix(0, 0)).Times(2)
		require.NoError(t, container.Provide(func() flamTime.Facade { return timeFacade }))

		data1 := "{\"timestamp\": \"1234-56-78 90:12:34 +0000 UTC\", \"config\": {\"field\": \"value\"}}"
		reader1 := func(b []byte) (int, error) {
			copy(b, data1)
			return len(data1), io.EOF
		}

		body1 := mocks.NewReadCloser(ctrl)
		body1.EXPECT().Read(gomock.Any()).DoAndReturn(reader1).Times(1)

		data2 := "{\"timestamp\": \"1234-56-78 90:12:35 +0000 UTC\", \"config\": {\"field\": \"value2\"}}"
		reader2 := func(b []byte) (int, error) {
			copy(b, data2)
			return len(data2), io.EOF
		}

		body2 := mocks.NewReadCloser(ctrl)
		body2.EXPECT().Read(gomock.Any()).DoAndReturn(reader2).Times(1)

		response1 := &http.Response{Body: body1}

		response2 := &http.Response{Body: body2}

		requester := mocks.NewRestRequester(ctrl)
		requester.EXPECT().Do(gomock.Any()).Return(response1, nil)
		requester.EXPECT().Do(gomock.Any()).Return(response2, nil)

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

			reloaded, e := got.(config.ObservableSource).Reload()
			assert.True(t, reloaded)
			assert.NoError(t, e)

			assert.Equal(t, "value2", got.Get("field"))
		}))
	})

	t.Run("should not reload config if the timestamp is less or equals to the stored timestamp", func(t *testing.T) {
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
				"driver":   config.SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		now := time.Now()
		timeFacade := mocks.NewTimeFacade(ctrl)
		timeFacade.EXPECT().Now().Return(now).Times(1)
		timeFacade.EXPECT().
			Parse(time.RFC3339, "1234-56-78 90:12:34 +0000 UTC").
			Return(now.Add(time.Hour*24), nil)
		timeFacade.EXPECT().
			Parse(time.RFC3339, "1234-56-78 90:12:35 +0000 UTC").
			Return(now.Add(time.Hour*24), nil)
		timeFacade.EXPECT().Unix(int64(0), int64(0)).Return(time.Unix(0, 0)).Times(2)
		require.NoError(t, container.Provide(func() flamTime.Facade { return timeFacade }))

		data1 := "{\"timestamp\": \"1234-56-78 90:12:34 +0000 UTC\", \"config\": {\"field\": \"value\"}}"
		reader1 := func(b []byte) (int, error) {
			copy(b, data1)
			return len(data1), io.EOF
		}

		body1 := mocks.NewReadCloser(ctrl)
		body1.EXPECT().Read(gomock.Any()).DoAndReturn(reader1).Times(1)

		data2 := "{\"timestamp\": \"1234-56-78 90:12:35 +0000 UTC\", \"config\": {\"field\": \"value2\"}}"
		reader2 := func(b []byte) (int, error) {
			copy(b, data2)
			return len(data2), io.EOF
		}

		body2 := mocks.NewReadCloser(ctrl)
		body2.EXPECT().Read(gomock.Any()).DoAndReturn(reader2).Times(1)

		response1 := &http.Response{Body: body1}

		response2 := &http.Response{Body: body2}

		requester := mocks.NewRestRequester(ctrl)
		requester.EXPECT().Do(gomock.Any()).Return(response1, nil)
		requester.EXPECT().Do(gomock.Any()).Return(response2, nil)

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

			reloaded, e := got.(config.ObservableSource).Reload()
			assert.False(t, reloaded)
			assert.NoError(t, e)

			assert.Equal(t, "value", got.Get("field"))
		}))
	})
}
