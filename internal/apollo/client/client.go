package client

import (
	"github.com/go-resty/resty/v2"
)

type Client interface {
}

type openClient struct {
	cc *resty.Client
}
