package hankoApiClient

import (
	"bytes"
	"encoding/json"
	"github.com/google/go-querystring/query"
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
	return "/" + client.apiVersion + "/webauthn"
}

func (client *HankoApiClient) GetWebAuthnRegistrationUrl() (url string) {
	return client.GetWebAuthnUrl() + "/registration"
}

func (client *HankoApiClient) GetWebAuthnAuthenticationUrl() (url string) {
	return client.GetWebAuthnUrl() + "/authentication"
}

func (client *HankoApiClient) GetWebAuthnTransactionUrl() (url string) {
	return client.GetWebAuthnUrl() + "/transaction"
}

func (client *HankoApiClient) GetWebAuthnCredentialsUrl() (url string) {
	return client.GetWebAuthnUrl() + "/credentials"
}

// WEBAUTHN ------------------------------------------------------------------------------------------------------------

// InitWebAuthnRegistration initiates the Registration of an Authenticator. Pass the result from the Hanko API to the
// WebAuthn API of the Browser to get it signed. The result has to be send back with FinalizeWebauthnOperation to finalize
// the Registration Flow.
func (client *HankoApiClient) InitWebAuthnRegistration(request *WebAuthnRegistrationInitializationRequest) (response *WebAuthnRegistrationInitializationResponse, err error) {
	response = &WebAuthnRegistrationInitializationResponse{}
	path := client.GetWebAuthnRegistrationUrl() + "/initialize"
	err = client.Request(http.MethodPost, path, request, response)
	return response, err
}

// FinalizeWebAuthnRegistration Is the last step to either Register or Authenticate an WebAuthn Request. Pass the result of
// the WebAuthn API call of the Browser to this method.
func (client *HankoApiClient) FinalizeWebAuthnRegistration(request *WebAuthnRegistrationFinalizationRequest) (response *WebAuthnRegistrationFinalizationResponse, err error) {
	response = &WebAuthnRegistrationFinalizationResponse{}
	path := client.GetWebAuthnRegistrationUrl() + "/finalize"
	err = client.Request(http.MethodPost, path, request, response)
	return response, err
}

// InitWebAuthnAuthentication initiates the Authentication Flow. Pass the challenge from the Hanko API to the
// WebAuthn API of the Browser to get it signed. The result has to be send back with FinalizeWebauthnOperation to finalize
// the Registration Flow.
func (client *HankoApiClient) InitWebAuthnAuthentication(request *WebAuthnAuthenticationInitializationRequest) (response *WebAuthnAuthenticationInitializationResponse, err error) {
	response = &WebAuthnAuthenticationInitializationResponse{}
	path := client.GetWebAuthnAuthenticationUrl() + "/initialize"
	err = client.Request(http.MethodPost, path, request, response)
	return response, err
}

// FinalizeWebAuthnRegistration Is the last step to either Register or Authenticate an WebAuthn Request. Pass the result of
// the WebAuthn API call of the Browser to this method.
func (client *HankoApiClient) FinalizeWebAuthnAuthentication(request *WebAuthnAuthenticationFinalizationRequest) (response *WebAuthnAuthenticationFinalizationResponse, err error) {
	response = &WebAuthnAuthenticationFinalizationResponse{}
	path := client.GetWebAuthnAuthenticationUrl() + "/finalize"
	err = client.Request(http.MethodPost, path, request, response)
	return response, err
}

// InitWebAuthnTransaction initiates the Authentication Flow. Pass the challenge from the Hanko API to the
// WebAuthn API of the Browser to get it signed. The result has to be send back with FinalizeWebauthnOperation to finalize
// the Registration Flow.
func (client *HankoApiClient) InitWebAuthnTransaction(request *WebAuthnTransactionInitializationRequest) (response *WebAuthnTransactionInitializationResponse, err error) {
	response = &WebAuthnTransactionInitializationResponse{}
	path := client.GetWebAuthnTransactionUrl() + "/initialize"
	err = client.Request(http.MethodPost, path, request, response)
	return response, err
}

// FinalizeWebAuthnRegistration Is the last step to either Register or Authenticate an WebAuthn Request. Pass the result of
// the WebAuthn API call of the Browser to this method.
func (client *HankoApiClient) FinalizeWebAuthnTransaction(request *WebAuthnTransactionFinalizationRequest) (response *WebAuthnTransactionFinalizationResponse, err error) {
	response = &WebAuthnTransactionFinalizationResponse{}
	path := client.GetWebAuthnTransactionUrl() + "/finalize"
	err = client.Request(http.MethodPost, path, request, response)
	return response, err
}

func (client *HankoApiClient) ListWebAuthnCredentials(q *WebAuthnCredentialQuery) (response *[]WebAuthnCredential, err error) {
	response = &[]WebAuthnCredential{}
	values, err := query.Values(q)
	if err != nil {
		return nil, err
	}
	path := client.GetWebAuthnTransactionUrl() + "?" + values.Encode()
	err = client.Request(http.MethodGet, path, nil, response)
	return response, err
}

func (client *HankoApiClient) GetWebAuthnCredential(id string) (response *WebAuthnCredential, err error) {
	response = &WebAuthnCredential{}
	path := client.GetWebAuthnTransactionUrl() + "/" + id
	err = client.Request(http.MethodGet, path, nil, response)
	return response, err
}

func (client *HankoApiClient) DeleteWebAuthnCredential(id string) (err error) {
	path := client.GetWebAuthnTransactionUrl() + "/" + id
	return client.Request(http.MethodDelete, path, nil, nil)
}

func (client *HankoApiClient) UpdateWebAuthnCredential(id string, request *WebAuthnCredentialUpdateRequest) (response *WebAuthnCredential, err error) {
	response = &WebAuthnCredential{}
	path := client.GetWebAuthnTransactionUrl() + "/" + id
	err = client.Request(http.MethodPut, path, request, response)
	return response, err
}

func (client *HankoApiClient) Request(method string, path string, request interface{}, response interface{}) (err error) {
	buf := new(bytes.Buffer)
	if request != nil {
		err = json.NewEncoder(buf).Encode(request)
		if err != nil {
			return errors.Wrap(err, "failed to encode the Request")
		}
	}

	req, err := http.NewRequest(method, client.baseUrl+path, buf)
	if err != nil {
		return errors.Wrapf(err, "failed to create Request %s %s", method, client.baseUrl+path)
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
		return errors.Wrapf(err, "could not do Request: %s %s%s", method, client.baseUrl, path)
	}

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("Request (%s %s%s) failed, got: %s", method, client.baseUrl, path, resp.Status)
	}

	if response != nil {
		body, err := ioutil.ReadAll(resp.Body)
		log.Printf("api raw response: %s", string(body))
		err = json.Unmarshal(body, response)
		if err != nil {
			return errors.Wrap(err, "failed to decode the response")
		}
	}

	return nil
}
