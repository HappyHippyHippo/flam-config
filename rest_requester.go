package config

import (
	"net/http"
)

type RestRequester interface {
	Do(req *http.Request) (*http.Response, error)
}
