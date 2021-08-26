// Package webauthn contains webauthn-specific functions to establish communication with the Hanko Authentication API
package webauthn

import (
	"fmt"
	"github.com/google/go-querystring/query"
	hankoClient "github.com/teamhanko/hanko-go/client"
	"net/http"
)

type urlPath string

const (
	pathWebauthnBase             urlPath = "webauthn"
	pathRegistrationInitialize   urlPath = "registration/initialize"
	pathRegistrationFinalize     urlPath = "registration/finalize"
	pathAuthenticationInitialize urlPath = "authentication/initialize"
	pathAuthenticationFinalize   urlPath = "authentication/finalize"
	pathTransactionInitialize    urlPath = "transaction/initialize"
	pathTransactionFinalize      urlPath = "transaction/finalize"
	pathCredentials              urlPath = "credentials"
)

// Client wraps a basic client.Client and provides methods for registration, authentication and webauthn
// credential management (i.e. credential retrieval, update, and deletion).
type Client struct {
	client *hankoClient.Client // the base client.Client to be extended by this package
}

// NewClient creates a new webauthn.Client. Provide the baseUrl of the Hanko Authentication API server and your API
// secret in order to make authenticated requests to the API.
//
// You can obtain API keys and the API's baseUrl through the Hanko Console (https://console.hanko.io).
// To create API keys (an API Key ID, API secret pair) in the console, log in, select your relying party and go to
// "General Settings" -> "ApiKeys" -> "Add new". Make sure you provide an API secret in when constructing the Client
// via NewClient() accordingly.
//
// Note: It is recommended to use HMAC authorization. See the Client.WithHmac option for more details.
func NewClient(baseUrl string, secret string) *Client {
	return &Client{client: hankoClient.NewClient(baseUrl, secret)}
}

// getUrl constructs an returns a full API WebAuthn request URL, e.g. "https://{baseUrl}/{apiVersion}/webauthn/{urlPath}"
// using a given urlPath.
func (c *Client) getUrl(p urlPath) string {
	return fmt.Sprintf("%s/%s/%s", c.client.GetUrl(), pathWebauthnBase, p)
}

// InitializeRegistration initializes the registration of a new credential using a RegistrationInitializationRequest.
// On successful initialization, the Hanko Authentication API returns a RegistrationInitializationResponse. Send
// the response to your client application in order to pass it to the browser's WebAuthn API's
// navigator.credentials.create() function.
func (c *Client) InitializeRegistration(requestBody *RegistrationInitializationRequest) (response *RegistrationInitializationResponse, err *hankoClient.ApiError) {
	response = &RegistrationInitializationResponse{}
	requestUrl := c.getUrl(pathRegistrationInitialize)
	err = c.client.Request("initialize webauthn registration", http.MethodPost, requestUrl, requestBody, response)
	return response, err
}

// FinalizeRegistration finalizes the registration request initiated by the InitializeRegistration method. Provide a
// RegistrationFinalizationRequest which represents the result of calling the browser's WebAuthn API's
// navigator.credentials.create() function.
func (c *Client) FinalizeRegistration(requestBody *RegistrationFinalizationRequest) (response *RegistrationFinalizationResponse, err *hankoClient.ApiError) {
	response = &RegistrationFinalizationResponse{}
	requestUrl := c.getUrl(pathRegistrationFinalize)
	err = c.client.Request("finalize webauthn registration", http.MethodPost, requestUrl, requestBody, response)
	return response, err
}

// InitializeAuthentication initializes an authentication with a registered credential using an
// AuthenticationInitializationRequest. On successful initialization, the Hanko Authentication API returns a
// AuthenticationInitializationResponse. Send the response to your client application in order to pass it to the
// browser's WebAuthn API's navigator.credentials.get() function.
func (c *Client) InitializeAuthentication(requestBody *AuthenticationInitializationRequest) (response *AuthenticationInitializationResponse, err *hankoClient.ApiError) {
	response = &AuthenticationInitializationResponse{}
	requestUrl := c.getUrl(pathAuthenticationInitialize)
	err = c.client.Request("initialize webauthn authentication", http.MethodPost, requestUrl, requestBody, response)
	return response, err
}

