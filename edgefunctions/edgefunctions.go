package edgefunctions

import (
	"encoding/json"
	"net/http"
)

type EdgeFunction struct {
	Name                 string                `json:"name,omitempty"`
	Description          string                `json:"description,omitempty"`
	FunctionArchive      []byte                `json:"functionArchive,omitempty"`
	Handler              string                `json:"handler,omitempty"`
	Runtime              string                `json:"runtime,omitempty"`
	Memory               int                   `json:"memory,omitempty"`
	Timeout              int                   `json:"timeout,omitempty"`
	CanDebug             bool                  `json:"canDebug,omitempty"`
	Sha256               string                `json:"sha256,omitempty"`
	EnvironmentVariables []EnvironmentVariable `json:"environmentVariables"`
	ReservedConcurrency  int                   `json:"reservedConcurrency,omitempty"`
	RevisionID           int                   `json:"revisionId,omitempty"`
	Version              int                   `json:"version,omitempty"`
}

type ReservedConcurrency struct {
	ReservedConcurrency int `json:"reservedConcurrency"`
}

type EnvironmentVariable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type EdgeFunctionAlias struct {
	Name            string `json:"name,omitempty"`
	Description     string `json:"description,omitempty"`
	Function        string `json:"function,omitempty"`
	FunctionVersion string `json:"functionVersion,omitempty"`
	RevisionID      int    `json:"revisionId,omitempty"`
}

func (c *EdgeFunctionsClient) GetEdgeFunction(name string, shortname string) (*EdgeFunction, *http.Response, error) {
	body, response, err := c.Auth.HTTPGet(c.BaseUrl + "/" + shortname + "/functions/" + name)

	if err != nil {
		return nil, response, err
	}

	edgeFunctionResponse := &EdgeFunction{}
	json.Unmarshal(body, edgeFunctionResponse)

	return edgeFunctionResponse, response, nil
}

func (c *EdgeFunctionsClient) CreateEdgeFunction(shortname string, edgeFunction *EdgeFunction) (*EdgeFunction, *http.Response, error) {
	jsonRequest, _ := json.Marshal(edgeFunction)

	body, response, err := c.Auth.HTTPPost(c.BaseUrl+"/"+shortname+"/functions", string(jsonRequest))

	if err != nil {
		return nil, response, err
	}

	edgeFunctionResponse := &EdgeFunction{}
	json.Unmarshal(body, edgeFunctionResponse)

	return edgeFunctionResponse, response, nil
}

func (c *EdgeFunctionsClient) UpdateEdgeFunctionCode(name string, shortname string, functionArchive []byte) (*EdgeFunction, *http.Response, error) {
	edgeFunction := &EdgeFunction{
		FunctionArchive: functionArchive,
	}

	jsonRequest, _ := json.Marshal(edgeFunction)

	body, response, err := c.Auth.HTTPPut(c.BaseUrl+"/"+shortname+"/functions/"+name, string(jsonRequest))

	if err != nil {
		return nil, response, err
	}

	edgeFunctionResponse := &EdgeFunction{}
	json.Unmarshal(body, edgeFunctionResponse)

	return edgeFunctionResponse, response, nil
}

func (c *EdgeFunctionsClient) UpdateEdgeFunctionConfiguration(name string, shortname string, edgeFunction *EdgeFunction) (*EdgeFunction, *http.Response, error) {
	jsonRequest, _ := json.Marshal(edgeFunction)

	body, response, err := c.Auth.HTTPPut(c.BaseUrl+"/"+shortname+"/functions/"+name+"/configuration", string(jsonRequest))

	if err != nil {
		return nil, response, err
	}

	edgeFunctionResponse := &EdgeFunction{}
	json.Unmarshal(body, edgeFunctionResponse)

	return edgeFunctionResponse, response, nil
}

func (c *EdgeFunctionsClient) DeleteEdgeFunction(name string, shortname string) (*http.Response, error) {
	_, response, err := c.Auth.HTTPDelete(c.BaseUrl + "/" + shortname + "/functions/" + name)

	if err != nil {
		return response, err
	}

	return response, nil
}

func (c *EdgeFunctionsClient) SetEdgeFunctionConcurrency(fnName string, shortname string, concurrency int) (*http.Response, error) {
	jsonRequest, err := json.Marshal(ReservedConcurrency{ReservedConcurrency: concurrency})
	if err != nil {
		return nil, err
	}
	_, response, err := c.Auth.HTTPPut(c.BaseUrl+"/"+shortname+"/functions/"+fnName+"/concurrency", string(jsonRequest))
	return response, err
}

func (c *EdgeFunctionsClient) CreateEdgeFunctionAlias(fnName, shortname string, alias *EdgeFunctionAlias) (*EdgeFunctionAlias, *http.Response, error) {
	jsonRequest, err := json.Marshal(alias)
	if err != nil {
		return nil, nil, err
	}

	body, response, err := c.Auth.HTTPPost(c.BaseUrl+"/"+shortname+"/functions/"+fnName+"/aliases", string(jsonRequest))
	if err != nil {
		return nil, response, err
	}

	aliasResponse := &EdgeFunctionAlias{}
	if err = json.Unmarshal(body, aliasResponse); err != nil {
		return nil, response, err
	}
	return aliasResponse, response, nil
}

func (c *EdgeFunctionsClient) UpdateEdgeFunctionAlias(fnName, shortname, aliasName string, alias *EdgeFunctionAlias) (*EdgeFunctionAlias, *http.Response, error) {
	jsonRequest, err := json.Marshal(alias)
	if err != nil {
		return nil, nil, err
	}

	body, response, err := c.Auth.HTTPPut(c.BaseUrl+"/"+shortname+"/functions/"+fnName+"/aliases/"+aliasName, string(jsonRequest))
	if err != nil {
		return nil, response, err
	}

	aliasResponse := &EdgeFunctionAlias{}
	if err = json.Unmarshal(body, aliasResponse); err != nil {
		return nil, response, err
	}
	return aliasResponse, response, nil
}

func (c *EdgeFunctionsClient) GetEdgeFunctionAlias(fnName, shortname, aliasName string) (*EdgeFunctionAlias, *http.Response, error) {

	body, response, err := c.Auth.HTTPGet(c.BaseUrl + "/" + shortname + "/functions/" + fnName + "/aliases/" + aliasName)
	if err != nil {
		return nil, response, err
	}

	aliasResponse := &EdgeFunctionAlias{}
	if err = json.Unmarshal(body, aliasResponse); err != nil {
		return nil, response, err
	}
	return aliasResponse, response, nil
}

func (c *EdgeFunctionsClient) DeleteEdgeFunctionAlias(fnName, shortname, aliasName string) (*http.Response, error) {

	_, response, err := c.Auth.HTTPDelete(c.BaseUrl + "/" + shortname + "/functions/" + fnName + "/aliases/" + aliasName)

	return response, err
}
