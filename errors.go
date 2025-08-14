package config

import (
	"errors"
	"fmt"

	flam "github.com/happyhippyhippo/flam"
)

var (
	ErrRestConfigNotFound    = errors.New("rest config data not found")
	ErrRestInvalidConfig     = errors.New("invalid rest config data")
	ErrRestTimestampNotFound = errors.New("rest config timestamp not found")
	ErrRestInvalidTimestamp  = errors.New("invalid rest config timestamp")
	ErrSourceNotFound        = errors.New("config source not found")
	ErrDuplicateSource       = errors.New("duplicate config source")
	ErrDuplicateObserver     = errors.New("duplicate config observer")
)

func newErrNilReference(
	field string,
) error {
	return flam.NewErrorFrom(
		flam.ErrNilReference,
		field)
}

func newErrRestConfigNotFound(
	path string,
	config flam.Bag,
) error {
	return flam.NewErrorFrom(
		ErrRestConfigNotFound,
		fmt.Sprintf("%s => %v", path, config))
}

func newErrRestInvalidConfig(
	path string,
	value any) error {
	return flam.NewErrorFrom(
		ErrRestInvalidConfig,
		fmt.Sprintf("%s => %v", path, value))
}

func newErrRestTimestampNotFound(
	path string,
	config flam.Bag,
) error {
	return flam.NewErrorFrom(
		ErrRestTimestampNotFound,
		fmt.Sprintf("%s => %v", path, config))
}

func newErrRestInvalidTimestamp(
	path string,
	value any,
) error {
	return flam.NewErrorFrom(
		ErrRestInvalidTimestamp,
		fmt.Sprintf("%s => %v", path, value))
}

func newErrSourceNotFound(
	id string,
) error {
	return flam.NewErrorFrom(
		ErrSourceNotFound,
		id)
}

func newErrDuplicateSource(
	id string,
) error {
	return flam.NewErrorFrom(
		ErrDuplicateSource,
		id)
}

func newErrDuplicateObserver(
	path string,
	id string,
) error {
	return flam.NewErrorFrom(
		ErrDuplicateObserver,
		fmt.Sprintf("%s => %s", path, id))
}
