package config

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

	"github.com/happyhippyhippo/flam"
	filesystem "github.com/happyhippyhippo/flam-filesystem"
	flamTime "github.com/happyhippyhippo/flam-time"
)

func Test_observableRestSource(t *testing.T) {
	t.Run("should ignore config without uri field", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		Defaults = flam.Bag{}
		_ = Defaults.Set(PathBoot, true)
		_ = Defaults.Set(PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   SourceDriverObservableRest,
				"parser":   "my_parser",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		assert.ErrorIs(
			t,
			NewProvider().(flam.BootableProvider).Boot(container),
			flam.ErrInvalidResourceConfig)
	})

	t.Run("should ignore config without path.config field", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		Defaults = flam.Bag{}
		_ = Defaults.Set(PathBoot, true)
		_ = Defaults.Set(PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://path/",
				"path":     flam.Bag{"timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		assert.ErrorIs(
			t,
			NewProvider().(flam.BootableProvider).Boot(container),
			flam.ErrInvalidResourceConfig)
	})

	t.Run("should ignore config without path.timestamp field", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		Defaults = flam.Bag{}
		_ = Defaults.Set(PathBoot, true)
		_ = Defaults.Set(PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://path/",
				"path":     flam.Bag{"config": "config"},
				"priority": 123,
			}})
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		assert.ErrorIs(
			t,
			NewProvider().(flam.BootableProvider).Boot(container),
			flam.ErrInvalidResourceConfig)
	})

	t.Run("should return requester generation error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		Defaults = flam.Bag{}
		_ = Defaults.Set(PathBoot, true)
		_ = Defaults.Set(PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://path/",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		expectedErr := errors.New("requester error")
		requestGenerator := NewRestRequesterGeneratorMock(ctrl)
		requestGenerator.EXPECT().Create().Return(nil, expectedErr).Times(1)
		require.NoError(t, container.Decorate(func(generator RestRequesterGenerator) RestRequesterGenerator {
			return requestGenerator
		}))

		assert.ErrorIs(
			t,
			NewProvider().(flam.BootableProvider).Boot(container),
			expectedErr)
	})

	t.Run("should return parser generation error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		Defaults = flam.Bag{}
		_ = Defaults.Set(PathBoot, true)
		_ = Defaults.Set(PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://path/",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		requester := NewRestRequesterMock(ctrl)

		requestGenerator := NewRestRequesterGeneratorMock(ctrl)
		requestGenerator.EXPECT().Create().Return(requester, nil).Times(1)
		require.NoError(t, container.Decorate(func(generator RestRequesterGenerator) RestRequesterGenerator {
			return requestGenerator
		}))

		assert.ErrorIs(
			t,
			NewProvider().(flam.BootableProvider).Boot(container),
			flam.ErrUnknownResource)
	})

	t.Run("should return request generation error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		Defaults = flam.Bag{}
		_ = Defaults.Set(PathBoot, true)
		_ = Defaults.Set(PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": ParserDriverJson,
			}})
		_ = Defaults.Set(PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      ":/uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		requester := NewRestRequesterMock(ctrl)

		requestGenerator := NewRestRequesterGeneratorMock(ctrl)
		requestGenerator.EXPECT().Create().Return(requester, nil).Times(1)
		require.NoError(t, container.Decorate(func(generator RestRequesterGenerator) RestRequesterGenerator {
			return requestGenerator
		}))

		assert.ErrorContains(
			t,
			NewProvider().(flam.BootableProvider).Boot(container),
			"missing protocol scheme")
	})

	t.Run("should return requester error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		Defaults = flam.Bag{}
		_ = Defaults.Set(PathBoot, true)
		_ = Defaults.Set(PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": ParserDriverJson,
			}})
		_ = Defaults.Set(PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		expectedErr := errors.New("requester error")
		requester := NewRestRequesterMock(ctrl)
		requester.EXPECT().Do(gomock.Any()).Return(nil, expectedErr).Times(1)

		requestGenerator := NewRestRequesterGeneratorMock(ctrl)
		requestGenerator.EXPECT().Create().Return(requester, nil).Times(1)
		require.NoError(t, container.Decorate(func(generator RestRequesterGenerator) RestRequesterGenerator {
			return requestGenerator
		}))

		assert.ErrorIs(
			t,
			NewProvider().(flam.BootableProvider).Boot(container),
			expectedErr)
	})

	t.Run("should return response read error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		Defaults = flam.Bag{}
		_ = Defaults.Set(PathBoot, true)
		_ = Defaults.Set(PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": ParserDriverJson,
			}})
		_ = Defaults.Set(PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		expectedErr := errors.New("requester error")
		body := NewReadCloserMock(ctrl)
		body.EXPECT().Read(gomock.Any()).Return(0, expectedErr).Times(1)

		response := &http.Response{Body: body}

		requester := NewRestRequesterMock(ctrl)
		requester.EXPECT().Do(gomock.Any()).Return(response, nil).Times(1)

		requestGenerator := NewRestRequesterGeneratorMock(ctrl)
		requestGenerator.EXPECT().Create().Return(requester, nil).Times(1)
		require.NoError(t, container.Decorate(func(generator RestRequesterGenerator) RestRequesterGenerator {
			return requestGenerator
		}))

		assert.ErrorIs(
			t,
			NewProvider().(flam.BootableProvider).Boot(container),
			expectedErr)
	})

	t.Run("should return response parsing error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		Defaults = flam.Bag{}
		_ = Defaults.Set(PathBoot, true)
		_ = Defaults.Set(PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": ParserDriverJson,
			}})
		_ = Defaults.Set(PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		data := "{"
		reader := func(b []byte) (int, error) {
			copy(b, data)
			return len(data), io.EOF
		}

		body := NewReadCloserMock(ctrl)
		body.EXPECT().Read(gomock.Any()).DoAndReturn(reader).Times(1)

		response := &http.Response{Body: body}

		requester := NewRestRequesterMock(ctrl)
		requester.EXPECT().Do(gomock.Any()).Return(response, nil).Times(1)

		requestGenerator := NewRestRequesterGeneratorMock(ctrl)
		requestGenerator.EXPECT().Create().Return(requester, nil).Times(1)
		require.NoError(t, container.Decorate(func(generator RestRequesterGenerator) RestRequesterGenerator {
			return requestGenerator
		}))

		assert.ErrorContains(
			t,
			NewProvider().(flam.BootableProvider).Boot(container),
			"unexpected end of JSON input")
	})

	t.Run("should return timestamp not found in response error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		Defaults = flam.Bag{}
		_ = Defaults.Set(PathBoot, true)
		_ = Defaults.Set(PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": ParserDriverJson,
			}})
		_ = Defaults.Set(PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		data := "{}"
		reader := func(b []byte) (int, error) {
			copy(b, data)
			return len(data), io.EOF
		}

		body := NewReadCloserMock(ctrl)
		body.EXPECT().Read(gomock.Any()).DoAndReturn(reader).Times(1)

		response := &http.Response{Body: body}

		requester := NewRestRequesterMock(ctrl)
		requester.EXPECT().Do(gomock.Any()).Return(response, nil).Times(1)

		requestGenerator := NewRestRequesterGeneratorMock(ctrl)
		requestGenerator.EXPECT().Create().Return(requester, nil).Times(1)
		require.NoError(t, container.Decorate(func(generator RestRequesterGenerator) RestRequesterGenerator {
			return requestGenerator
		}))

		assert.ErrorIs(
			t,
			NewProvider().(flam.BootableProvider).Boot(container),
			ErrRestTimestampNotFound)
	})

	t.Run("should return invalid timestamp in response error (type)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		Defaults = flam.Bag{}
		_ = Defaults.Set(PathBoot, true)
		_ = Defaults.Set(PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": ParserDriverJson,
			}})
		_ = Defaults.Set(PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		data := "{\"timestamp\": 123}"
		reader := func(b []byte) (int, error) {
			copy(b, data)
			return len(data), io.EOF
		}

		body := NewReadCloserMock(ctrl)
		body.EXPECT().Read(gomock.Any()).DoAndReturn(reader).Times(1)

		response := &http.Response{Body: body}

		requester := NewRestRequesterMock(ctrl)
		requester.EXPECT().Do(gomock.Any()).Return(response, nil).Times(1)

		requestGenerator := NewRestRequesterGeneratorMock(ctrl)
		requestGenerator.EXPECT().Create().Return(requester, nil).Times(1)
		require.NoError(t, container.Decorate(func(generator RestRequesterGenerator) RestRequesterGenerator {
			return requestGenerator
		}))

		assert.ErrorIs(
			t,
			NewProvider().(flam.BootableProvider).Boot(container),
			ErrRestInvalidTimestamp)
	})

	t.Run("should return invalid timestamp in response error (string parsing)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		Defaults = flam.Bag{}
		_ = Defaults.Set(PathBoot, true)
		_ = Defaults.Set(PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": ParserDriverJson,
			}})
		_ = Defaults.Set(PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		expectedErr := errors.New("parse error")
		timeFacade := NewTimeFacadeMock(ctrl)
		timeFacade.EXPECT().Now().Return(time.Now()).Times(1)
		timeFacade.EXPECT().Parse(time.RFC3339, "invalid").Return(time.Time{}, expectedErr).Times(1)
		require.NoError(t, container.Provide(func() flamTime.Facade { return timeFacade }))

		data := "{\"timestamp\": \"invalid\"}"
		reader := func(b []byte) (int, error) {
			copy(b, data)
			return len(data), io.EOF
		}

		body := NewReadCloserMock(ctrl)
		body.EXPECT().Read(gomock.Any()).DoAndReturn(reader).Times(1)

		response := &http.Response{Body: body}

		requester := NewRestRequesterMock(ctrl)
		requester.EXPECT().Do(gomock.Any()).Return(response, nil).Times(1)

		requestGenerator := NewRestRequesterGeneratorMock(ctrl)
		requestGenerator.EXPECT().Create().Return(requester, nil).Times(1)
		require.NoError(t, container.Decorate(func(generator RestRequesterGenerator) RestRequesterGenerator {
			return requestGenerator
		}))

		assert.ErrorIs(
			t,
			NewProvider().(flam.BootableProvider).Boot(container),
			expectedErr)
	})

	t.Run("should return config not found in response error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		Defaults = flam.Bag{}
		_ = Defaults.Set(PathBoot, true)
		_ = Defaults.Set(PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": ParserDriverJson,
			}})
		_ = Defaults.Set(PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		now := time.Now()
		timeFacade := NewTimeFacadeMock(ctrl)
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

		body := NewReadCloserMock(ctrl)
		body.EXPECT().Read(gomock.Any()).DoAndReturn(reader).Times(1)

		response := &http.Response{Body: body}

		requester := NewRestRequesterMock(ctrl)
		requester.EXPECT().Do(gomock.Any()).Return(response, nil).Times(1)

		requestGenerator := NewRestRequesterGeneratorMock(ctrl)
		requestGenerator.EXPECT().Create().Return(requester, nil).Times(1)
		require.NoError(t, container.Decorate(func(generator RestRequesterGenerator) RestRequesterGenerator {
			return requestGenerator
		}))

		assert.ErrorIs(
			t,
			NewProvider().(flam.BootableProvider).Boot(container),
			ErrRestConfigNotFound)
	})

	t.Run("should return invalid config in response error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		Defaults = flam.Bag{}
		_ = Defaults.Set(PathBoot, true)
		_ = Defaults.Set(PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": ParserDriverJson,
			}})
		_ = Defaults.Set(PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		now := time.Now()
		timeFacade := NewTimeFacadeMock(ctrl)
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

		body := NewReadCloserMock(ctrl)
		body.EXPECT().Read(gomock.Any()).DoAndReturn(reader).Times(1)

		response := &http.Response{Body: body}

		requester := NewRestRequesterMock(ctrl)
		requester.EXPECT().Do(gomock.Any()).Return(response, nil).Times(1)

		requestGenerator := NewRestRequesterGeneratorMock(ctrl)
		requestGenerator.EXPECT().Create().Return(requester, nil).Times(1)
		require.NoError(t, container.Decorate(func(generator RestRequesterGenerator) RestRequesterGenerator {
			return requestGenerator
		}))

		assert.ErrorIs(
			t,
			NewProvider().(flam.BootableProvider).Boot(container),
			ErrRestInvalidConfig)
	})

	t.Run("should correctly load the config", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		Defaults = flam.Bag{}
		_ = Defaults.Set(PathBoot, true)
		_ = Defaults.Set(PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": ParserDriverJson,
			}})
		_ = Defaults.Set(PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		now := time.Now()
		timeFacade := NewTimeFacadeMock(ctrl)
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

		body := NewReadCloserMock(ctrl)
		body.EXPECT().Read(gomock.Any()).DoAndReturn(reader).Times(1)

		response := &http.Response{Body: body}

		requester := NewRestRequesterMock(ctrl)
		requester.EXPECT().Do(gomock.Any()).Return(response, nil).Times(1)

		requestGenerator := NewRestRequesterGeneratorMock(ctrl)
		requestGenerator.EXPECT().Create().Return(requester, nil).Times(1)
		require.NoError(t, container.Decorate(func(generator RestRequesterGenerator) RestRequesterGenerator {
			return requestGenerator
		}))

		require.NoError(t, NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
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

		Defaults = flam.Bag{}
		_ = Defaults.Set(PathBoot, true)
		_ = Defaults.Set(PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": ParserDriverJson,
			}})
		_ = Defaults.Set(PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		now := time.Now()
		timeFacade := NewTimeFacadeMock(ctrl)
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

		body := NewReadCloserMock(ctrl)
		body.EXPECT().Read(gomock.Any()).DoAndReturn(reader).Times(1)

		response := &http.Response{Body: body}

		expectedErr := errors.New("requester error")
		requester := NewRestRequesterMock(ctrl)
		requester.EXPECT().Do(gomock.Any()).Return(response, nil)
		requester.EXPECT().Do(gomock.Any()).Return(nil, expectedErr)

		requestGenerator := NewRestRequesterGeneratorMock(ctrl)
		requestGenerator.EXPECT().Create().Return(requester, nil).Times(1)
		require.NoError(t, container.Decorate(func(generator RestRequesterGenerator) RestRequesterGenerator {
			return requestGenerator
		}))

		require.NoError(t, NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			got, e := facade.GetSource("my_source")
			require.NotNil(t, got)
			require.NoError(t, e)

			reloaded, e := got.(ObservableSource).Reload()
			assert.False(t, reloaded)
			assert.ErrorIs(t, e, expectedErr)
		}))
	})

	t.Run("should return response read error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		Defaults = flam.Bag{}
		_ = Defaults.Set(PathBoot, true)
		_ = Defaults.Set(PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": ParserDriverJson,
			}})
		_ = Defaults.Set(PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		now := time.Now()
		timeFacade := NewTimeFacadeMock(ctrl)
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

		body1 := NewReadCloserMock(ctrl)
		body1.EXPECT().Read(gomock.Any()).DoAndReturn(reader).Times(1)

		expectedErr := errors.New("reader error")
		body2 := NewReadCloserMock(ctrl)
		body2.EXPECT().Read(gomock.Any()).Return(0, expectedErr).Times(1)

		response1 := &http.Response{Body: body1}

		response2 := &http.Response{Body: body2}

		requester := NewRestRequesterMock(ctrl)
		requester.EXPECT().Do(gomock.Any()).Return(response1, nil)
		requester.EXPECT().Do(gomock.Any()).Return(response2, nil)

		requestGenerator := NewRestRequesterGeneratorMock(ctrl)
		requestGenerator.EXPECT().Create().Return(requester, nil).Times(1)
		require.NoError(t, container.Decorate(func(generator RestRequesterGenerator) RestRequesterGenerator {
			return requestGenerator
		}))

		require.NoError(t, NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			got, e := facade.GetSource("my_source")
			require.NotNil(t, got)
			require.NoError(t, e)

			reloaded, e := got.(ObservableSource).Reload()
			assert.False(t, reloaded)
			assert.ErrorIs(t, e, expectedErr)
		}))
	})

	t.Run("should return response parsing error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		Defaults = flam.Bag{}
		_ = Defaults.Set(PathBoot, true)
		_ = Defaults.Set(PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": ParserDriverJson,
			}})
		_ = Defaults.Set(PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		now := time.Now()
		timeFacade := NewTimeFacadeMock(ctrl)
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

		body1 := NewReadCloserMock(ctrl)
		body1.EXPECT().Read(gomock.Any()).DoAndReturn(reader1).Times(1)

		data2 := "{"
		reader2 := func(b []byte) (int, error) {
			copy(b, data2)
			return len(data2), io.EOF
		}

		body2 := NewReadCloserMock(ctrl)
		body2.EXPECT().Read(gomock.Any()).DoAndReturn(reader2).Times(1)

		response1 := &http.Response{Body: body1}

		response2 := &http.Response{Body: body2}

		requester := NewRestRequesterMock(ctrl)
		requester.EXPECT().Do(gomock.Any()).Return(response1, nil)
		requester.EXPECT().Do(gomock.Any()).Return(response2, nil)

		requestGenerator := NewRestRequesterGeneratorMock(ctrl)
		requestGenerator.EXPECT().Create().Return(requester, nil).Times(1)
		require.NoError(t, container.Decorate(func(generator RestRequesterGenerator) RestRequesterGenerator {
			return requestGenerator
		}))

		require.NoError(t, NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			got, e := facade.GetSource("my_source")
			require.NotNil(t, got)
			require.NoError(t, e)

			reloaded, e := got.(ObservableSource).Reload()
			assert.False(t, reloaded)
			assert.ErrorContains(t, e, "unexpected end of JSON input")
		}))
	})

	t.Run("should return timestamp not found in response error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		Defaults = flam.Bag{}
		_ = Defaults.Set(PathBoot, true)
		_ = Defaults.Set(PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": ParserDriverJson,
			}})
		_ = Defaults.Set(PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		now := time.Now()
		timeFacade := NewTimeFacadeMock(ctrl)
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

		body1 := NewReadCloserMock(ctrl)
		body1.EXPECT().Read(gomock.Any()).DoAndReturn(reader1).Times(1)

		data2 := "{}"
		reader2 := func(b []byte) (int, error) {
			copy(b, data2)
			return len(data2), io.EOF
		}

		body2 := NewReadCloserMock(ctrl)
		body2.EXPECT().Read(gomock.Any()).DoAndReturn(reader2).Times(1)

		response1 := &http.Response{Body: body1}

		response2 := &http.Response{Body: body2}

		requester := NewRestRequesterMock(ctrl)
		requester.EXPECT().Do(gomock.Any()).Return(response1, nil)
		requester.EXPECT().Do(gomock.Any()).Return(response2, nil)

		requestGenerator := NewRestRequesterGeneratorMock(ctrl)
		requestGenerator.EXPECT().Create().Return(requester, nil).Times(1)
		require.NoError(t, container.Decorate(func(generator RestRequesterGenerator) RestRequesterGenerator {
			return requestGenerator
		}))

		require.NoError(t, NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			got, e := facade.GetSource("my_source")
			require.NotNil(t, got)
			require.NoError(t, e)

			reloaded, e := got.(ObservableSource).Reload()
			assert.False(t, reloaded)
			assert.ErrorIs(t, e, ErrRestTimestampNotFound)
		}))
	})

	t.Run("should return invalid timestamp in response error (type)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		Defaults = flam.Bag{}
		_ = Defaults.Set(PathBoot, true)
		_ = Defaults.Set(PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": ParserDriverJson,
			}})
		_ = Defaults.Set(PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		now := time.Now()
		timeFacade := NewTimeFacadeMock(ctrl)
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

		body1 := NewReadCloserMock(ctrl)
		body1.EXPECT().Read(gomock.Any()).DoAndReturn(reader1).Times(1)

		data2 := "{\"timestamp\": 1234567890}"
		reader2 := func(b []byte) (int, error) {
			copy(b, data2)
			return len(data2), io.EOF
		}

		body2 := NewReadCloserMock(ctrl)
		body2.EXPECT().Read(gomock.Any()).DoAndReturn(reader2).Times(1)

		response1 := &http.Response{Body: body1}

		response2 := &http.Response{Body: body2}

		requester := NewRestRequesterMock(ctrl)
		requester.EXPECT().Do(gomock.Any()).Return(response1, nil)
		requester.EXPECT().Do(gomock.Any()).Return(response2, nil)

		requestGenerator := NewRestRequesterGeneratorMock(ctrl)
		requestGenerator.EXPECT().Create().Return(requester, nil).Times(1)
		require.NoError(t, container.Decorate(func(generator RestRequesterGenerator) RestRequesterGenerator {
			return requestGenerator
		}))

		require.NoError(t, NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			got, e := facade.GetSource("my_source")
			require.NotNil(t, got)
			require.NoError(t, e)

			reloaded, e := got.(ObservableSource).Reload()
			assert.False(t, reloaded)
			assert.ErrorIs(t, e, ErrRestInvalidTimestamp)
		}))
	})

	t.Run("should return invalid timestamp in response error (string parsing)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		Defaults = flam.Bag{}
		_ = Defaults.Set(PathBoot, true)
		_ = Defaults.Set(PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": ParserDriverJson,
			}})
		_ = Defaults.Set(PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		expectedErr := errors.New("invalid timestamp")
		now := time.Now()
		timeFacade := NewTimeFacadeMock(ctrl)
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

		body1 := NewReadCloserMock(ctrl)
		body1.EXPECT().Read(gomock.Any()).DoAndReturn(reader1).Times(1)

		data2 := "{\"timestamp\": \"invalid\"}"
		reader2 := func(b []byte) (int, error) {
			copy(b, data2)
			return len(data2), io.EOF
		}

		body2 := NewReadCloserMock(ctrl)
		body2.EXPECT().Read(gomock.Any()).DoAndReturn(reader2).Times(1)

		response1 := &http.Response{Body: body1}

		response2 := &http.Response{Body: body2}

		requester := NewRestRequesterMock(ctrl)
		requester.EXPECT().Do(gomock.Any()).Return(response1, nil)
		requester.EXPECT().Do(gomock.Any()).Return(response2, nil)

		requestGenerator := NewRestRequesterGeneratorMock(ctrl)
		requestGenerator.EXPECT().Create().Return(requester, nil).Times(1)
		require.NoError(t, container.Decorate(func(generator RestRequesterGenerator) RestRequesterGenerator {
			return requestGenerator
		}))

		require.NoError(t, NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			got, e := facade.GetSource("my_source")
			require.NotNil(t, got)
			require.NoError(t, e)

			reloaded, e := got.(ObservableSource).Reload()
			assert.False(t, reloaded)
			assert.ErrorIs(t, e, expectedErr)
		}))
	})

	t.Run("should return config not found in response error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		Defaults = flam.Bag{}
		_ = Defaults.Set(PathBoot, true)
		_ = Defaults.Set(PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": ParserDriverJson,
			}})
		_ = Defaults.Set(PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		now := time.Now()
		timeFacade := NewTimeFacadeMock(ctrl)
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

		body1 := NewReadCloserMock(ctrl)
		body1.EXPECT().Read(gomock.Any()).DoAndReturn(reader1).Times(1)

		data2 := "{\"timestamp\": \"1234-56-78 90:12:35 +0000 UTC\"}"
		reader2 := func(b []byte) (int, error) {
			copy(b, data2)
			return len(data2), io.EOF
		}

		body2 := NewReadCloserMock(ctrl)
		body2.EXPECT().Read(gomock.Any()).DoAndReturn(reader2).Times(1)

		response1 := &http.Response{Body: body1}

		response2 := &http.Response{Body: body2}

		requester := NewRestRequesterMock(ctrl)
		requester.EXPECT().Do(gomock.Any()).Return(response1, nil)
		requester.EXPECT().Do(gomock.Any()).Return(response2, nil)

		requestGenerator := NewRestRequesterGeneratorMock(ctrl)
		requestGenerator.EXPECT().Create().Return(requester, nil).Times(1)
		require.NoError(t, container.Decorate(func(generator RestRequesterGenerator) RestRequesterGenerator {
			return requestGenerator
		}))

		require.NoError(t, NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			got, e := facade.GetSource("my_source")
			require.NotNil(t, got)
			require.NoError(t, e)

			reloaded, e := got.(ObservableSource).Reload()
			assert.False(t, reloaded)
			assert.ErrorIs(t, e, ErrRestConfigNotFound)
		}))
	})

	t.Run("should return invalid config in response error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		Defaults = flam.Bag{}
		_ = Defaults.Set(PathBoot, true)
		_ = Defaults.Set(PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": ParserDriverJson,
			}})
		_ = Defaults.Set(PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		now := time.Now()
		timeFacade := NewTimeFacadeMock(ctrl)
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

		body1 := NewReadCloserMock(ctrl)
		body1.EXPECT().Read(gomock.Any()).DoAndReturn(reader1).Times(1)

		data2 := "{\"timestamp\": \"1234-56-78 90:12:35 +0000 UTC\", \"config\": 123}"
		reader2 := func(b []byte) (int, error) {
			copy(b, data2)
			return len(data2), io.EOF
		}

		body2 := NewReadCloserMock(ctrl)
		body2.EXPECT().Read(gomock.Any()).DoAndReturn(reader2).Times(1)

		response1 := &http.Response{Body: body1}

		response2 := &http.Response{Body: body2}

		requester := NewRestRequesterMock(ctrl)
		requester.EXPECT().Do(gomock.Any()).Return(response1, nil)
		requester.EXPECT().Do(gomock.Any()).Return(response2, nil)

		requestGenerator := NewRestRequesterGeneratorMock(ctrl)
		requestGenerator.EXPECT().Create().Return(requester, nil).Times(1)
		require.NoError(t, container.Decorate(func(generator RestRequesterGenerator) RestRequesterGenerator {
			return requestGenerator
		}))

		require.NoError(t, NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			got, e := facade.GetSource("my_source")
			require.NotNil(t, got)
			require.NoError(t, e)

			reloaded, e := got.(ObservableSource).Reload()
			assert.False(t, reloaded)
			assert.ErrorIs(t, e, ErrRestInvalidConfig)
		}))
	})

	t.Run("should correctly reload the source", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		Defaults = flam.Bag{}
		_ = Defaults.Set(PathBoot, true)
		_ = Defaults.Set(PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": ParserDriverJson,
			}})
		_ = Defaults.Set(PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		now := time.Now()
		timeFacade := NewTimeFacadeMock(ctrl)
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

		body1 := NewReadCloserMock(ctrl)
		body1.EXPECT().Read(gomock.Any()).DoAndReturn(reader1).Times(1)

		data2 := "{\"timestamp\": \"1234-56-78 90:12:35 +0000 UTC\", \"config\": {\"field\": \"value2\"}}"
		reader2 := func(b []byte) (int, error) {
			copy(b, data2)
			return len(data2), io.EOF
		}

		body2 := NewReadCloserMock(ctrl)
		body2.EXPECT().Read(gomock.Any()).DoAndReturn(reader2).Times(1)

		response1 := &http.Response{Body: body1}

		response2 := &http.Response{Body: body2}

		requester := NewRestRequesterMock(ctrl)
		requester.EXPECT().Do(gomock.Any()).Return(response1, nil)
		requester.EXPECT().Do(gomock.Any()).Return(response2, nil)

		requestGenerator := NewRestRequesterGeneratorMock(ctrl)
		requestGenerator.EXPECT().Create().Return(requester, nil).Times(1)
		require.NoError(t, container.Decorate(func(generator RestRequesterGenerator) RestRequesterGenerator {
			return requestGenerator
		}))

		require.NoError(t, NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			got, e := facade.GetSource("my_source")
			require.NotNil(t, got)
			require.NoError(t, e)

			reloaded, e := got.(ObservableSource).Reload()
			assert.True(t, reloaded)
			assert.NoError(t, e)

			assert.Equal(t, "value2", got.Get("field"))
		}))
	})

	t.Run("should not reload config if the timestamp is less or equals to the stored timestamp", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		Defaults = flam.Bag{}
		_ = Defaults.Set(PathBoot, true)
		_ = Defaults.Set(PathParsers, flam.Bag{
			"my_parser": flam.Bag{
				"driver": ParserDriverJson,
			}})
		_ = Defaults.Set(PathSources, flam.Bag{
			"my_source": flam.Bag{
				"driver":   SourceDriverObservableRest,
				"parser":   "my_parser",
				"uri":      "http://uri",
				"path":     flam.Bag{"config": "config", "timestamp": "timestamp"},
				"priority": 123,
			}})
		defer func() { Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		now := time.Now()
		timeFacade := NewTimeFacadeMock(ctrl)
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

		body1 := NewReadCloserMock(ctrl)
		body1.EXPECT().Read(gomock.Any()).DoAndReturn(reader1).Times(1)

		data2 := "{\"timestamp\": \"1234-56-78 90:12:35 +0000 UTC\", \"config\": {\"field\": \"value2\"}}"
		reader2 := func(b []byte) (int, error) {
			copy(b, data2)
			return len(data2), io.EOF
		}

		body2 := NewReadCloserMock(ctrl)
		body2.EXPECT().Read(gomock.Any()).DoAndReturn(reader2).Times(1)

		response1 := &http.Response{Body: body1}

		response2 := &http.Response{Body: body2}

		requester := NewRestRequesterMock(ctrl)
		requester.EXPECT().Do(gomock.Any()).Return(response1, nil)
		requester.EXPECT().Do(gomock.Any()).Return(response2, nil)

		requestGenerator := NewRestRequesterGeneratorMock(ctrl)
		requestGenerator.EXPECT().Create().Return(requester, nil).Times(1)
		require.NoError(t, container.Decorate(func(generator RestRequesterGenerator) RestRequesterGenerator {
			return requestGenerator
		}))

		require.NoError(t, NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			got, e := facade.GetSource("my_source")
			require.NotNil(t, got)
			require.NoError(t, e)

			reloaded, e := got.(ObservableSource).Reload()
			assert.False(t, reloaded)
			assert.NoError(t, e)

			assert.Equal(t, "value", got.Get("field"))
		}))
	})
}
