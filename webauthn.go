package hankoApiClient

import (
	"fmt"
	"github.com/google/go-querystring/query"
	"net/http"
)

func (c *HankoApiClient) GetWebAuthnUrl() (url string) {
	return fmt.Sprintf("%s/%s/webauthn", c.baseUrl, c.apiVersion)
}

func (c *HankoApiClient) GetWebAuthnRegistrationUrl() (url string) {
	return c.GetWebAuthnUrl() + "/registration"
}

func (c *HankoApiClient) GetWebAuthnAuthenticationUrl() (url string) {
	return c.GetWebAuthnUrl() + "/authentication"
}

func (c *HankoApiClient) GetWebAuthnTransactionUrl() (url string) {
	return c.GetWebAuthnUrl() + "/transaction"
}

func (c *HankoApiClient) GetWebAuthnCredentialsUrl() (url string) {
	return c.GetWebAuthnUrl() + "/credentials"
}

// WEBAUTHN ------------------------------------------------------------------------------------------------------------

// InitializeWebAuthnRegistration initiates the Registration of an Authenticator. Pass the result from the Hanko API to the
// WebAuthn API of the Browser to get it signed. The result has to be send back with FinalizeWebauthnOperation to finalize
// the Registration Flow.
func (c *HankoApiClient) InitializeWebAuthnRegistration(requestBody *WebAuthnRegistrationInitializationRequest) (response *WebAuthnRegistrationInitializationResponse, err error) {
	response = &WebAuthnRegistrationInitializationResponse{}
	requestUrl := c.GetWebAuthnRegistrationUrl() + "/initialize"
	err = c.run("initialize webauthn registration", http.MethodPost, requestUrl, requestBody, response)
	return response, err
}

// FinalizeWebAuthnRegistration Is the last step to either Register or Authenticate an WebAuthn do. Pass the result of
// the WebAuthn API call of the Browser to this method.
func (c *HankoApiClient) FinalizeWebAuthnRegistration(requestBody *WebAuthnRegistrationFinalizationRequest) (response *WebAuthnRegistrationFinalizationResponse, err error) {
	response = &WebAuthnRegistrationFinalizationResponse{}
	requestUrl := c.GetWebAuthnRegistrationUrl() + "/finalize"
	err = c.run("finalize webauthn registration", http.MethodPost, requestUrl, requestBody, response)
	return response, err
}

// InitializeWebAuthnAuthentication initiates the Authentication Flow. Pass the challenge from the Hanko API to the
// WebAuthn API of the Browser to get it signed. The result has to be send back with FinalizeWebauthnOperation to finalize
// the Registration Flow.
func (c *HankoApiClient) InitializeWebAuthnAuthentication(requestBody *WebAuthnAuthenticationInitializationRequest) (response *WebAuthnAuthenticationInitializationResponse, err error) {
	response = &WebAuthnAuthenticationInitializationResponse{}
	requestUrl := c.GetWebAuthnAuthenticationUrl() + "/initialize"
	err = c.run("initialize webauthn authentication", http.MethodPost, requestUrl, requestBody, response)
	return response, err
}

// FinalizeWebAuthnRegistration Is the last step to either Register or Authenticate an WebAuthn do. Pass the result of
// the WebAuthn API call of the Browser to this method.
func (c *HankoApiClient) FinalizeWebAuthnAuthentication(requestBody *WebAuthnAuthenticationFinalizationRequest) (response *WebAuthnAuthenticationFinalizationResponse, err error) {
	response = &WebAuthnAuthenticationFinalizationResponse{}
	requestUrl := c.GetWebAuthnAuthenticationUrl() + "/finalize"
	err = c.run("finalize webauthn authentication", http.MethodPost, requestUrl, requestBody, response)
	return response, err
}

// InitializeWebAuthnTransaction initiates the Authentication Flow. Pass the challenge from the Hanko API to the
// WebAuthn API of the Browser to get it signed. The result has to be send back with FinalizeWebauthnOperation to finalize
// the Registration Flow.
func (c *HankoApiClient) InitializeWebAuthnTransaction(requestBody *WebAuthnTransactionInitializationRequest) (response *WebAuthnTransactionInitializationResponse, err error) {
	response = &WebAuthnTransactionInitializationResponse{}
	requestUrl := c.GetWebAuthnTransactionUrl() + "/initialize"
	err = c.run("initialize webauthn transaction", http.MethodPost, requestUrl, requestBody, response)
	return response, err
}


// FinalizeWebAuthnRegistration Is the last step to either Register or Authenticate an WebAuthn do. Pass the result of
// the WebAuthn API call of the Browser to this method.
func (c *HankoApiClient) FinalizeWebAuthnTransaction(requestBody *WebAuthnTransactionFinalizationRequest) (response *WebAuthnTransactionFinalizationResponse, err error) {
	response = &WebAuthnTransactionFinalizationResponse{}
	requestUrl := c.GetWebAuthnTransactionUrl() + "/finalize"
	err = c.run("finalize webauthn transaction", http.MethodPost, requestUrl, requestBody, response)
	return response, err
}

// ListWebAuthnCredentials
// TODO: Docs (pagination defaults...)
func (c *HankoApiClient) ListWebAuthnCredentials(credentialQuery *WebAuthnCredentialQuery) (response *[]WebauthnCredential, err error) {
	response = &[]WebauthnCredential{}
	requestUrl := c.GetWebAuthnCredentialsUrl()

	values, err := query.Values(credentialQuery)
	if err == nil {
		requestUrl += "?" + values.Encode()
	}

	err = c.run("list webauthn credentials", http.MethodGet, requestUrl, nil, response)
	return response, err
}

// TODO: Docs
func (c *HankoApiClient) GetWebAuthnCredential(credentialId string) (response *WebauthnCredential, err error) {
	response = &WebauthnCredential{}
	requestUrl := fmt.Sprintf("%s/%s", c.GetWebAuthnCredentialsUrl(), credentialId)
	err = c.run("get webauthn credential", http.MethodGet, requestUrl, nil, response)
	return response, err
}

// TODO: Docs
func (c *HankoApiClient) DeleteWebAuthnCredential(credentialId string) (err error) {
	requestUrl := fmt.Sprintf("%s/%s", c.GetWebAuthnCredentialsUrl(), credentialId)
	return c.run("delete webauthn credential", http.MethodDelete, requestUrl, nil, nil)
}

// TODO: Docs
func (c *HankoApiClient) UpdateWebAuthnCredential(credentialId string, requestBody *WebAuthnCredentialUpdateRequest) (response *WebauthnCredential, err error) {
	response = &WebauthnCredential{}
	requestUrl := fmt.Sprintf("%s/%s", c.GetWebAuthnCredentialsUrl(), credentialId)
	err = c.run("update webauthn credential", http.MethodPut, requestUrl, requestBody, response)
	return response, err
}
