package api

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/go-resty/resty/v2"
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

// NotAllowedFormat use whitelist for allowed formats.
func NotAllowedFormat(ext Format) bool {
	switch ext {
	case Format_TEXT, Format_JSON, Format_XML, Format_YAML, Format_YML:
		return false
	}

	return true
}

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

type Exception struct {
	Exception string `json:"exception"`
	Message   string `json:"message"`
	// FIXED(@yeqown): status in response body is float. ops!
	// use http response status code as instead.
	Status    int    `json:"-"`
	Timestamp string `json:"timestamp"`
}

func (e Exception) Error() string {
	return strconv.Itoa(e.Status) + ": " + e.Message
}

func handleResponseError(r *resty.Response) error {
	if r == nil {
		return nil
	}

	if r.StatusCode() == http.StatusOK {
		return nil
	}

	e := new(Exception)
	if err := json.Unmarshal(r.Body(), e); err != nil {
		log.
			WithFields(log.Fields{
				"rawBody": string(r.Body()),
				"err":     err,
			}).
			Error("handleResponseError could not unmarshal Exception")
		return nil
	}
	e.Status = r.StatusCode()

	return e
}
