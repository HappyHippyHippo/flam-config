package tests

import (
	"errors"
	"os"
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

func Test_Facade_Entries(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	container := dig.New()
	require.NoError(t, config.NewProvider().Register(container))

	data := flam.Bag{"field1": "value1", "field2": "value2"}
	source := mocks.NewSource(ctrl)
	source.EXPECT().Get("", flam.Bag{}).Return(data).Times(1)

	assert.NoError(t, container.Invoke(func(facade config.Facade) {
		require.NoError(t, facade.AddSource("source", source))

		assert.ElementsMatch(t, []string{"field1", "field2"}, facade.Entries())
	}))
}

func Test_Facade_Has(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	container := dig.New()
	require.NoError(t, config.NewProvider().Register(container))

	data := flam.Bag{"field1": "value1", "field2": "value2"}
	source := mocks.NewSource(ctrl)
	source.EXPECT().Get("", flam.Bag{}).Return(data).Times(1)

	assert.NoError(t, container.Invoke(func(facade config.Facade) {
		require.NoError(t, facade.AddSource("source", source))

		assert.True(t, facade.Has("field1"))
		assert.False(t, facade.Has("invalid"))
	}))
}

func Test_Facade_Get(t *testing.T) {
	scenarios := []struct {
		name     string
		data     flam.Bag
		path     string
		def      []any
		expected any
	}{
		{
			name:     "should return nil on non-existent path",
			data:     flam.Bag{"field1": "value1", "field2": flam.Bag{"field3": "value3"}},
			path:     "invalid",
			expected: nil,
		},
		{
			name:     "should return passed default on non-existent path",
			data:     flam.Bag{"field1": "value1", "field2": flam.Bag{"field3": "value3"}},
			path:     "invalid",
			def:      []any{"default"},
			expected: "default",
		},
		{
			name:     "should return value on existent path",
			data:     flam.Bag{"field1": "value1", "field2": flam.Bag{"field3": "value3"}},
			path:     "field1",
			expected: "value1",
		},
		{
			name:     "should return value on existing inner path",
			data:     flam.Bag{"field1": "value1", "field2": flam.Bag{"field3": "value3"}},
			path:     "field2.field3",
			expected: "value3",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			container := dig.New()
			require.NoError(t, config.NewProvider().Register(container))

			source := mocks.NewSource(ctrl)
			source.EXPECT().Get("", flam.Bag{}).Return(scenario.data).Times(1)

			assert.NoError(t, container.Invoke(func(facade config.Facade) {
				require.NoError(t, facade.AddSource("source", source))

				assert.Equal(t, scenario.expected, facade.Get(scenario.path, scenario.def...))
			}))
		})
	}
}

func Test_Facade_Bool(t *testing.T) {
	scenarios := []struct {
		name     string
		data     flam.Bag
		path     string
		def      []bool
		expected bool
	}{
		{
			name:     "should return false on non-existent path",
			data:     flam.Bag{"field1": true, "field2": flam.Bag{"field3": true}},
			path:     "invalid",
			expected: false,
		},
		{
			name:     "should return passed default on non-existent path",
			data:     flam.Bag{"field1": false, "field2": flam.Bag{"field3": false}},
			path:     "invalid",
			def:      []bool{true},
			expected: true,
		},
		{
			name:     "should return value on existent path",
			data:     flam.Bag{"field1": true, "field2": flam.Bag{"field3": true}},
			path:     "field1",
			expected: true,
		},
		{
			name:     "should return value on existing inner path",
			data:     flam.Bag{"field1": true, "field2": flam.Bag{"field3": true}},
			path:     "field2.field3",
			expected: true,
		},
		{
			name:     "should return false on existing path that hold a non-boolean value",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": true}},
			path:     "field1",
			expected: false,
		},
		{
			name:     "should return default on existing path that hold a non-boolean value",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": true}},
			path:     "field1",
			def:      []bool{true},
			expected: true,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			container := dig.New()
			require.NoError(t, config.NewProvider().Register(container))

			source := mocks.NewSource(ctrl)
			source.EXPECT().Get("", flam.Bag{}).Return(scenario.data).Times(1)

			assert.NoError(t, container.Invoke(func(facade config.Facade) {
				require.NoError(t, facade.AddSource("source", source))

				assert.Equal(t, scenario.expected, facade.Bool(scenario.path, scenario.def...))
			}))
		})
	}
}

func Test_Facade_Int(t *testing.T) {
	scenarios := []struct {
		name     string
		data     flam.Bag
		path     string
		def      []int
		expected int
	}{
		{
			name:     "should return false on non-existent path",
			data:     flam.Bag{"field1": 1, "field2": flam.Bag{"field3": 1}},
			path:     "invalid",
			expected: 0,
		},
		{
			name:     "should return passed default on non-existent path",
			data:     flam.Bag{"field1": 0, "field2": flam.Bag{"field3": 0}},
			path:     "invalid",
			def:      []int{1},
			expected: 1,
		},
		{
			name:     "should return value on existent path",
			data:     flam.Bag{"field1": 1, "field2": flam.Bag{"field3": 2}},
			path:     "field1",
			expected: 1,
		},
		{
			name:     "should return value on existing inner path",
			data:     flam.Bag{"field1": 1, "field2": flam.Bag{"field3": 2}},
			path:     "field2.field3",
			expected: 2,
		},
		{
			name:     "should return 0 on existing path that hold a non-integer value",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": 1}},
			path:     "field1",
			expected: 0,
		},
		{
			name:     "should return default on existing path that hold a non-integer value",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": 1}},
			path:     "field1",
			def:      []int{1},
			expected: 1,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			container := dig.New()
			require.NoError(t, config.NewProvider().Register(container))

			source := mocks.NewSource(ctrl)
			source.EXPECT().Get("", flam.Bag{}).Return(scenario.data).Times(1)

			assert.NoError(t, container.Invoke(func(facade config.Facade) {
				require.NoError(t, facade.AddSource("source", source))

				assert.Equal(t, scenario.expected, facade.Int(scenario.path, scenario.def...))
			}))
		})
	}
}

func Test_Facade_Int8(t *testing.T) {
	scenarios := []struct {
		name     string
		data     flam.Bag
		path     string
		def      []int8
		expected int8
	}{
		{
			name:     "should return false on non-existent path",
			data:     flam.Bag{"field1": int8(1), "field2": flam.Bag{"field3": int8(1)}},
			path:     "invalid",
			expected: 0,
		},
		{
			name:     "should return passed default on non-existent path",
			data:     flam.Bag{"field1": int8(0), "field2": flam.Bag{"field3": int8(0)}},
			path:     "invalid",
			def:      []int8{1},
			expected: 1,
		},
		{
			name:     "should return value on existent path",
			data:     flam.Bag{"field1": int8(1), "field2": flam.Bag{"field3": int8(2)}},
			path:     "field1",
			expected: 1,
		},
		{
			name:     "should return value on existing inner path",
			data:     flam.Bag{"field1": int8(1), "field2": flam.Bag{"field3": int8(2)}},
			path:     "field2.field3",
			expected: 2,
		},
		{
			name:     "should return 0 on existing path that hold a non-integer value",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": 1}},
			path:     "field1",
			expected: 0,
		},
		{
			name:     "should return default on existing path that hold a non-integer value",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": 1}},
			path:     "field1",
			def:      []int8{1},
			expected: 1,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			container := dig.New()
			require.NoError(t, config.NewProvider().Register(container))

			source := mocks.NewSource(ctrl)
			source.EXPECT().Get("", flam.Bag{}).Return(scenario.data).Times(1)

			assert.NoError(t, container.Invoke(func(facade config.Facade) {
				require.NoError(t, facade.AddSource("source", source))

				assert.Equal(t, scenario.expected, facade.Int8(scenario.path, scenario.def...))
			}))
		})
	}
}

