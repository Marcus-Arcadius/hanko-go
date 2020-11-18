package hankoApiClient

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
)

// HankoApiClient Provides Methods for interacting with the Hanko API
type HankoApiClient struct {
	baseUrl    string
	secret     string
	apiKeyId   string
	httpClient http.Client
	apiVersion string
}

// Returns a HankoApiClient give it the base url e.g. https://api.hanko.io and your API Secret
func NewHankoApiClient(baseUrl string, secret string) *HankoApiClient {
	return &HankoApiClient{
		baseUrl:    baseUrl,
		secret:     secret,
		apiVersion: "v1",
	}
}

// Returns new client with capabilities for calculating an HMAC
func NewHankoHmacClient(baseUrl string, secret string, apiKeyId string) *HankoApiClient {
	return &HankoApiClient{
		baseUrl:    baseUrl,
		secret:     secret,
		apiKeyId:   apiKeyId,
		apiVersion: "v1",
	}
}

func (client *HankoApiClient) GetWebAuthnUrl() (url string) {
	return "/" + client.apiVersion + "/webauthn/requests"
}

func (client *HankoApiClient) GetUafUrl() (url string) {
	return "/" + client.apiVersion + "/uaf/requests"
}

// WEBAUTHN ------------------------------------------------------------------------------------------------------------

// InitWebauthnRegistration initiates the Registration of an Authenticator. Pass the result from the Hanko API to the
// WebAuthn API of the Browser to get it signed. The result has to be send back with FinalizeWebauthnOperation to finalize
// the Registration Flow.
func (client *HankoApiClient) InitWebauthnRegistration(userId string, userName string) (*Response, error) {
	req := &Request{
		Operation: REG,
		Username:  userName,
		UserId:    userId,
	}
	return client.InitOperation(client.GetWebAuthnUrl(), req)
}

// InitWebAuthnAuthentication initiates the Authentication Flow. Pass the challenge from the Hanko API to the
// WebAuthn API of the Browser to get it signed. The result has to be send back with FinalizeWebauthnOperation to finalize
// the Registration Flow.
func (client *HankoApiClient) InitWebAuthnAuthentication(userId string, userName string) (*Response, error) {
	req := &Request{
		Operation: AUTH,
		Username:  userName,
		UserId:    userId,
	}
	return client.InitOperation(client.GetWebAuthnUrl(), req)
}

// InitWebAuthnDeRegistration de-registers an Authenticator Device. This Operation doesn't need to be finalized.
func (client *HankoApiClient) InitWebAuthnDeRegistration(userId string, userName string) (*Response, error) {
	req := &Request{
		Operation: DEREG,
		Username:  userName,
		UserId:    userId,
	}
	return client.InitOperation(client.GetWebAuthnUrl(), req)
}

// FinalizeWebAuthnOperation Is the last step to either Register or Authenticate an WebAuthn Request. Pass the result of
// the WebAuthn API call of the Browser to this method.
func (client *HankoApiClient) FinalizeWebAuthnOperation(requestId string, request *HankoWebAuthnCredentialRequest) (*Response, error) {
	return client.FinalizeOperation(client.GetWebAuthnUrl(), requestId, request)
}

// GetWebauthnRequestStatus Returns a status Response of a running request
func (client *HankoApiClient) GetWebauthnRequestStatus(requestId string) (*Response, error) {
	return client.GetRequestStatus(client.GetWebAuthnUrl(), requestId)
}

func (client *HankoApiClient) CancelWebAuthnRequest(requestId string) (*Response, error) {
	return client.CancelOperation(client.GetWebAuthnUrl(), requestId)
}

// UAF -----------------------------------------------------------------------------------------------------------------

func (client *HankoApiClient) InitUafRegistration(userId string, userName string) (*Response, error) {
	req := &Request{
		Operation: REG,
		Username:  userName,
		UserId:    userId,
	}
	return client.InitOperation(client.GetUafUrl(), req)
}

func (client *HankoApiClient) InitUafAuthentication(userId string, userName string) (*Response, error) {
	req := &Request{
		Operation: AUTH,
		Username:  userName,
		UserId:    userId,
	}
	return client.InitOperation(client.GetUafUrl(), req)
}

func (client *HankoApiClient) InitUafDeRegistration(userId string, userName string) (*Response, error) {
	req := &Request{
		Operation: DEREG,
		Username:  userName,
		UserId:    userId,
	}
	return client.InitOperation(client.GetUafUrl(), req)
}

func (client *HankoApiClient) FinalizeUafOperation(requestId string, request *HankoUafCredentialRequest) (*Response, error) {
	return client.FinalizeOperation(client.GetUafUrl(), requestId, request)
}

func (client *HankoApiClient) GetUafRequestStatus(requestId string) (*Response, error) {
	return client.GetRequestStatus(client.GetUafUrl(), requestId)
}

func (client *HankoApiClient) CancelUafRequest(requestId string) (*Response, error) {
	return client.CancelOperation(client.GetUafUrl(), requestId)
}

// GENERIC -------------------------------------------------------------------------------------------------------------

func (client *HankoApiClient) InitOperation(path string, request *Request) (*Response, error) {
	return client.Request(http.MethodPost, path, request)
}

func (client *HankoApiClient) FinalizeOperation(path string, requestId string, request interface{}) (*Response, error) {
	return client.Request(http.MethodPut, path+"/"+requestId, request)
}

func (client *HankoApiClient) GetRequestStatus(path string, requestId string) (*Response, error) {
	return client.Request(http.MethodGet, path+"/"+requestId, nil)
}

func(client *HankoApiClient) CancelOperation(path string, requestId string) (*Response, error) {
	return client.Request(http.MethodDelete, path+"/"+requestId, nil)
}

// Request does an AUTH/REG/DEREG based Request to the Hanko API
func (client *HankoApiClient) Request(method string, path string, request interface{}) (*Response, error) {
	resp, err := client.doRequest(method, path, request)
	if err != nil {
		return nil, err
	}

	apiResp := &Response{}
	body, err := ioutil.ReadAll(resp.Body)
	log.Printf("api raw response: %s", string(body))

	err = json.Unmarshal(body, apiResp)
	//dec := json.NewDecoder(resp.Body)
	//err = dec.Decode(apiResp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode the response")
	}

	return apiResp, nil
}

// doRequest does a generic Request to the Hanko API
func (client *HankoApiClient) doRequest(method string, path string, request interface{}) (*http.Response, error) {
	var err error

	buf := new(bytes.Buffer)
	if request != nil {
		err = json.NewEncoder(buf).Encode(request)
		if err != nil {
			return nil, errors.Wrap(err, "failed to encode the request")
		}
	}

	req, err := http.NewRequest(method, client.baseUrl+path, buf)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create request %s %s", method, client.baseUrl+path)
	}

	if client.apiKeyId != "" {
		hmac := CalculateHmac(&HmacMessageData{
			apiSecret:     client.secret,
			apiKeyId:      client.apiKeyId,
			requestMethod: method,
			requestPath:   path,
			requestBody:   buf.String(),
		})

		req.Header.Add("Authorization", "HANKO "+hmac)
	} else {
		req.Header.Add("Authorization", "secret "+client.secret)
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := client.httpClient.Do(req)

	if err != nil {
		return nil, errors.Wrapf(err, "could not do request: %s %s%s", method, client.baseUrl, path)
	}

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, "could not read error response body")
		}
		errorResponse := &HankoError{}
		err = json.Unmarshal(body, errorResponse)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode error response")
		}

		return nil, errorResponse
	}

	return resp, nil
}
