package configuration

import (
	"sync"
	"time"

	"github.com/llnw/llnw-sdk-go"
)

type ConfigurationClient struct {
	Auth                             *llnw.Auth
	BaseUrl                          string
	rateLimiter                      <-chan time.Time
	configOptionLock                 sync.Mutex
	configOptionArgumentIntegerCache map[string][]bool
}

func NewClient(apiUser string, apiKey string) *ConfigurationClient {
	return NewClientOverrideBaseUrl(apiUser, apiKey, "https://apis.llnw.com/config-api/v1")
}

func NewClientOverrideBaseUrl(apiUser string, apiKey string, baseUrl string) *ConfigurationClient {
	a := &llnw.Auth{}
	a.APIUser = apiUser
	a.APIKey = apiKey

	c := &ConfigurationClient{}
	c.Auth = a
	c.BaseUrl = baseUrl

	c.rateLimiter = time.Tick(1200 * time.Millisecond)

	return c
}

func (c *ConfigurationClient) SetUserAgent(userAgent string) {
	c.Auth.UserAgent = userAgent
}