func Test_Facade_Int16(t *testing.T) {
	scenarios := []struct {
		name     string
		data     flam.Bag
		path     string
		def      []int16
		expected int16
	}{
		{
			name:     "should return false on non-existent path",
			data:     flam.Bag{"field1": int16(1), "field2": flam.Bag{"field3": int16(1)}},
			path:     "invalid",
			expected: 0,
		},
		{
			name:     "should return passed default on non-existent path",
			data:     flam.Bag{"field1": int16(0), "field2": flam.Bag{"field3": int16(0)}},
			path:     "invalid",
			def:      []int16{1},
			expected: 1,
		},
		{
			name:     "should return value on existent path",
			data:     flam.Bag{"field1": int16(1), "field2": flam.Bag{"field3": int16(2)}},
			path:     "field1",
			expected: 1,
		},
		{
			name:     "should return value on existing inner path",
			data:     flam.Bag{"field1": int16(1), "field2": flam.Bag{"field3": int16(2)}},
			path:     "field2.field3",
			expected: 2,
		},
		{
			name:     "should return 0 on existing path that hold a non-integer value",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": 1}},
			path:     "field1",
			expected: 0,
		},
		{
			name:     "should return default on existing path that hold a non-integer value",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": 1}},
			path:     "field1",
			def:      []int16{1},
			expected: 1,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			container := dig.New()
			require.NoError(t, config.NewProvider().Register(container))

			source := mocks.NewSource(ctrl)
			source.EXPECT().Get("", flam.Bag{}).Return(scenario.data).Times(1)

			assert.NoError(t, container.Invoke(func(facade config.Facade) {
				require.NoError(t, facade.AddSource("source", source))

				assert.Equal(t, scenario.expected, facade.Int16(scenario.path, scenario.def...))
			}))
		})
	}
}

func Test_Facade_Int32(t *testing.T) {
	scenarios := []struct {
		name     string
		data     flam.Bag
		path     string
		def      []int32
		expected int32
	}{
		{
			name:     "should return false on non-existent path",
			data:     flam.Bag{"field1": int32(1), "field2": flam.Bag{"field3": int32(1)}},
			path:     "invalid",
			expected: 0,
		},
		{
			name:     "should return passed default on non-existent path",
			data:     flam.Bag{"field1": int32(0), "field2": flam.Bag{"field3": int32(0)}},
			path:     "invalid",
			def:      []int32{1},
			expected: 1,
		},
		{
			name:     "should return value on existent path",
			data:     flam.Bag{"field1": int32(1), "field2": flam.Bag{"field3": int32(2)}},
			path:     "field1",
			expected: 1,
		},
		{
			name:     "should return value on existing inner path",
			data:     flam.Bag{"field1": int32(1), "field2": flam.Bag{"field3": int32(2)}},
			path:     "field2.field3",
			expected: 2,
		},
		{
			name:     "should return 0 on existing path that hold a non-integer value",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": 1}},
			path:     "field1",
			expected: 0,
		},
		{
			name:     "should return default on existing path that hold a non-integer value",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": 1}},
			path:     "field1",
			def:      []int32{1},
			expected: 1,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			container := dig.New()
			require.NoError(t, config.NewProvider().Register(container))

			source := mocks.NewSource(ctrl)
			source.EXPECT().Get("", flam.Bag{}).Return(scenario.data).Times(1)

			assert.NoError(t, container.Invoke(func(facade config.Facade) {
				require.NoError(t, facade.AddSource("source", source))

				assert.Equal(t, scenario.expected, facade.Int32(scenario.path, scenario.def...))
			}))
		})
	}
}

func Test_Facade_Int64(t *testing.T) {
	scenarios := []struct {
		name     string
		data     flam.Bag
		path     string
		def      []int64
		expected int64
	}{
		{
			name:     "should return false on non-existent path",
			data:     flam.Bag{"field1": int64(1), "field2": flam.Bag{"field3": int64(1)}},
			path:     "invalid",
			expected: 0,
		},
		{
			name:     "should return passed default on non-existent path",
			data:     flam.Bag{"field1": int64(0), "field2": flam.Bag{"field3": int64(0)}},
			path:     "invalid",
			def:      []int64{1},
			expected: 1,
		},
		{
			name:     "should return value on existent path",
			data:     flam.Bag{"field1": int64(1), "field2": flam.Bag{"field3": int64(2)}},
			path:     "field1",
			expected: 1,
		},
		{
			name:     "should return value on existing inner path",
			data:     flam.Bag{"field1": int64(1), "field2": flam.Bag{"field3": int64(2)}},
			path:     "field2.field3",
			expected: 2,
		},
		{
			name:     "should return 0 on existing path that hold a non-integer value",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": 1}},
			path:     "field1",
			expected: 0,
		},
		{
			name:     "should return default on existing path that hold a non-integer value",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": 1}},
			path:     "field1",
			def:      []int64{1},
			expected: 1,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			container := dig.New()
			require.NoError(t, config.NewProvider().Register(container))

			source := mocks.NewSource(ctrl)
			source.EXPECT().Get("", flam.Bag{}).Return(scenario.data).Times(1)

			assert.NoError(t, container.Invoke(func(facade config.Facade) {
				require.NoError(t, facade.AddSource("source", source))

				assert.Equal(t, scenario.expected, facade.Int64(scenario.path, scenario.def...))
			}))
		})
	}
}

func Test_Facade_Uint(t *testing.T) {
	scenarios := []struct {
		name     string
		data     flam.Bag
		path     string
		def      []uint
		expected uint
	}{
		{
			name:     "should return false on non-existent path",
			data:     flam.Bag{"field1": uint(1), "field2": flam.Bag{"field3": uint(1)}},
			path:     "invalid",
			expected: 0,
		},
		{
			name:     "should return passed default on non-existent path",
			data:     flam.Bag{"field1": uint(0), "field2": flam.Bag{"field3": uint(0)}},
			path:     "invalid",
			def:      []uint{1},
			expected: 1,
		},
		{
			name:     "should return value on existent path",
			data:     flam.Bag{"field1": uint(1), "field2": flam.Bag{"field3": uint(2)}},
			path:     "field1",
			expected: 1,
		},
		{
			name:     "should return value on existing inner path",
			data:     flam.Bag{"field1": uint(1), "field2": flam.Bag{"field3": uint(2)}},
			path:     "field2.field3",
			expected: 2,
		},
		{
			name:     "should return 0 on existing path that hold a non-integer value",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": 1}},
			path:     "field1",
			expected: 0,
		},
		{
			name:     "should return default on existing path that hold a non-integer value",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": 1}},
			path:     "field1",
			def:      []uint{1},
			expected: 1,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			container := dig.New()
			require.NoError(t, config.NewProvider().Register(container))

			source := mocks.NewSource(ctrl)
			source.EXPECT().Get("", flam.Bag{}).Return(scenario.data).Times(1)

			assert.NoError(t, container.Invoke(func(facade config.Facade) {
				require.NoError(t, facade.AddSource("source", source))

				assert.Equal(t, scenario.expected, facade.Uint(scenario.path, scenario.def...))
			}))
		})
	}
}

func Test_Facade_Uint8(t *testing.T) {
	scenarios := []struct {
		name     string
		data     flam.Bag
		path     string
		def      []uint8
		expected uint8
	}{
		{
			name:     "should return false on non-existent path",
			data:     flam.Bag{"field1": uint8(1), "field2": flam.Bag{"field3": uint8(1)}},
			path:     "invalid",
			expected: 0,
		},
		{
			name:     "should return passed default on non-existent path",
			data:     flam.Bag{"field1": uint8(0), "field2": flam.Bag{"field3": uint8(0)}},
			path:     "invalid",
			def:      []uint8{1},
			expected: 1,
		},
		{
			name:     "should return value on existent path",
			data:     flam.Bag{"field1": uint8(1), "field2": flam.Bag{"field3": uint8(2)}},
			path:     "field1",
			expected: 1,
		},
		{
			name:     "should return value on existing inner path",
			data:     flam.Bag{"field1": uint8(1), "field2": flam.Bag{"field3": uint8(2)}},
			path:     "field2.field3",
			expected: 2,
		},
		{
			name:     "should return 0 on existing path that hold a non-integer value",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": 1}},
			path:     "field1",
			expected: 0,
		},
		{
			name:     "should return default on existing path that hold a non-integer value",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": 1}},
			path:     "field1",
			def:      []uint8{1},
			expected: 1,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			container := dig.New()
			require.NoError(t, config.NewProvider().Register(container))

			source := mocks.NewSource(ctrl)
			source.EXPECT().Get("", flam.Bag{}).Return(scenario.data).Times(1)

			assert.NoError(t, container.Invoke(func(facade config.Facade) {
				require.NoError(t, facade.AddSource("source", source))

				assert.Equal(t, scenario.expected, facade.Uint8(scenario.path, scenario.def...))
			}))
		})
	}
}

