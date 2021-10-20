package openapi

import (
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/yeqown/log"
)

type Format string

const (
	Format_Properties Format = "properties"
	Format_XML        Format = "xml"
	Format_JSON       Format = "json"
	Format_YML        Format = "yml"
	Format_YAML       Format = "yaml"
	Format_TEXT       Format = "txt"
)

func getMapper(input map[string]string) func(k string) string {
	return func(k string) string {
		v, ok := input[k]
		if !ok {
			return ""
		}

		return v
	}
}

func expand(uri string, input map[string]string) string {
	mapper := getMapper(input)
	out := os.Expand(uri, mapper)

	log.
		WithFields(log.Fields{
			"in":     uri,
			"mapper": input,
			"out":    out,
		}).
		Debug("expand uri")

	return out
}

func handleResponseError(r *resty.Response) error {
	if r == nil {
		return nil
	}

	if r.StatusCode() == http.StatusOK {
		return nil
	}

	switch r.StatusCode() {
	case http.StatusBadRequest:
		return errors.New("Bad Request")
	case http.StatusUnauthorized:
		return errors.New("Unauthorized")
	case http.StatusForbidden:
		return errors.New("Forbidden")
	case http.StatusNotFound:
		return errors.New("Not Found")
	case http.StatusMethodNotAllowed:
		return errors.New("Method not Allowed")
	case http.StatusInternalServerError:
		return errors.New("Internal Server Error")
	}

	return nil
}
