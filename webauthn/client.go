package webauthn

import (
	"fmt"
	"github.com/google/go-querystring/query"
	hankoClient "github.com/teamhanko/hanko-sdk-golang/client"
	"net/http"
)

type Client struct {
	*hankoClient.Client
}

// Returns a Client give it the base url e.g. https://api.hanko.io and your API Secret
func NewClient(baseUrl string, secret string, opts ...hankoClient.Option) (client *Client) {
	return &Client{Client: hankoClient.NewClient(baseUrl, secret, opts...)}
}

func (c *Client) GetUrl() (url string) {
	return fmt.Sprintf("%s/%s/webauthn", c.BaseUrl, c.ApiVersion)
}

func (c *Client) GetRegistrationUrl() (url string) {
	return c.GetUrl() + "/registration"
}

func (c *Client) GetAuthenticationUrl() (url string) {
	return c.GetUrl() + "/authentication"
}

func (c *Client) GetTransactionUrl() (url string) {
	return c.GetUrl() + "/transaction"
}

func (c *Client) GetCredentialsUrl() (url string) {
	return c.GetUrl() + "/credentials"
}

// WEBAUTHN ------------------------------------------------------------------------------------------------------------

// InitializeRegistration initiates the Registration of an Authenticator. Pass the result from the Hanko API to the
// WebAuthn API of the Browser to get it signed. The result has to be send back with FinalizeWebauthnOperation to finalize
// the Registration Flow.
func (c *Client) InitializeRegistration(requestBody *RegistrationInitializationRequest) (response *RegistrationInitializationResponse, err error) {
	response = &RegistrationInitializationResponse{}
	requestUrl := c.GetRegistrationUrl() + "/initialize"
	err = c.Request("initialize webauthn registration", http.MethodPost, requestUrl, requestBody, response)
	return response, err
}

// FinalizeRegistration Is the last step to either Register or Authenticate an WebAuthn do. Pass the result of
// the WebAuthn API call of the Browser to this method.
func (c *Client) FinalizeRegistration(requestBody *RegistrationFinalizationRequest) (response *RegistrationFinalizationResponse, err error) {
	response = &RegistrationFinalizationResponse{}
	requestUrl := c.GetRegistrationUrl() + "/finalize"
	err = c.Request("finalize webauthn registration", http.MethodPost, requestUrl, requestBody, response)
	return response, err
}

// InitializeAuthentication initiates the Authentication Flow. Pass the challenge from the Hanko API to the
// WebAuthn API of the Browser to get it signed. The result has to be send back with FinalizeWebauthnOperation to finalize
// the Registration Flow.
func (c *Client) InitializeAuthentication(requestBody *AuthenticationInitializationRequest) (response *AuthenticationInitializationResponse, err error) {
	response = &AuthenticationInitializationResponse{}
	requestUrl := c.GetAuthenticationUrl() + "/initialize"
	err = c.Request("initialize webauthn authentication", http.MethodPost, requestUrl, requestBody, response)
	return response, err
}

// FinalizeRegistration Is the last step to either Register or Authenticate an WebAuthn do. Pass the result of
// the WebAuthn API call of the Browser to this method.
func (c *Client) FinalizeAuthentication(requestBody *AuthenticationFinalizationRequest) (response *AuthenticationFinalizationResponse, err error) {
	response = &AuthenticationFinalizationResponse{}
	requestUrl := c.GetAuthenticationUrl() + "/finalize"
	err = c.Request("finalize webauthn authentication", http.MethodPost, requestUrl, requestBody, response)
	return response, err
}

// InitializeTransaction initiates the Authentication Flow. Pass the challenge from the Hanko API to the
// WebAuthn API of the Browser to get it signed. The result has to be send back with FinalizeWebauthnOperation to finalize
// the Registration Flow.
func (c *Client) InitializeTransaction(requestBody *TransactionInitializationRequest) (response *TransactionInitializationResponse, err error) {
	response = &TransactionInitializationResponse{}
	requestUrl := c.GetTransactionUrl() + "/initialize"
	err = c.Request("initialize webauthn transaction", http.MethodPost, requestUrl, requestBody, response)
	return response, err
}

// FinalizeRegistration Is the last step to either Register or Authenticate an WebAuthn do. Pass the result of
// the WebAuthn API call of the Browser to this method.
func (c *Client) FinalizeTransaction(requestBody *TransactionFinalizationRequest) (response *TransactionFinalizationResponse, err error) {
	response = &TransactionFinalizationResponse{}
	requestUrl := c.GetTransactionUrl() + "/finalize"
	err = c.Request("finalize webauthn transaction", http.MethodPost, requestUrl, requestBody, response)
	return response, err
}

// ListCredentials
// TODO: Docs (pagination defaults...)
func (c *Client) ListCredentials(credentialQuery *CredentialQuery) (response *[]Credential, err error) {
	response = &[]Credential{}
	requestUrl := c.GetCredentialsUrl()

	values, err := query.Values(credentialQuery)
	if err == nil {
		requestUrl += "?" + values.Encode()
	}

	err = c.Request("list webauthn credentials", http.MethodGet, requestUrl, nil, response)
	return response, err
}

// TODO: Docs
func (c *Client) GetCredential(credentialId string) (response *Credential, err error) {
	response = &Credential{}
	requestUrl := fmt.Sprintf("%s/%s", c.GetCredentialsUrl(), credentialId)
	err = c.Request("get webauthn credential", http.MethodGet, requestUrl, nil, response)
	return response, err
}

// TODO: Docs
func (c *Client) DeleteCredential(credentialId string) (err error) {
	requestUrl := fmt.Sprintf("%s/%s", c.GetCredentialsUrl(), credentialId)
	return c.Request("delete webauthn credential", http.MethodDelete, requestUrl, nil, nil)
}

// TODO: Docs
func (c *Client) UpdateCredential(credentialId string, requestBody *CredentialUpdateRequest) (response *Credential, err error) {
	response = &Credential{}
	requestUrl := fmt.Sprintf("%s/%s", c.GetCredentialsUrl(), credentialId)
	err = c.Request("update webauthn credential", http.MethodPut, requestUrl, requestBody, response)
	return response, err
}