func Test_Facade_Uint16(t *testing.T) {
	scenarios := []struct {
		name     string
		data     flam.Bag
		path     string
		def      []uint16
		expected uint16
	}{
		{
			name:     "should return false on non-existent path",
			data:     flam.Bag{"field1": uint16(1), "field2": flam.Bag{"field3": uint16(1)}},
			path:     "invalid",
			expected: 0,
		},
		{
			name:     "should return passed default on non-existent path",
			data:     flam.Bag{"field1": uint16(0), "field2": flam.Bag{"field3": uint16(0)}},
			path:     "invalid",
			def:      []uint16{1},
			expected: 1,
		},
		{
			name:     "should return value on existent path",
			data:     flam.Bag{"field1": uint16(1), "field2": flam.Bag{"field3": uint16(2)}},
			path:     "field1",
			expected: 1,
		},
		{
			name:     "should return value on existing inner path",
			data:     flam.Bag{"field1": uint16(1), "field2": flam.Bag{"field3": uint16(2)}},
			path:     "field2.field3",
			expected: 2,
		},
		{
			name:     "should return 0 on existing path that hold a non-integer value",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": 1}},
			path:     "field1",
			expected: 0,
		},
		{
			name:     "should return default on existing path that hold a non-integer value",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": 1}},
			path:     "field1",
			def:      []uint16{1},
			expected: 1,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			container := dig.New()
			require.NoError(t, config.NewProvider().Register(container))

			source := mocks.NewSource(ctrl)
			source.EXPECT().Get("", flam.Bag{}).Return(scenario.data).Times(1)

			assert.NoError(t, container.Invoke(func(facade config.Facade) {
				require.NoError(t, facade.AddSource("source", source))

				assert.Equal(t, scenario.expected, facade.Uint16(scenario.path, scenario.def...))
			}))
		})
	}
}

func Test_Facade_Uint32(t *testing.T) {
	scenarios := []struct {
		name     string
		data     flam.Bag
		path     string
		def      []uint32
		expected uint32
	}{
		{
			name:     "should return false on non-existent path",
			data:     flam.Bag{"field1": uint32(1), "field2": flam.Bag{"field3": uint32(1)}},
			path:     "invalid",
			expected: 0,
		},
		{
			name:     "should return passed default on non-existent path",
			data:     flam.Bag{"field1": uint32(0), "field2": flam.Bag{"field3": uint32(0)}},
			path:     "invalid",
			def:      []uint32{1},
			expected: 1,
		},
		{
			name:     "should return value on existent path",
			data:     flam.Bag{"field1": uint32(1), "field2": flam.Bag{"field3": uint32(2)}},
			path:     "field1",
			expected: 1,
		},
		{
			name:     "should return value on existing inner path",
			data:     flam.Bag{"field1": uint32(1), "field2": flam.Bag{"field3": uint32(2)}},
			path:     "field2.field3",
			expected: 2,
		},
		{
			name:     "should return 0 on existing path that hold a non-integer value",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": 1}},
			path:     "field1",
			expected: 0,
		},
		{
			name:     "should return default on existing path that hold a non-integer value",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": 1}},
			path:     "field1",
			def:      []uint32{1},
			expected: 1,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			container := dig.New()
			require.NoError(t, config.NewProvider().Register(container))

			source := mocks.NewSource(ctrl)
			source.EXPECT().Get("", flam.Bag{}).Return(scenario.data).Times(1)

			assert.NoError(t, container.Invoke(func(facade config.Facade) {
				require.NoError(t, facade.AddSource("source", source))

				assert.Equal(t, scenario.expected, facade.Uint32(scenario.path, scenario.def...))
			}))
		})
	}
}

func Test_Facade_Uint64(t *testing.T) {
	scenarios := []struct {
		name     string
		data     flam.Bag
		path     string
		def      []uint64
		expected uint64
	}{
		{
			name:     "should return false on non-existent path",
			data:     flam.Bag{"field1": uint64(1), "field2": flam.Bag{"field3": uint64(1)}},
			path:     "invalid",
			expected: 0,
		},
		{
			name:     "should return passed default on non-existent path",
			data:     flam.Bag{"field1": uint64(0), "field2": flam.Bag{"field3": uint64(0)}},
			path:     "invalid",
			def:      []uint64{1},
			expected: 1,
		},
		{
			name:     "should return value on existent path",
			data:     flam.Bag{"field1": uint64(1), "field2": flam.Bag{"field3": uint64(2)}},
			path:     "field1",
			expected: 1,
		},
		{
			name:     "should return value on existing inner path",
			data:     flam.Bag{"field1": uint64(1), "field2": flam.Bag{"field3": uint64(2)}},
			path:     "field2.field3",
			expected: 2,
		},
		{
			name:     "should return 0 on existing path that hold a non-integer value",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": 1}},
			path:     "field1",
			expected: 0,
		},
		{
			name:     "should return default on existing path that hold a non-integer value",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": 1}},
			path:     "field1",
			def:      []uint64{1},
			expected: 1,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			container := dig.New()
			require.NoError(t, config.NewProvider().Register(container))

			source := mocks.NewSource(ctrl)
			source.EXPECT().Get("", flam.Bag{}).Return(scenario.data).Times(1)

			assert.NoError(t, container.Invoke(func(facade config.Facade) {
				require.NoError(t, facade.AddSource("source", source))

				assert.Equal(t, scenario.expected, facade.Uint64(scenario.path, scenario.def...))
			}))
		})
	}
}

func Test_Facade_Float32(t *testing.T) {
	scenarios := []struct {
		name     string
		data     flam.Bag
		path     string
		def      []float32
		expected float32
	}{
		{
			name:     "should return false on non-existent path",
			data:     flam.Bag{"field1": float32(1), "field2": flam.Bag{"field3": float32(1)}},
			path:     "invalid",
			expected: 0,
		},
		{
			name:     "should return passed default on non-existent path",
			data:     flam.Bag{"field1": float32(0), "field2": flam.Bag{"field3": float32(0)}},
			path:     "invalid",
			def:      []float32{1},
			expected: 1,
		},
		{
			name:     "should return value on existent path",
			data:     flam.Bag{"field1": float32(1), "field2": flam.Bag{"field3": float32(2)}},
			path:     "field1",
			expected: 1,
		},
		{
			name:     "should return value on existing inner path",
			data:     flam.Bag{"field1": float32(1), "field2": flam.Bag{"field3": float32(2)}},
			path:     "field2.field3",
			expected: 2,
		},
		{
			name:     "should return 0 on existing path that hold a non-integer value",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": 1}},
			path:     "field1",
			expected: 0,
		},
		{
			name:     "should return default on existing path that hold a non-integer value",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": 1}},
			path:     "field1",
			def:      []float32{1},
			expected: 1,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			container := dig.New()
			require.NoError(t, config.NewProvider().Register(container))

			source := mocks.NewSource(ctrl)
			source.EXPECT().Get("", flam.Bag{}).Return(scenario.data).Times(1)

			assert.NoError(t, container.Invoke(func(facade config.Facade) {
				require.NoError(t, facade.AddSource("source", source))

				assert.Equal(t, scenario.expected, facade.Float32(scenario.path, scenario.def...))
			}))
		})
	}
}

func Test_Facade_Float64(t *testing.T) {
	scenarios := []struct {
		name     string
		data     flam.Bag
		path     string
		def      []float64
		expected float64
	}{
		{
			name:     "should return false on non-existent path",
			data:     flam.Bag{"field1": float64(1), "field2": flam.Bag{"field3": float64(1)}},
			path:     "invalid",
			expected: 0,
		},
		{
			name:     "should return passed default on non-existent path",
			data:     flam.Bag{"field1": float64(0), "field2": flam.Bag{"field3": float64(0)}},
			path:     "invalid",
			def:      []float64{1},
			expected: 1,
		},
		{
			name:     "should return value on existent path",
			data:     flam.Bag{"field1": float64(1), "field2": flam.Bag{"field3": float64(2)}},
			path:     "field1",
			expected: 1,
		},
		{
			name:     "should return value on existing inner path",
			data:     flam.Bag{"field1": float64(1), "field2": flam.Bag{"field3": float64(2)}},
			path:     "field2.field3",
			expected: 2,
		},
		{
			name:     "should return 0 on existing path that hold a non-integer value",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": 1}},
			path:     "field1",
			expected: 0,
		},
		{
			name:     "should return default on existing path that hold a non-integer value",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": 1}},
			path:     "field1",
			def:      []float64{1},
			expected: 1,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			container := dig.New()
			require.NoError(t, config.NewProvider().Register(container))

			source := mocks.NewSource(ctrl)
			source.EXPECT().Get("", flam.Bag{}).Return(scenario.data).Times(1)

			assert.NoError(t, container.Invoke(func(facade config.Facade) {
				require.NoError(t, facade.AddSource("source", source))

				assert.Equal(t, scenario.expected, facade.Float64(scenario.path, scenario.def...))
			}))
		})
	}
}

