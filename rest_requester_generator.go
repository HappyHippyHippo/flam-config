package config

import (
	"net/http"
)

type RestRequesterGenerator interface {
	Create() (RestRequester, error)
}

type restRequesterGenerator struct{}

func newRestRequesterGenerator() RestRequesterGenerator {
	return &restRequesterGenerator{}
}

func (restRequesterGenerator) Create() (RestRequester, error) {
	return &http.Client{}, nil
}
