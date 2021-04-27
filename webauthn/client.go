package webauthn

import (
	"fmt"
	"github.com/google/go-querystring/query"
	hankoClient "github.com/teamhanko/hanko-sdk-golang/client"
	"net/http"
)

type path string

const (
	pathWebauthnBase             path = "webauthn"
	pathRegistrationInitialize   path = "registration/initialize"
	pathRegistrationFinalize     path = "registration/finalize"
	pathAuthenticationInitialize path = "authentication/initialize"
	pathAuthenticationFinalize   path = "authentication/finalize"
	pathTransactionInitialize    path = "transaction/initialize"
	pathTransactionFinalize      path = "transaction/finalize"
	pathCredentials              path = "credentials"
)

type Client struct {
	client *hankoClient.Client
}

func NewClient(baseUrl string, secret string) *Client {
	return &Client{client: hankoClient.NewClient(baseUrl, secret)}
}

func (c *Client) getUrl(path path) string {
	return fmt.Sprintf("%s/%s/%s", c.client.GetUrl(), pathWebauthnBase, path)
}

func (c *Client) InitializeRegistration(requestBody *RegistrationInitializationRequest) (response *RegistrationInitializationResponse, err *hankoClient.ApiError) {
	response = &RegistrationInitializationResponse{}
	requestUrl := c.getUrl(pathRegistrationInitialize)
	err = c.client.Request("initialize webauthn registration", http.MethodPost, requestUrl, requestBody, response)
	return response, err
}

func (c *Client) FinalizeRegistration(requestBody *RegistrationFinalizationRequest) (response *RegistrationFinalizationResponse, err *hankoClient.ApiError) {
	response = &RegistrationFinalizationResponse{}
	requestUrl := c.getUrl(pathRegistrationFinalize)
	err = c.client.Request("finalize webauthn registration", http.MethodPost, requestUrl, requestBody, response)
	return response, err
}

func (c *Client) InitializeAuthentication(requestBody *AuthenticationInitializationRequest) (response *AuthenticationInitializationResponse, err *hankoClient.ApiError) {
	response = &AuthenticationInitializationResponse{}
	requestUrl := c.getUrl(pathAuthenticationInitialize)
	err = c.client.Request("initialize webauthn authentication", http.MethodPost, requestUrl, requestBody, response)
	return response, err
}

func (c *Client) FinalizeAuthentication(requestBody *AuthenticationFinalizationRequest) (response *AuthenticationFinalizationResponse, err *hankoClient.ApiError) {
	response = &AuthenticationFinalizationResponse{}
	requestUrl := c.getUrl(pathAuthenticationFinalize)
	err = c.client.Request("finalize webauthn authentication", http.MethodPost, requestUrl, requestBody, response)
	return response, err
}

func (c *Client) InitializeTransaction(requestBody *TransactionInitializationRequest) (response *TransactionInitializationResponse, err *hankoClient.ApiError) {
	response = &TransactionInitializationResponse{}
	requestUrl := c.getUrl(pathTransactionInitialize)
	err = c.client.Request("initialize webauthn transaction", http.MethodPost, requestUrl, requestBody, response)
	return response, err
}

func (c *Client) FinalizeTransaction(requestBody *TransactionFinalizationRequest) (response *TransactionFinalizationResponse, err *hankoClient.ApiError) {
	response = &TransactionFinalizationResponse{}
	requestUrl := c.getUrl(pathTransactionFinalize)
	err = c.client.Request("finalize webauthn transaction", http.MethodPost, requestUrl, requestBody, response)
	return response, err
}

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

func (c *Client) GetCredential(credentialId string) (response *Credential, err *hankoClient.ApiError) {
	response = &Credential{}
	requestUrl := fmt.Sprintf("%s/%s", c.getUrl(pathCredentials), credentialId)
	err = c.client.Request("get webauthn credential", http.MethodGet, requestUrl, nil, response)
	return response, err
}

func (c *Client) DeleteCredential(credentialId string) (err *hankoClient.ApiError) {
	requestUrl := fmt.Sprintf("%s/%s", c.getUrl(pathCredentials), credentialId)
	return c.client.Request("delete webauthn credential", http.MethodDelete, requestUrl, nil, nil)
}

func (c *Client) UpdateCredential(credentialId string, requestBody *CredentialUpdateRequest) (response *Credential, err *hankoClient.ApiError) {
	response = &Credential{}
	requestUrl := fmt.Sprintf("%s/%s", c.getUrl(pathCredentials), credentialId)
	err = c.client.Request("update webauthn credential", http.MethodPut, requestUrl, requestBody, response)
	return response, err
}