func Test_Facade_String(t *testing.T) {
	scenarios := []struct {
		name     string
		data     flam.Bag
		path     string
		def      []string
		expected string
	}{
		{
			name:     "should return empty string on non-existent path",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": "string"}},
			path:     "invalid",
			expected: "",
		},
		{
			name:     "should return passed default on non-existent path",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": "string"}},
			path:     "invalid",
			def:      []string{"1"},
			expected: "1",
		},
		{
			name:     "should return value on existent path",
			data:     flam.Bag{"field1": "1", "field2": flam.Bag{"field3": "2"}},
			path:     "field1",
			expected: "1",
		},
		{
			name:     "should return value on existing inner path",
			data:     flam.Bag{"field1": "1", "field2": flam.Bag{"field3": "2"}},
			path:     "field2.field3",
			expected: "2",
		},
		{
			name:     "should return empty string on existing path that hold a non-integer value",
			data:     flam.Bag{"field1": 1, "field2": flam.Bag{"field3": "2"}},
			path:     "field1",
			expected: "",
		},
		{
			name:     "should return default on existing path that hold a non-integer value",
			data:     flam.Bag{"field1": 1, "field2": flam.Bag{"field3": "2"}},
			path:     "field1",
			def:      []string{"1"},
			expected: "1",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			container := dig.New()
			require.NoError(t, config.NewProvider().Register(container))

			source := mocks.NewSource(ctrl)
			source.EXPECT().Get("", flam.Bag{}).Return(scenario.data).Times(1)

			assert.NoError(t, container.Invoke(func(facade config.Facade) {
				require.NoError(t, facade.AddSource("source", source))

				assert.Equal(t, scenario.expected, facade.String(scenario.path, scenario.def...))
			}))
		})
	}
}

func Test_Facade_StringMap(t *testing.T) {
	scenarios := []struct {
		name     string
		data     flam.Bag
		path     string
		def      []map[string]any
		expected map[string]any
	}{
		{
			name:     "should return empty string map on non-existent path",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": "string"}},
			path:     "invalid",
			expected: map[string]any(nil),
		},
		{
			name:     "should return passed default on non-existent path",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": "string"}},
			path:     "invalid",
			def:      []map[string]any{{"1": "1"}},
			expected: map[string]any{"1": "1"},
		},
		{
			name:     "should return value on existent path",
			data:     flam.Bag{"field1": map[string]any{"1": "1"}, "field2": flam.Bag{"field3": "2"}},
			path:     "field1",
			expected: map[string]any{"1": "1"},
		},
		{
			name:     "should return value on existing inner path",
			data:     flam.Bag{"field1": map[string]any{"1": "1"}, "field2": flam.Bag{"field3": map[string]any{"2": "2"}}},
			path:     "field2.field3",
			expected: map[string]any{"2": "2"},
		},
		{
			name:     "should return empty string map on existing path that hold a non-string-map value",
			data:     flam.Bag{"field1": 1, "field2": flam.Bag{"field3": "2"}},
			path:     "field1",
			expected: map[string]any(nil),
		},
		{
			name:     "should return default on existing path that hold a non-string-map value",
			data:     flam.Bag{"field1": 1, "field2": flam.Bag{"field3": "2"}},
			path:     "field1",
			def:      []map[string]any{{"1": "1"}},
			expected: map[string]any{"1": "1"},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			container := dig.New()
			require.NoError(t, config.NewProvider().Register(container))

			source := mocks.NewSource(ctrl)
			source.EXPECT().Get("", flam.Bag{}).Return(scenario.data).Times(1)

			assert.NoError(t, container.Invoke(func(facade config.Facade) {
				require.NoError(t, facade.AddSource("source", source))

				assert.Equal(t, scenario.expected, facade.StringMap(scenario.path, scenario.def...))
			}))
		})
	}
}

func Test_Facade_StringMapString(t *testing.T) {
	scenarios := []struct {
		name     string
		data     flam.Bag
		path     string
		def      []map[string]string
		expected map[string]string
	}{
		{
			name:     "should return empty string map on non-existent path",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": "string"}},
			path:     "invalid",
			expected: map[string]string(nil),
		},
		{
			name:     "should return passed default on non-existent path",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": "string"}},
			path:     "invalid",
			def:      []map[string]string{{"1": "1"}},
			expected: map[string]string{"1": "1"},
		},
		{
			name:     "should return value on existent path",
			data:     flam.Bag{"field1": map[string]string{"1": "1"}, "field2": flam.Bag{"field3": "2"}},
			path:     "field1",
			expected: map[string]string{"1": "1"},
		},
		{
			name:     "should return value on existing inner path",
			data:     flam.Bag{"field1": map[string]string{"1": "1"}, "field2": flam.Bag{"field3": map[string]string{"2": "2"}}},
			path:     "field2.field3",
			expected: map[string]string{"2": "2"},
		},
		{
			name:     "should return empty string map on existing path that hold a non-string-map value",
			data:     flam.Bag{"field1": 1, "field2": flam.Bag{"field3": "2"}},
			path:     "field1",
			expected: map[string]string(nil),
		},
		{
			name:     "should return default on existing path that hold a non-string-map value",
			data:     flam.Bag{"field1": 1, "field2": flam.Bag{"field3": "2"}},
			path:     "field1",
			def:      []map[string]string{{"1": "1"}},
			expected: map[string]string{"1": "1"},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			container := dig.New()
			require.NoError(t, config.NewProvider().Register(container))

			source := mocks.NewSource(ctrl)
			source.EXPECT().Get("", flam.Bag{}).Return(scenario.data).Times(1)

			assert.NoError(t, container.Invoke(func(facade config.Facade) {
				require.NoError(t, facade.AddSource("source", source))

				assert.Equal(t, scenario.expected, facade.StringMapString(scenario.path, scenario.def...))
			}))
		})
	}
}

func Test_Facade_Slice(t *testing.T) {
	scenarios := []struct {
		name     string
		data     flam.Bag
		path     string
		def      [][]any
		expected []any
	}{
		{
			name:     "should return empty slice on non-existent path",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": "string"}},
			path:     "invalid",
			expected: []any(nil),
		},
		{
			name:     "should return passed default on non-existent path",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": "string"}},
			path:     "invalid",
			def:      [][]any{{"1"}},
			expected: []any{"1"},
		},
		{
			name:     "should return value on existent path",
			data:     flam.Bag{"field1": []any{"1"}, "field2": flam.Bag{"field3": "2"}},
			path:     "field1",
			expected: []any{"1"},
		},
		{
			name:     "should return value on existing inner path",
			data:     flam.Bag{"field1": []any{"1"}, "field2": flam.Bag{"field3": []any{"2"}}},
			path:     "field2.field3",
			expected: []any{"2"},
		},
		{
			name:     "should return empty slice on existing path that hold a non-string-map value",
			data:     flam.Bag{"field1": 1, "field2": flam.Bag{"field3": "2"}},
			path:     "field1",
			expected: []any(nil),
		},
		{
			name:     "should return default on existing path that hold a non-string-map value",
			data:     flam.Bag{"field1": 1, "field2": flam.Bag{"field3": "2"}},
			path:     "field1",
			def:      [][]any{{"1"}},
			expected: []any{"1"},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			container := dig.New()
			require.NoError(t, config.NewProvider().Register(container))

			source := mocks.NewSource(ctrl)
			source.EXPECT().Get("", flam.Bag{}).Return(scenario.data).Times(1)

			assert.NoError(t, container.Invoke(func(facade config.Facade) {
				require.NoError(t, facade.AddSource("source", source))

				assert.Equal(t, scenario.expected, facade.Slice(scenario.path, scenario.def...))
			}))
		})
	}
}

func Test_Facade_StringSlice(t *testing.T) {
	scenarios := []struct {
		name     string
		data     flam.Bag
		path     string
		def      [][]string
		expected []string
	}{
		{
			name:     "should return empty slice on non-existent path",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": "string"}},
			path:     "invalid",
			expected: []string(nil),
		},
		{
			name:     "should return passed default on non-existent path",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": "string"}},
			path:     "invalid",
			def:      [][]string{{"1"}},
			expected: []string{"1"},
		},
		{
			name:     "should return value on existent path",
			data:     flam.Bag{"field1": []string{"1"}, "field2": flam.Bag{"field3": "2"}},
			path:     "field1",
			expected: []string{"1"},
		},
		{
			name:     "should return value on existing inner path",
			data:     flam.Bag{"field1": []string{"1"}, "field2": flam.Bag{"field3": []string{"2"}}},
			path:     "field2.field3",
			expected: []string{"2"},
		},
		{
			name:     "should return empty slice on existing path that hold a non-string-map value",
			data:     flam.Bag{"field1": 1, "field2": flam.Bag{"field3": "2"}},
			path:     "field1",
			expected: []string(nil),
		},
		{
			name:     "should return default on existing path that hold a non-string-map value",
			data:     flam.Bag{"field1": 1, "field2": flam.Bag{"field3": "2"}},
			path:     "field1",
			def:      [][]string{{"1"}},
			expected: []string{"1"},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			container := dig.New()
			require.NoError(t, config.NewProvider().Register(container))

			source := mocks.NewSource(ctrl)
			source.EXPECT().Get("", flam.Bag{}).Return(scenario.data).Times(1)

			assert.NoError(t, container.Invoke(func(facade config.Facade) {
				require.NoError(t, facade.AddSource("source", source))

				assert.Equal(t, scenario.expected, facade.StringSlice(scenario.path, scenario.def...))
			}))
		})
	}
}

