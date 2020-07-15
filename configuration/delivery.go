package configuration

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Start - DeliveryServiceInstance types
type DeliveryServiceInstance struct {
	UUID      string                      `json:"uuid"`
	IsLatest  bool                        `json:"isLatest"`
	IsEnabled bool                        `json:"isEnabled"`
	Revision  Revision                    `json:"revision"`
	Accounts  []Account                   `json:"accounts"`
	Shortname string                      `json:"shortname"`
	Body      DeliveryServiceInstanceBody `json:"body"`
}

type DeliveryServiceInstanceCreateRequest struct {
	Body     DeliveryServiceInstanceBody `json:"body"`
	Accounts []Account                   `json:"accounts"`
}

type DeliveryServiceInstanceUpdateRequest struct {
	UUID     string                      `json:"uuid"`
	Body     DeliveryServiceInstanceBody `json:"body"`
	Accounts []Account                   `json:"accounts"`
}

type Revision struct {
	CreatedBy     string `json:"createdBy"`
	CreatedDate   int64  `json:"createdDate"`
	VersionNumber int    `json:"versionNumber"`
}

type DeliveryServiceInstanceBody struct {
	ServiceProfileName string        `json:"serviceProfileName"`
	ProtocolSets       []ProtocolSet `json:"protocolSets"`
	PublishedHostname  string        `json:"publishedHostname"`
	SourceHostname     string        `json:"sourceHostname"`
	PublishedURLPath   string        `json:"publishedUrlPath"`
	SourceURLPath      string        `json:"sourceUrlPath"`
	ServiceKey         ServiceKey    `json:"serviceKey"`
}

type ProtocolSet struct {
	PublishedProtocol string   `json:"publishedProtocol"`
	SourceProtocol    string   `json:"sourceProtocol"`
	SourcePort        *int     `json:"sourcePort"`
	Options           []Option `json:"options"`
}

type Option struct {
	Name       string        `json:"name"`
	Parameters []interface{} `json:"parameters"`
}

type ServiceKey struct {
	Name string `json:"name"`
}

type Account struct {
	Shortname string `json:"shortname"`
}

// End - DeliveryServiceInstance types

// Start - ConfigOption types
type ConfigOptionsResponse struct {
	Results []ConfigOption `json:"results"`
}

type ConfigOption struct {
	Body ConfigOptionBody `json:"body"`
}

type ConfigOptionBody struct {
	Name    string              `json:"optionName"`
	Details ConfigOptionDetails `json:"optionDetails"`
}

type ConfigOptionDetails struct {
	Arguments []ConfigOptionArgument `json:"argumentList"`
}

type ConfigOptionArgument struct {
	Type string `json:"type"`
}

// End - ConfigOption types

func (c *ConfigurationClient) GetConfigurationOptions(shortname string, profileName string) ([]ConfigOption, *http.Response, error) {
	<-c.rateLimiter
	body, response, err := c.Auth.HTTPGet(fmt.Sprintf("%s/configoption/shortname/%s/svcProf/%s", c.BaseUrl, shortname, profileName))

	if err != nil {
		return nil, response, err
	}

	configOptionsResponse := &ConfigOptionsResponse{}
	json.Unmarshal(body, configOptionsResponse)

	return configOptionsResponse.Results, response, nil
}

func (c *ConfigurationClient) IsOptionArgumentInteger(shortname string, profileName string, optionName string, argumentPosition int) (bool, error) {
	c.configOptionLock.Lock()
	defer c.configOptionLock.Unlock()

	if c.configOptionArgumentIntegerCache == nil {
		configOptions, _, err := c.GetConfigurationOptions(shortname, profileName)
		if err != nil {
			return false, err
		}
		optionsMap := map[string][]bool{}
		for _, option := range configOptions {
			var argumentIntegerList []bool
			for _, argument := range option.Body.Details.Arguments {
				argumentIntegerList = append(argumentIntegerList, argument.Type == "Int")
			}
			optionsMap[option.Body.Name] = argumentIntegerList
		}
		c.configOptionArgumentIntegerCache = optionsMap
	}

	if argumentIntegerList, ok := c.configOptionArgumentIntegerCache[optionName]; ok {
		if argumentPosition >= len(argumentIntegerList) {
			return false, nil
		} else {
			return argumentIntegerList[argumentPosition], nil
		}
	} else {
		return false, nil
	}
}

func (c *ConfigurationClient) GetDeliveryServiceInstance(uuid string) (*DeliveryServiceInstance, *http.Response, error) {
	<-c.rateLimiter
	deliveryServiceInstance := &DeliveryServiceInstance{}

	body, response, err := c.Auth.HTTPGet(c.BaseUrl + "/svcinst/delivery/" + uuid)

	if err != nil {
		return nil, response, err
	}

	json.Unmarshal(body, deliveryServiceInstance)

	return deliveryServiceInstance, response, nil
}

func (c *ConfigurationClient) CreateDeliveryServiceInstance(body *DeliveryServiceInstanceBody, shortname string) (*DeliveryServiceInstance, *http.Response, error) {
	<-c.rateLimiter
	request := &DeliveryServiceInstanceCreateRequest{
		Body: *body,
		Accounts: []Account{
			Account{
				Shortname: shortname,
			},
		},
	}

	jsonRequest, _ := json.Marshal(request)

	respBody, response, err := c.Auth.HTTPPost(c.BaseUrl+"/svcinst/delivery", string(jsonRequest))

	if err != nil {
		return nil, response, err
	}

	deliveryServiceInstance := &DeliveryServiceInstance{}
	json.Unmarshal(respBody, deliveryServiceInstance)

	return deliveryServiceInstance, response, nil
}

func (c *ConfigurationClient) UpdateDeliveryServiceInstance(uuid string, body *DeliveryServiceInstanceBody, shortname string) (*DeliveryServiceInstance, *http.Response, error) {
	<-c.rateLimiter
	request := &DeliveryServiceInstanceUpdateRequest{
		UUID: uuid,
		Body: *body,
		Accounts: []Account{
			Account{
				Shortname: shortname,
			},
		},
	}

	jsonRequest, _ := json.Marshal(request)

	respBody, response, err := c.Auth.HTTPPut(c.BaseUrl+"/svcinst/delivery/"+uuid, string(jsonRequest))

	if err != nil {
		return nil, response, err
	}

	deliveryServiceInstance := &DeliveryServiceInstance{}
	json.Unmarshal(respBody, deliveryServiceInstance)

	return deliveryServiceInstance, response, nil
}

func (c *ConfigurationClient) DeleteDeliveryServiceInstance(uuid string) (*DeliveryServiceInstance, *http.Response, error) {
	<-c.rateLimiter
	body, response, err := c.Auth.HTTPDelete(c.BaseUrl + "/svcinst/delivery/" + uuid)

	if err != nil {
		return nil, response, err
	}

	deliveryServiceInstance := &DeliveryServiceInstance{}
	json.Unmarshal(body, deliveryServiceInstance)

	return deliveryServiceInstance, response, nil
}