// FinalizeAuthentication finalizes the authentication request initiated by the InitializeAuthentication method. Provide
// a AuthenticationFinalizationRequest which represents the result of calling the browser's WebAuthn API's
// navigator.credentials.get() function.
func (c *Client) FinalizeAuthentication(requestBody *AuthenticationFinalizationRequest) (response *AuthenticationFinalizationResponse, err *hankoClient.ApiError) {
	response = &AuthenticationFinalizationResponse{}
	requestUrl := c.getUrl(pathAuthenticationFinalize)
	err = c.client.Request("finalize webauthn authentication", http.MethodPost, requestUrl, requestBody, response)
	return response, err
}

// InitializeTransaction initiates a transaction. A transaction operation is analogous to the authentication operation,
// with the main difference being that a transaction context must be provided in the form of a string. This value will
// become part of the challenge an authenticator signs over during the operation.
//
// Initialize a transaction using a TransactionInitializationRequest. On successful initialization, the Hanko
// Authentication API returns a TransactionInitializationResponse. Send the response to your client application in order
// to pass it to the browser's WebAuthn API's navigator.credentials.get() function.
func (c *Client) InitializeTransaction(requestBody *TransactionInitializationRequest) (response *TransactionInitializationResponse, err *hankoClient.ApiError) {
	response = &TransactionInitializationResponse{}
	requestUrl := c.getUrl(pathTransactionInitialize)
	err = c.client.Request("initialize webauthn transaction", http.MethodPost, requestUrl, requestBody, response)
	return response, err
}

// FinalizeTransaction finalizes the transaction request initiated by the InitializeTransaction method. Provide
// a TransactionFinalizationRequest which represents the result of calling of the browser's WebAuthn API's
// navigator.credentials.get() function.
func (c *Client) FinalizeTransaction(requestBody *TransactionFinalizationRequest) (response *TransactionFinalizationResponse, err *hankoClient.ApiError) {
	response = &TransactionFinalizationResponse{}
	requestUrl := c.getUrl(pathTransactionFinalize)
	err = c.client.Request("finalize webauthn transaction", http.MethodPost, requestUrl, requestBody, response)
	return response, err
}

// ListCredentials returns a list of Credential. Filter by userId and paginate results using a CredentialQuery.
// The value for PageSize defaults to 10 and the value for Page to 1.
func (c *Client) ListCredentials(credentialQuery *CredentialQuery) (response *[]Credential, err *hankoClient.ApiError) {
	response = &[]Credential{}
	requestUrl := c.getUrl(pathCredentials)
	values, _ := query.Values(credentialQuery)
	if values != nil {
		requestUrl += "?" + values.Encode()
	}
	err = c.client.Request("list webauthn credentials", http.MethodGet, requestUrl, nil, response)
	return response, err
}

// GetCredential returns the Credential with the specified credentialId.
func (c *Client) GetCredential(credentialId string) (response *Credential, err *hankoClient.ApiError) {
	response = &Credential{}
	requestUrl := fmt.Sprintf("%s/%s", c.getUrl(pathCredentials), credentialId)
	err = c.client.Request("get webauthn credential", http.MethodGet, requestUrl, nil, response)
	return response, err
}

// DeleteCredential deletes the Credential with the specified credentialId.
func (c *Client) DeleteCredential(credentialId string) (err *hankoClient.ApiError) {
	requestUrl := fmt.Sprintf("%s/%s", c.getUrl(pathCredentials), credentialId)
	return c.client.Request("delete webauthn credential", http.MethodDelete, requestUrl, nil, nil)
}

// UpdateCredential updates the Credential with the specified credentialId. Provide a CredentialUpdateRequest with the
// updated data. Currently, you can only update the name of a Credential.
func (c *Client) UpdateCredential(credentialId string, requestBody *CredentialUpdateRequest) (response *Credential, err *hankoClient.ApiError) {
	response = &Credential{}
	requestUrl := fmt.Sprintf("%s/%s", c.getUrl(pathCredentials), credentialId)
	err = c.client.Request("update webauthn credential", http.MethodPut, requestUrl, requestBody, response)
	return response, err
}