func Test_Facade_Duration(t *testing.T) {
	scenarios := []struct {
		name     string
		data     flam.Bag
		path     string
		def      []time.Duration
		expected time.Duration
	}{
		{
			name:     "should return zero duration on non-existent path",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": "string"}},
			path:     "invalid",
			expected: time.Duration(0),
		},
		{
			name:     "should return passed default on non-existent path",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": "string"}},
			path:     "invalid",
			def:      []time.Duration{1},
			expected: time.Duration(1),
		},
		{
			name:     "should return value on existent path",
			data:     flam.Bag{"field1": time.Duration(1), "field2": flam.Bag{"field3": "2"}},
			path:     "field1",
			expected: time.Duration(1),
		},
		{
			name:     "should return value on existent path (int conversion)",
			data:     flam.Bag{"field1": 1000, "field2": flam.Bag{"field3": "2"}},
			path:     "field1",
			expected: time.Duration(1000) * time.Millisecond,
		},
		{
			name:     "should return value on existent path (int64 conversion)",
			data:     flam.Bag{"field1": int64(1000), "field2": flam.Bag{"field3": "2"}},
			path:     "field1",
			expected: time.Duration(1000) * time.Millisecond,
		},
		{
			name:     "should return value on existing inner path",
			data:     flam.Bag{"field1": time.Duration(1), "field2": flam.Bag{"field3": time.Duration(2)}},
			path:     "field2.field3",
			expected: time.Duration(2),
		},
		{
			name:     "should return zero duration on existing path that hold a non-duration value",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": "2"}},
			path:     "field1",
			expected: time.Duration(0),
		},
		{
			name:     "should return default on existing path that hold a non-duration value",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": "2"}},
			path:     "field1",
			def:      []time.Duration{1},
			expected: time.Duration(1),
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			container := dig.New()
			require.NoError(t, config.NewProvider().Register(container))

			source := mocks.NewSource(ctrl)
			source.EXPECT().Get("", flam.Bag{}).Return(scenario.data).Times(1)

			assert.NoError(t, container.Invoke(func(facade config.Facade) {
				require.NoError(t, facade.AddSource("source", source))

				assert.Equal(t, scenario.expected, facade.Duration(scenario.path, scenario.def...))
			}))
		})
	}
}

func Test_Facade_Bag(t *testing.T) {
	scenarios := []struct {
		name     string
		data     flam.Bag
		path     string
		def      []flam.Bag
		expected flam.Bag
	}{
		{
			name:     "should return empty bag on non-existent path",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": "string"}},
			path:     "invalid",
			expected: flam.Bag(nil),
		},
		{
			name:     "should return passed default on non-existent path",
			data:     flam.Bag{"field1": "string", "field2": flam.Bag{"field3": "string"}},
			path:     "invalid",
			def:      []flam.Bag{{"1": "1"}},
			expected: flam.Bag{"1": "1"},
		},
		{
			name:     "should return value on existent path",
			data:     flam.Bag{"field1": flam.Bag{"1": "1"}, "field2": flam.Bag{"field3": "2"}},
			path:     "field1",
			expected: flam.Bag{"1": "1"},
		},
		{
			name:     "should return value on existing inner path",
			data:     flam.Bag{"field1": flam.Bag{"1": "1"}, "field2": flam.Bag{"field3": flam.Bag{"2": "2"}}},
			path:     "field2.field3",
			expected: flam.Bag{"2": "2"},
		},
		{
			name:     "should return empty bag on existing path that hold a non-string-map value",
			data:     flam.Bag{"field1": 1, "field2": flam.Bag{"field3": "2"}},
			path:     "field1",
			expected: flam.Bag(nil),
		},
		{
			name:     "should return default on existing path that hold a non-string-map value",
			data:     flam.Bag{"field1": 1, "field2": flam.Bag{"field3": "2"}},
			path:     "field1",
			def:      []flam.Bag{{"1": "1"}},
			expected: flam.Bag{"1": "1"},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			container := dig.New()
			require.NoError(t, config.NewProvider().Register(container))

			source := mocks.NewSource(ctrl)
			source.EXPECT().Get("", flam.Bag{}).Return(scenario.data).Times(1)

			assert.NoError(t, container.Invoke(func(facade config.Facade) {
				require.NoError(t, facade.AddSource("source", source))

				assert.Equal(t, scenario.expected, facade.Bag(scenario.path, scenario.def...))
			}))
		})
	}
}

func Test_Facade_Set(t *testing.T) {
	t.Run("should return error if passed an invalid path to save", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			assert.ErrorIs(t, facade.Set("", flam.Bag{}), flam.ErrBagInvalidPath)
		}))
	})

	t.Run("should correctly save a value", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			require.NoError(t, facade.Set("field1", "value1"))
			assert.Equal(t, "value1", facade.Get("field1"))
		}))
	})
}

func Test_Facade_Populate(t *testing.T) {
	type simpleStruct struct {
		Field int
	}

	type complexStruct struct {
		Name     string `mapstructure:"name"`
		Value    int    `mapstructure:"value"`
		Nested   simpleStruct
		Children []string
	}

	t.Run("without path", func(t *testing.T) {
		scenarios := []struct {
			test        string
			data        flam.Bag
			target      any
			expected    any
			expectedErr error
		}{
			{
				test:     "should populate a struct with a simple scalar value",
				data:     flam.Bag{"field": 123},
				target:   &simpleStruct{},
				expected: &simpleStruct{Field: 123},
			},
			{
				test: "should populate a complex struct with tags",
				data: flam.Bag{
					"name":  "test_name",
					"value": 999,
					"Nested": flam.Bag{
						"Field": 789,
					},
					"Children": []any{"child1", "child2"},
				},
				target: &complexStruct{},
				expected: &complexStruct{
					Name:     "test_name",
					Value:    999,
					Nested:   simpleStruct{Field: 789},
					Children: []string{"child1", "child2"},
				},
			},
		}

		for _, scenario := range scenarios {
			t.Run(scenario.test, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				container := dig.New()
				require.NoError(t, config.NewProvider().Register(container))

				source := mocks.NewSource(ctrl)
				source.EXPECT().Get("", flam.Bag{}).Return(scenario.data).Times(1)

				assert.NoError(t, container.Invoke(func(facade config.Facade) {
					require.NoError(t, facade.AddSource("source", source))

					e := facade.Populate(scenario.target)

					if scenario.expectedErr != nil {
						assert.ErrorIs(t, e, scenario.expectedErr)
						return
					}

					assert.NoError(t, e)
					assert.Equal(t, scenario.expected, scenario.target)
				}))
			})
		}
	})

	t.Run("with path", func(t *testing.T) {
		scenarios := []struct {
			test        string
			data        flam.Bag
			path        string
			target      any
			expected    any
			expectedErr error
		}{
			{
				test:     "should populate a struct from a nested path",
				data:     flam.Bag{"config": flam.Bag{"field": 456}},
				path:     "config",
				target:   &simpleStruct{},
				expected: &simpleStruct{Field: 456},
			},
			{
				test:        "should return an error for an invalid path",
				data:        flam.Bag{"config": flam.Bag{"field": 456}},
				path:        "invalid.path",
				target:      &simpleStruct{},
				expectedErr: flam.ErrBagInvalidPath,
			},
			{
				test: "should populate a complex struct from a nested path",
				data: flam.Bag{
					"data": flam.Bag{
						"name":  "nested_name",
						"value": 111,
						"Nested": flam.Bag{
							"Field": 222,
						},
						"Children": []any{"c1"},
					},
				},
				path:   "data",
				target: &complexStruct{},
				expected: &complexStruct{
					Name:     "nested_name",
					Value:    111,
					Nested:   simpleStruct{Field: 222},
					Children: []string{"c1"},
				},
			},
		}

		for _, scenario := range scenarios {
			t.Run(scenario.test, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				container := dig.New()
				require.NoError(t, config.NewProvider().Register(container))

				source := mocks.NewSource(ctrl)
				source.EXPECT().Get("", flam.Bag{}).Return(scenario.data).Times(1)

				assert.NoError(t, container.Invoke(func(facade config.Facade) {
					require.NoError(t, facade.AddSource("source", source))

					e := facade.Populate(scenario.target, scenario.path)

					if scenario.expectedErr != nil {
						assert.ErrorIs(t, e, scenario.expectedErr)
						return
					}

					assert.NoError(t, e)
					assert.Equal(t, scenario.expected, scenario.target)
				}))
			})
		}
	})
}

