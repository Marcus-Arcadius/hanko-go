package hankoApiClient

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

// HankoApiClient Provides Methods for interacting with the Hanko API
type HankoApiClient struct {
	baseUrl    string
	secret     string
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

// WEBAUTHN ------------------------------------------------------------------------------------------------------------

// InitWebauthnRegistration initiates the Registration of an Authenticator. Pass the result from the Hanko API to the
// WebAuthn API of the Browser to get it signed. The result has to be send back with FinalizeWebauthnOperation to finalize
// the Registration Flow.
func (client *HankoApiClient) InitWebauthnRegistration(userId string, userName string) (*Response, error) {
	req := &Request{
		Operation: REG,
		Username:  userName,
		UserId:    userId,
		AuthenticatorSelectionCriteria: &AuthenticatorSelectionCriteria{
			UserVerification:        "preferred",
			AuthenticatorAttachment: "platform",
		},
	}
	return client.InitOperation("/"+client.apiVersion+"/webauthn/requests", req)
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
	return client.InitOperation("/"+client.apiVersion+"/webauthn/requests", req)
}

// InitWebAuthnDeRegistration de-registers an Authenticator Device. This Operation doesn't need to be finalized.
func (client *HankoApiClient) InitWebAuthnDeRegistration(userId string, userName string) (*Response, error) {
	req := &Request{
		Operation: DEREG,
		Username:  userName,
		UserId:    userId,
	}
	return client.InitOperation("/"+client.apiVersion+"/webauthn/requests", req)
}

// FinalizeWebAuthnOperation Is the last step to either Register or Authenticate an WebAuthn Request. Pass the result of
// the WebAuthn API call of the Browser to this method.
func (client *HankoApiClient) FinalizeWebAuthnOperation(requestId string, request *HankoCredentialRequest) (*Response, error) {
	return client.FinalizeOperation("/"+client.apiVersion+"/webauthn/requests", requestId, request)
}

// GetWebauthnRequestStatus Returns a status Response of a running request
func (client *HankoApiClient) GetWebauthnRequestStatus(requestId string) (*Response, error) {
	return client.GetRequestStatus("/"+client.apiVersion+"/webauthn/requests", requestId)
}

// UAF -----------------------------------------------------------------------------------------------------------------

func (client *HankoApiClient) InitUafRegistration(userId string, userName string) (*Response, error) {
	req := &Request{
		Operation: REG,
		Username:  userName,
		UserId:    userId,
	}
	return client.InitOperation("/"+client.apiVersion+"/uaf/requests", req)
}

func (client *HankoApiClient) InitUafAuthentication(userId string, userName string) (*Response, error) {
	req := &Request{
		Operation: AUTH,
		Username:  userName,
		UserId:    userId,
	}
	return client.InitOperation("/"+client.apiVersion+"/uaf/requests", req)
}

func (client *HankoApiClient) InitUafDeRegistration(userId string, userName string) (*Response, error) {
	req := &Request{
		Operation: DEREG,
		Username:  userName,
		UserId:    userId,
	}
	return client.InitOperation("/"+client.apiVersion+"/uaf/requests", req)
}

func (client *HankoApiClient) GetUafRequestStatus(requestId string) (*Response, error) {
	return client.GetRequestStatus("/"+client.apiVersion+"/uaf/requests", requestId)
}

// GENERIC -------------------------------------------------------------------------------------------------------------

func (client *HankoApiClient) InitOperation(path string, request *Request) (*Response, error) {
	return client.Request(http.MethodPost, path, request)
}

func (client *HankoApiClient) FinalizeOperation(path string, requestId string, request *HankoCredentialRequest) (*Response, error) {
	return client.Request(http.MethodPut, path+"/"+requestId, request)
}

func (client *HankoApiClient) GetRequestStatus(path string, requestId string) (*Response, error) {
	return client.Request(http.MethodGet, path+"/"+requestId, nil)
}

// Request does an AUTH/REG/DEREG based Request to the Hanko API
func (client *HankoApiClient) Request(method string, path string, request interface{}) (*Response, error) {
	resp, err := client.doRequest(method, path, request)

	apiResp := &Response{}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(apiResp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode the response")
	}

	return apiResp, nil
}

// doRequest does a generic Request to the Hanko API
func (client *HankoApiClient) doRequest(method string, path string, request interface{}) (response *http.Response, error error) {
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
	return resp, nil
}
