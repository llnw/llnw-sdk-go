package configuration

import (
	"encoding/json"
	"net/http"
)

type IPAllowList struct {
	IPRanges []string `json:"ipAllowList"`
	Version  int      `json:"version"`
}

func (c ConfigurationClient) GetIPAllowList() (*IPAllowList, *http.Response, error) {
	object := &IPAllowList{}
	body, response, err := c.Auth.HTTPGet("https://control.llnw.com/aportal/api/ipam/getIpAllowList.do")

	if err != nil {
		return nil, response, err
	}

	json.Unmarshal(body, &object)
	return object, response, nil
}