func Test_Facade_HasParser(t *testing.T) {
	t.Run("should return false on unknown parser", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			assert.False(t, facade.HasParser("unknown"))
		}))
	})

	t.Run("should return true on known parser in config", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		cfg := flam.Bag{"parser": flam.Bag{}}
		factoryConfig := mocks.NewFactoryConfig(ctrl)
		factoryConfig.EXPECT().Get(config.PathParsers).Return(cfg).Times(1)
		require.NoError(t, container.Decorate(func(flam.FactoryConfig) flam.FactoryConfig {
			return factoryConfig
		}))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			assert.True(t, facade.HasParser("parser"))
		}))
	})

	t.Run("should return true on an added parser", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		cfg := flam.Bag{}
		factoryConfig := mocks.NewFactoryConfig(ctrl)
		factoryConfig.EXPECT().Get(config.PathParsers).Return(cfg).Times(1)
		require.NoError(t, container.Decorate(func(flam.FactoryConfig) flam.FactoryConfig {
			return factoryConfig
		}))

		parser := mocks.NewParser(ctrl)

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			require.NoError(t, facade.AddParser("parser", parser))

			assert.True(t, facade.HasParser("parser"))
		}))
	})
}

func Test_Facade_ListParsers(t *testing.T) {
	t.Run("should return an empty list if no parsers has been registered", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		cfg := flam.Bag{}
		factoryConfig := mocks.NewFactoryConfig(ctrl)
		factoryConfig.EXPECT().Get(config.PathParsers).Return(cfg).Times(1)
		require.NoError(t, container.Decorate(func(flam.FactoryConfig) flam.FactoryConfig {
			return factoryConfig
		}))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			assert.Empty(t, facade.ListParsers())
		}))
	})

	t.Run("should return a list of registered parsers (added)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		cfg := flam.Bag{}
		factoryConfig := mocks.NewFactoryConfig(ctrl)
		factoryConfig.EXPECT().Get(config.PathParsers).Return(cfg).Times(3)
		require.NoError(t, container.Decorate(func(flam.FactoryConfig) flam.FactoryConfig {
			return factoryConfig
		}))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			require.NoError(t, facade.AddParser("parser1", mocks.NewParser(ctrl)))
			require.NoError(t, facade.AddParser("parser2", mocks.NewParser(ctrl)))

			assert.ElementsMatch(t, []string{"parser1", "parser2"}, facade.ListParsers())
		}))
	})

	t.Run("should return a list of registered parsers (in config)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		cfg := flam.Bag{"parser1": flam.Bag{}, "parser2": flam.Bag{}}
		factoryConfig := mocks.NewFactoryConfig(ctrl)
		factoryConfig.EXPECT().Get(config.PathParsers).Return(cfg).Times(1)
		require.NoError(t, container.Decorate(func(flam.FactoryConfig) flam.FactoryConfig {
			return factoryConfig
		}))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			assert.ElementsMatch(t, []string{"parser1", "parser2"}, facade.ListParsers())
		}))
	})

	t.Run("should return a list of registered parsers (in config and added)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		cfg := flam.Bag{"parser1": flam.Bag{}, "parser3": flam.Bag{}}
		factoryConfig := mocks.NewFactoryConfig(ctrl)
		factoryConfig.EXPECT().Get(config.PathParsers).Return(cfg).Times(2)
		require.NoError(t, container.Decorate(func(flam.FactoryConfig) flam.FactoryConfig {
			return factoryConfig
		}))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			require.NoError(t, facade.AddParser("parser2", mocks.NewParser(ctrl)))

			assert.ElementsMatch(t, []string{"parser1", "parser2", "parser3"}, facade.ListParsers())
		}))
	})
}

func Test_Facade_GetParser(t *testing.T) {
	t.Run("should return error on unknown parser", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		cfg := flam.Bag{}
		factoryConfig := mocks.NewFactoryConfig(ctrl)
		factoryConfig.EXPECT().Get(config.PathParsers).Return(cfg).Times(1)
		require.NoError(t, container.Decorate(func(flam.FactoryConfig) flam.FactoryConfig {
			return factoryConfig
		}))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			got, e := facade.GetParser("unknown")
			assert.Nil(t, got)
			assert.ErrorIs(t, e, flam.ErrUnknownResource)
		}))
	})

	t.Run("should return error when unable to generate the parser instance", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		cfg := flam.Bag{"parser": flam.Bag{"driver": "mock"}}
		factoryConfig := mocks.NewFactoryConfig(ctrl)
		factoryConfig.EXPECT().Get(config.PathParsers).Return(cfg).Times(1)
		require.NoError(t, container.Decorate(func(flam.FactoryConfig) flam.FactoryConfig {
			return factoryConfig
		}))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			got, e := facade.GetParser("parser")
			assert.Nil(t, got)
			assert.ErrorIs(t, e, flam.ErrInvalidResourceConfig)
		}))
	})

	t.Run("should return error originated when generating the parser instance", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		cfg := flam.Bag{"parser": flam.Bag{"driver": "mock"}}
		factoryConfig := mocks.NewFactoryConfig(ctrl)
		factoryConfig.EXPECT().Get(config.PathParsers).Return(cfg).Times(1)
		require.NoError(t, container.Decorate(func(flam.FactoryConfig) flam.FactoryConfig {
			return factoryConfig
		}))

		parserCreatorConfig := flam.Bag{"id": "parser", "driver": "mock"}
		expectedErr := errors.New("expected error")
		parserCreator := mocks.NewParserCreator(ctrl)
		parserCreator.EXPECT().Accept(parserCreatorConfig).Return(true).Times(1)
		parserCreator.EXPECT().Create(parserCreatorConfig).Return(nil, expectedErr).Times(1)
		require.NoError(t, container.Provide(func() config.ParserCreator {
			return parserCreator
		}, dig.Group(config.ParserCreatorGroup)))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			got, e := facade.GetParser("parser")
			assert.Nil(t, got)
			assert.ErrorIs(t, e, expectedErr)
		}))
	})

	t.Run("should return 'json' parser", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		cfg := flam.Bag{"parser": flam.Bag{"driver": config.ParserDriverJson}}
		factoryConfig := mocks.NewFactoryConfig(ctrl)
		factoryConfig.EXPECT().Get(config.PathParsers).Return(cfg).Times(1)
		require.NoError(t, container.Decorate(func(flam.FactoryConfig) flam.FactoryConfig {
			return factoryConfig
		}))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			got, e := facade.GetParser("parser")
			assert.NotNil(t, got)
			assert.NoError(t, e)
		}))
	})

	t.Run("should return 'yaml' parser", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		cfg := flam.Bag{"parser": flam.Bag{"driver": config.ParserDriverYaml}}
		factoryConfig := mocks.NewFactoryConfig(ctrl)
		factoryConfig.EXPECT().Get(config.PathParsers).Return(cfg).Times(1)
		require.NoError(t, container.Decorate(func(flam.FactoryConfig) flam.FactoryConfig {
			return factoryConfig
		}))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			got, e := facade.GetParser("parser")
			assert.NotNil(t, got)
			assert.NoError(t, e)
		}))
	})
}

