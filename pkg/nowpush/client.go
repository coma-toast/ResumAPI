package nowpush

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

// Client is the PrintifyAPI Client
type Client struct {
	client *http.Client
	Token  string
}

// Error is the error
type Error struct {
	Result       bool   `json:"result"`
	ErrorText    string `json:"error"`
	CallResponse *http.Response
}

// Result is the API Result
type Result struct {
	code string `json:"code"`
}

var baseURL = "https://www.api.nowpush.app/v3/"

func (e Error) Error() string {
	message := fmt.Sprintf("NowPush API Error: %s \n Status Code: %d", e.ErrorText, e.CallResponse.StatusCode)
	return message
}

func (c *Client) call(method string, destinationURL string, payload []byte, queryParams map[string]string, target interface{}) ([]byte, error) {
	destinationURL = baseURL + destinationURL
	if c.client == nil {
		c.client = &http.Client{}
	}

	var err error
	var request *http.Request

	if payload != nil {
		if method == http.MethodPost {
			request, err = http.NewRequest(method, destinationURL, bytes.NewBuffer(payload))
			request.Header.Set("Content-Type", "application/json")
		}
	} else if queryParams != nil {
		params := url.Values{}
		for param, value := range queryParams {
			params.Set(param, value)
			log.Debugf("%s: %s", param, value)
		}
		request, err = http.NewRequest(method, destinationURL+"?"+params.Encode(), nil)
		log.Debugf("Params: %s", params.Encode())
	} else {
		log.Debugf("No params or payload: %s", destinationURL)
		request, err = http.NewRequest(method, destinationURL, nil)
	}

	request.Header.Set("Authorization", "Bearer "+c.Token)
	if err != nil {
		return []byte{}, err
	}

	resp, err := c.client.Do(request)
	log.Debugf("Sending request to URL: %s", request.URL)
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	//TODO: this can all be one error function, take responseBody and do all the error checks
	errorTarget := Error{}
	if target != nil {
		err = json.Unmarshal(responseBody, &target)
		if err != nil {
			return responseBody, err
		}
	}

	if errorTarget.ErrorText != "" {
		errorTarget.CallResponse = resp
		return responseBody, errorTarget
	}
	// TODO: ^ to here
	if resp.StatusCode >= 400 {
		err := fmt.Errorf("HTTP Error: %d", resp.StatusCode)
		return responseBody, err
	}

	return responseBody, nil
}
