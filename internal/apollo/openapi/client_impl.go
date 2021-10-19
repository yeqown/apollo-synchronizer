package openapi

import "github.com/go-resty/resty/v2"

type openapiClient struct {
	config *Config
	cc     *resty.Client
}

func New(c *Config) *openapiClient {
	client := resty.
		New().
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		SetAuthToken(c.Token)

	return &openapiClient{
		config: c,
		cc:     client,
	}
}