func Test_Facade_AddParser(t *testing.T) {
	t.Run("should return error on nil parser", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		factoryConfig := mocks.NewFactoryConfig(ctrl)
		require.NoError(t, container.Decorate(func(flam.FactoryConfig) flam.FactoryConfig {
			return factoryConfig
		}))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			assert.ErrorIs(t, facade.AddParser("parser", nil), flam.ErrNilReference)
		}))
	})

	t.Run("should store a parser", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		cfg := flam.Bag{}
		factoryConfig := mocks.NewFactoryConfig(ctrl)
		factoryConfig.EXPECT().Get(config.PathParsers).Return(cfg).Times(1)
		require.NoError(t, container.Decorate(func(flam.FactoryConfig) flam.FactoryConfig {
			return factoryConfig
		}))

		parser := mocks.NewParser(ctrl)

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			require.NoError(t, facade.AddParser("parser", parser))

			got, e := facade.GetParser("parser")
			assert.Same(t, parser, got)
			assert.NoError(t, e)
		}))
	})

	t.Run("should return error when trying to store an already existing parser", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		cfg := flam.Bag{"parser": flam.Bag{"driver": "mock"}}
		factoryConfig := mocks.NewFactoryConfig(ctrl)
		factoryConfig.EXPECT().Get(config.PathParsers).Return(cfg).Times(1)
		require.NoError(t, container.Decorate(func(flam.FactoryConfig) flam.FactoryConfig {
			return factoryConfig
		}))

		parser := mocks.NewParser(ctrl)

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			assert.ErrorIs(t, facade.AddParser("parser", parser), flam.ErrDuplicateResource)
		}))
	})
}

func Test_Facade_HasSource(t *testing.T) {
	t.Run("should return false on unknown source", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			assert.False(t, facade.HasSource("unknown"))
		}))
	})

	t.Run("should return false if the source is on config but not loaded", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			require.NoError(t, facade.Set(config.PathSources, flam.Bag{
				"my_source": flam.Bag{
					"driver": "mock",
				}}))

			assert.False(t, facade.HasSource("my_source"))
		}))
	})

	t.Run("should return true if the source was added directly", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		source := mocks.NewSource(ctrl)
		source.EXPECT().Get("", flam.Bag{}).Return(flam.Bag{}).Times(1)

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			require.NoError(t, facade.AddSource("source", source))

			assert.True(t, facade.HasSource("source"))
		}))
	})
}

func Test_Facade_ListSources(t *testing.T) {
	t.Run("should return an empty list if no sources were added", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			assert.Empty(t, facade.ListSources())
		}))
	})

	t.Run("should return a ordered list of sources", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		source1 := mocks.NewSource(ctrl)
		source1.EXPECT().Get("", flam.Bag{}).Return(flam.Bag{}).AnyTimes()
		source1.EXPECT().GetPriority().Return(1).AnyTimes()

		source2 := mocks.NewSource(ctrl)
		source2.EXPECT().Get("", flam.Bag{}).Return(flam.Bag{}).AnyTimes()
		source2.EXPECT().GetPriority().Return(2).AnyTimes()

		source3 := mocks.NewSource(ctrl)
		source3.EXPECT().Get("", flam.Bag{}).Return(flam.Bag{}).AnyTimes()
		source3.EXPECT().GetPriority().Return(3).AnyTimes()

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			require.NoError(t, facade.AddSource("zulu", source3))
			require.NoError(t, facade.AddSource("alpha", source1))
			require.NoError(t, facade.AddSource("charlie", source2))

			assert.ElementsMatch(t, []string{"alpha", "charlie", "zulu"}, facade.ListSources())
		}))
	})
}

func Test_Facade_GetSource(t *testing.T) {
	t.Run("should return error on unknown source", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			got, e := facade.GetSource("unknown")
			assert.Nil(t, got)
			assert.ErrorIs(t, e, config.ErrSourceNotFound)
		}))
	})

	t.Run("should return requested source", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		source := mocks.NewSource(ctrl)
		source.EXPECT().Get("", flam.Bag{}).Return(flam.Bag{}).AnyTimes()
		source.EXPECT().GetPriority().Return(1).AnyTimes()

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			require.NoError(t, facade.AddSource("source", source))

			got, e := facade.GetSource("source")
			assert.Same(t, source, got)
			assert.NoError(t, e)
		}))
	})
}

func Test_Facade_AddSource(t *testing.T) {
	t.Run("should return error on nil source", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			assert.ErrorIs(t, facade.AddSource("source", nil), flam.ErrNilReference)
		}))
	})

	t.Run("should return error on duplicate source", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		source := mocks.NewSource(ctrl)
		source.EXPECT().Get("", flam.Bag{}).Return(flam.Bag{}).AnyTimes()
		source.EXPECT().GetPriority().Return(1).AnyTimes()

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			require.NoError(t, facade.AddSource("source", source))
			assert.ErrorIs(t, facade.AddSource("source", source), config.ErrDuplicateSource)
		}))
	})

	t.Run("should add source", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		source := mocks.NewSource(ctrl)
		source.EXPECT().Get("", flam.Bag{}).Return(flam.Bag{}).AnyTimes()
		source.EXPECT().GetPriority().Return(1).AnyTimes()

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			require.NoError(t, facade.AddSource("source", source))

			got, e := facade.GetSource("source")
			assert.Same(t, source, got)
			assert.NoError(t, e)
		}))
	})

	t.Run("should override existing source values if with higher priority", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		source1 := mocks.NewSource(ctrl)
		source1.EXPECT().Get("", flam.Bag{}).Return(flam.Bag{"field": "value1"}).AnyTimes()
		source1.EXPECT().GetPriority().Return(1).AnyTimes()

		source2 := mocks.NewSource(ctrl)
		source2.EXPECT().Get("", flam.Bag{}).Return(flam.Bag{"field": "value2"}).AnyTimes()
		source2.EXPECT().GetPriority().Return(2).AnyTimes()

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			require.NoError(t, facade.AddSource("source1", source1))

			got, e := facade.GetSource("source1")
			assert.Same(t, source1, got)
			assert.NoError(t, e)
			assert.Equal(t, "value1", facade.Get("field"))

			require.NoError(t, facade.AddSource("source2", source2))

			got, e = facade.GetSource("source2")
			assert.Same(t, source2, got)
			assert.NoError(t, e)
			assert.Equal(t, "value2", facade.Get("field"))
		}))
	})
}

func Test_Facade_SetSourcePriority(t *testing.T) {
	t.Run("should return error on invalid source", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			assert.ErrorIs(t, facade.SetSourcePriority("source", 1), config.ErrSourceNotFound)
		}))
	})

	t.Run("should override existing source priority and rearrange sources", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		source1 := mocks.NewSource(ctrl)
		source1.EXPECT().Get("", flam.Bag{}).Return(flam.Bag{"field": "value1"}).AnyTimes()
		gomock.InOrder(
			source1.EXPECT().GetPriority().Return(1),
			source1.EXPECT().GetPriority().Return(3),
		)
		source1.EXPECT().SetPriority(3).AnyTimes()

		source2 := mocks.NewSource(ctrl)
		source2.EXPECT().Get("", flam.Bag{}).Return(flam.Bag{"field": "value2"}).AnyTimes()
		source2.EXPECT().GetPriority().Return(2).AnyTimes()

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			require.NoError(t, facade.AddSource("source1", source1))
			assert.Equal(t, "value1", facade.Get("field"))

			require.NoError(t, facade.AddSource("source2", source2))
			assert.Equal(t, "value2", facade.Get("field"))

			require.NoError(t, facade.SetSourcePriority("source1", 3))
			assert.Equal(t, "value1", facade.Get("field"))
		}))
	})

	t.Run("real source test", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		require.NoError(t, os.Setenv("ENV_FIELD_1", "value_1"))
		require.NoError(t, os.Setenv("ENV_FIELD_2", "value_2"))
		defer func() {
			_ = os.Unsetenv("ENV_FIELD_1")
			_ = os.Unsetenv("ENV_FIELD_2")
		}()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(config.PathSources, flam.Bag{
			"my_source_1": flam.Bag{
				"driver":   config.SourceDriverEnv,
				"priority": 1,
				"files":    []string{},
				"mappings": flam.Bag{
					"ENV_FIELD_1": "env.field",
				},
			},
			"my_source_2": flam.Bag{
				"driver":   config.SourceDriverEnv,
				"priority": 2,
				"files":    []string{},
				"mappings": flam.Bag{
					"ENV_FIELD_2": "env.field",
				},
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			assert.Equal(t, "value_2", facade.Get("env.field"))

			require.NoError(t, facade.SetSourcePriority("my_source_1", 3))
			assert.Equal(t, "value_1", facade.Get("env.field"))
		}))
	})
}

