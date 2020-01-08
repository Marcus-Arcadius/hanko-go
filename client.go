package hankoApiClient

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

type HankoApiClient struct {
	baseUrl    string
	secret     string
	httpClient http.Client
}

func NewHankoApiClient(baseUrl string, secret string) *HankoApiClient {
	return &HankoApiClient{
		baseUrl: baseUrl,
		secret:  secret,
	}
}

func (client *HankoApiClient) Request(method string, path string, request interface{}) (*Response, error) {
	buf := new(bytes.Buffer)
	if request != nil {
		err := json.NewEncoder(buf).Encode(request)
		if err != nil {
			return nil, errors.Wrap(err, "failed to encode the request")
		}
	}

	req, err := http.NewRequest(method, client.baseUrl+path, buf)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create request %s %s", method, client.baseUrl+path)
	}

	req.Header.Add("Authorization", "secret "+client.secret)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "could not do request")
	}

	apiResp := &Response{}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(apiResp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode the response")
	}

	return apiResp, nil
}
