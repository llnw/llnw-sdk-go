package llnw

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Auth attributes for API
type Auth struct {
	APIUser   string
	APIKey    string
	UserAgent string
}

// HTTPGet performs a GET on the said url
func (a Auth) HTTPGet(url string) ([]byte, *http.Response, error) {
	return a.httpRequestWithoutBody("GET", url)
}

// HTTPPost performs a POST on the said url with the said requestBody
func (a Auth) HTTPPost(url string, requestBody string) ([]byte, *http.Response, error) {
	return a.httpRequestWithBody("POST", url, requestBody)
}

// HTTPPut performs a PUT on the said url with the said requestBody
func (a Auth) HTTPPut(url string, requestBody string) ([]byte, *http.Response, error) {
	return a.httpRequestWithBody("PUT", url, requestBody)
}

// HTTPDelete performs a DELETE on the said url
func (a Auth) HTTPDelete(url string) ([]byte, *http.Response, error) {
	return a.httpRequestWithoutBody("DELETE", url)
}

func (a Auth) httpRequestWithoutBody(method string, url string) ([]byte, *http.Response, error) {
	var client = &http.Client{
		Timeout: time.Second * 30,
	}

	req, _ := http.NewRequest(method, url, nil)

	headers := a.buildAuthHeaders(url, method, "")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	if a.UserAgent != "" {
		req.Header.Set("User-Agent", a.UserAgent)
	}

	log.Printf("[DEBUG] %s %s", method, url)
	resp, err := client.Do(req)

	if err != nil {
		return nil, resp, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, resp, fmt.Errorf("non-2XX status code from API, got status %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	return bodyBytes, resp, nil
}

func (a Auth) httpRequestWithBody(method string, url string, requestBody string) ([]byte, *http.Response, error) {
	var client = &http.Client{
		Timeout: time.Second * 30,
	}

	req, _ := http.NewRequest(method, url, bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")

	headers := a.buildAuthHeaders(url, method, requestBody)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	if a.UserAgent != "" {
		req.Header.Set("User-Agent", a.UserAgent)
	}

	log.Printf("[DEBUG] %s %s", method, url)
	resp, err := client.Do(req)

	if err != nil {
		return nil, resp, err
	}

	if resp.StatusCode != 200 {
		errorBodyBytes, _ := ioutil.ReadAll(resp.Body)
		return nil, resp, fmt.Errorf("non-200 status code from API, got status %d: %s", resp.StatusCode, string(errorBodyBytes))
	}

	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	return bodyBytes, resp, nil
}

func (a Auth) buildAuthHeaders(url string, method string, requestBody string) map[string]string {
	splitURL := strings.SplitN(url, "?", 2)

	authURL := splitURL[0]
	var queryString string
	if len(splitURL) == 2 {
		queryString = splitURL[1]
	} else {
		queryString = ""
	}
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1000000, 10)

	data := method + authURL + queryString + timestamp + requestBody

	decodedAPIKey, _ := hex.DecodeString(a.APIKey)

	tokenHmac := hmac.New(sha256.New, []byte(decodedAPIKey))
	tokenHmac.Write([]byte(data))

	token := hex.EncodeToString(tokenHmac.Sum(nil))

	return map[string]string{
		"X-LLNW-Security-Principal": a.APIUser,
		"X-LLNW-Security-Timestamp": timestamp,
		"X-LLNW-Security-Token":     token,
	}
}