func Test_Facade_RemoveSource(t *testing.T) {
	t.Run("should return error on invalid source", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			assert.ErrorIs(t, facade.RemoveSource("source"), config.ErrSourceNotFound)
		}))
	})

	t.Run("should return the error if the source returns any on closing", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		source1 := mocks.NewSource(ctrl)
		source1.EXPECT().Get("", flam.Bag{}).Return(flam.Bag{"field": "value1"}).AnyTimes()
		source1.EXPECT().GetPriority().Return(1).AnyTimes()

		expectedErr := errors.New("close error")
		source2 := mocks.NewSource(ctrl)
		source2.EXPECT().Get("", flam.Bag{}).Return(flam.Bag{"field": "value2"}).AnyTimes()
		source2.EXPECT().GetPriority().Return(2).AnyTimes()
		source2.EXPECT().Close().Return(expectedErr).Times(1)

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			require.NoError(t, facade.AddSource("source1", source1))
			assert.Equal(t, "value1", facade.Get("field"))

			require.NoError(t, facade.AddSource("source2", source2))
			assert.Equal(t, "value2", facade.Get("field"))

			require.ErrorIs(t, facade.RemoveSource("source2"), expectedErr)
			assert.Equal(t, "value2", facade.Get("field"))
		}))
	})

	t.Run("should remove source and rearrange values", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		source1 := mocks.NewSource(ctrl)
		source1.EXPECT().Get("", flam.Bag{}).Return(flam.Bag{"field": "value1"}).AnyTimes()
		source1.EXPECT().GetPriority().Return(1).AnyTimes()
		source1.EXPECT().Close().Return(nil).AnyTimes()

		source2 := mocks.NewSource(ctrl)
		source2.EXPECT().Get("", flam.Bag{}).Return(flam.Bag{"field": "value2"}).AnyTimes()
		source2.EXPECT().GetPriority().Return(2).AnyTimes()
		source2.EXPECT().Close().Return(nil).AnyTimes()

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			require.NoError(t, facade.AddSource("source1", source1))
			assert.Equal(t, "value1", facade.Get("field"))

			require.NoError(t, facade.AddSource("source2", source2))
			assert.Equal(t, "value2", facade.Get("field"))

			require.NoError(t, facade.RemoveSource("source2"))
			assert.Equal(t, "value1", facade.Get("field"))
		}))
	})
}

func Test_Facade_RemoveAllSources(t *testing.T) {
	t.Run("should return the error if the source returns any on closing", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		expectedErr := errors.New("close error")
		source := mocks.NewSource(ctrl)
		source.EXPECT().Get("", flam.Bag{}).Return(flam.Bag{"field": "value1"})
		source.EXPECT().Close().Return(expectedErr)

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			require.NoError(t, facade.AddSource("source", source))

			assert.ErrorIs(t, facade.RemoveAllSources(), expectedErr)
		}))
	})

	t.Run("should remove all sources", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		source1 := mocks.NewSource(ctrl)
		source1.EXPECT().Get("", flam.Bag{}).Return(flam.Bag{"field": "value1"}).AnyTimes()
		source1.EXPECT().GetPriority().Return(1).AnyTimes()
		source1.EXPECT().Close().Return(nil).Times(1)

		source2 := mocks.NewSource(ctrl)
		source2.EXPECT().Get("", flam.Bag{}).Return(flam.Bag{"field": "value2"}).AnyTimes()
		source2.EXPECT().GetPriority().Return(1).AnyTimes()
		source2.EXPECT().Close().Return(nil).Times(1)

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			require.NoError(t, facade.AddSource("source1", source1))
			require.NoError(t, facade.AddSource("source2", source2))

			require.NoError(t, facade.RemoveAllSources())

			assert.Empty(t, facade.ListSources())
			assert.False(t, facade.Has("field"))
		}))
	})
}

func Test_Facade_ReloadSources(t *testing.T) {
	t.Run("should no-op if none of the sources is a observable source", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		source := mocks.NewSource(ctrl)
		source.EXPECT().Get("", flam.Bag{}).Return(flam.Bag{"field": "value1"}).AnyTimes()
		source.EXPECT().GetPriority().Return(1).AnyTimes()

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			require.NoError(t, facade.AddSource("source", source))

			assert.NoError(t, facade.ReloadSources())
		}))
	})

	t.Run("should no-op if no source reloads", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		source := mocks.NewObservableSource(ctrl)
		source.EXPECT().Get("", flam.Bag{}).Return(flam.Bag{"field": "value1"}).Times(1)
		source.EXPECT().Reload().Return(false, nil).Times(1)

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			require.NoError(t, facade.AddSource("source", source))

			assert.NoError(t, facade.ReloadSources())
		}))
	})

	t.Run("should return the source reload error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		expectedErr := errors.New("reload error")
		source := mocks.NewObservableSource(ctrl)
		source.EXPECT().Get("", flam.Bag{}).Return(flam.Bag{"field": "value1"}).Times(1)
		source.EXPECT().Reload().Return(false, expectedErr).Times(1)

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			require.NoError(t, facade.AddSource("source", source))

			assert.ErrorIs(t, facade.ReloadSources(), expectedErr)
		}))
	})

	t.Run("should reload all sources", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		source1 := mocks.NewObservableSource(ctrl)
		source1.EXPECT().Get("", flam.Bag{}).Return(flam.Bag{"field1": "value-y"})
		source1.EXPECT().Get("", flam.Bag{}).Return(flam.Bag{"field1": "value-y"})
		source1.EXPECT().Get("", flam.Bag{}).Return(flam.Bag{"field1": "value-y", "field2": "value2"})
		source1.EXPECT().GetPriority().Return(1).AnyTimes()
		source1.EXPECT().Reload().Return(true, nil).Times(1)

		source2 := mocks.NewObservableSource(ctrl)
		source2.EXPECT().Get("", flam.Bag{}).Return(flam.Bag{"field1": "value-x"})
		source2.EXPECT().Get("", flam.Bag{}).Return(flam.Bag{"field1": "value-x", "field3": "value3"})
		source2.EXPECT().GetPriority().Return(2).AnyTimes()
		source2.EXPECT().Reload().Return(true, nil).Times(1)

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			require.NoError(t, facade.AddSource("source1", source1))
			require.NoError(t, facade.AddSource("source2", source2))

			assert.Equal(t, "value-x", facade.Get("field1"))
			assert.Nil(t, facade.Get("field2"))
			assert.Nil(t, facade.Get("field3"))

			require.NoError(t, facade.ReloadSources())
			assert.Equal(t, "value-x", facade.Get("field1"))
			assert.Equal(t, "value2", facade.Get("field2"))
			assert.Equal(t, "value3", facade.Get("field3"))
		}))
	})
}

func Test_Facade_HasObserver(t *testing.T) {
	t.Run("should return false if the observer is not present", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		observer := config.Observer(func(old any, new any) {})

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			require.NoError(t, facade.AddObserver("id", "field", observer))
			assert.False(t, facade.HasObserver("other", "field"))
			assert.False(t, facade.HasObserver("id", "other"))
		}))
	})

	t.Run("should return true if the observer is present", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		observer := config.Observer(func(old any, new any) {})

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			require.NoError(t, facade.AddObserver("id", "field", observer))
			assert.True(t, facade.HasObserver("id", "field"))
		}))
	})
}

func Test_Facade_AddObserver(t *testing.T) {
	t.Run("should return error on nil callback", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			assert.ErrorIs(t, facade.AddObserver("", "", nil), flam.ErrNilReference)
		}))
	})

	t.Run("should store observer and be called on value change", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		observer := config.Observer(func(old any, new any) {})

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			require.NoError(t, facade.AddObserver("id", "field", observer))
			assert.ErrorIs(t, facade.AddObserver("id", "field", observer), config.ErrDuplicateObserver)
		}))
	})

	t.Run("should store observer and be called on value change", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		called := false
		observer := config.Observer(func(old any, new any) {
			assert.Nil(t, old)
			assert.Equal(t, "value1", new)
			called = true
		})

		source := mocks.NewObservableSource(ctrl)
		source.EXPECT().Get("", flam.Bag{}).Return(flam.Bag{"field": "value1"}).Times(1)

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			require.NoError(t, facade.AddObserver("id", "field", observer))

			require.NoError(t, facade.AddSource("source", source))
			assert.True(t, called)
			assert.Equal(t, "value1", facade.Get("field"))
		}))
	})
}

func Test_Facade_RemoveObserver(t *testing.T) {
	t.Run("should no error if observer does not exist", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			require.NoError(t, facade.RemoveObserver("id"))
		}))
	})

	t.Run("should remove observer", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, config.NewProvider().Register(container))

		called := false
		observer := config.Observer(func(old any, new any) {
			called = true
		})

		source := mocks.NewObservableSource(ctrl)
		source.EXPECT().Get("", flam.Bag{}).Return(flam.Bag{"field": "value1"}).Times(1)

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			require.NoError(t, facade.AddObserver("id", "field", observer))
			require.NoError(t, facade.RemoveObserver("id"))

			require.NoError(t, facade.AddSource("source", source))

			assert.False(t, called)
			assert.Equal(t, "value1", facade.Get("field"))
		}))
	})
}
