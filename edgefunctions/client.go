package edgefunctions

import (
	"github.com/llnw/llnw-sdk-go"
)

type EdgeFunctionsClient struct {
	Auth    *llnw.Auth
	BaseUrl string
}

func NewClient(apiUser string, apiKey string) *EdgeFunctionsClient {
	return NewClientOverrideBaseUrl(apiUser, apiKey, "https://api.faas.llnw.net/ef-api/v1")
}

func NewClientOverrideBaseUrl(apiUser string, apiKey string, baseUrl string) *EdgeFunctionsClient {
	a := &llnw.Auth{}
	a.APIUser = apiUser
	a.APIKey = apiKey

	c := &EdgeFunctionsClient{}
	c.Auth = a
	c.BaseUrl = baseUrl

	return c
}

func (c *EdgeFunctionsClient) SetUserAgent(userAgent string) {
	c.Auth.UserAgent = userAgent
}
