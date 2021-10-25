// Package passlink provides core definitions supporting passwordless authentication through the Hanko
// Authentication API using Passlinks.
package passlink

import (
	"fmt"
	hankoClient "github.com/teamhanko/hanko-go/client"
	"net/http"
)

type urlPath string

const (
	pathPasslinkBase       urlPath = "passlink"
	pathPasslinkInitialize urlPath = "initialize"
	pathPasslinkFinalize   urlPath = "%s/finalize"
)

// Client wraps a basic client.Client and provides methods for initializing and finalizing Passlink-based authentication
// flows with the Hanko Authentication API.
type Client struct {
	client *hankoClient.Client
}

// NewClient creates a new passlink.Client. Provide the baseUrl of the Hanko Authentication API server and your API
// secret in order to make authenticated requests to the API.
//
// You can obtain API keys and the API's baseUrl through the Hanko Console (https://console.hanko.io).
// To create API keys (an API Key ID, API secret pair) in the console, log in, select your relying party and go to
// "General Settings" -> "ApiKeys" -> "Add new". Make sure you provide an API secret when constructing the Client
// via NewClient() accordingly.
//
// Note: It is recommended to use HMAC authorization. See the Client.WithHmac option for more details.
func NewClient(baseUrl string, secret string) *Client {
	return &Client{client: hankoClient.NewClient(baseUrl, secret)}
}

// getUrl constructs and returns a full Passlink API request URL, e.g. "https://{baseUrl}/{apiVersion}/passlink/{urlPath}"
// using a given urlPath.
func (c *Client) getUrl(p urlPath) string {
	return fmt.Sprintf("%s/%s/%s", c.client.GetUrl(), pathPasslinkBase, p)
}

// InitializePasslink triggers the creation of a new Passlink using a LinkRequest.
// On successful initialization, the Hanko Authentication API will send a message containing a link to the recipient
// specified in the requestBody LinkRequest and returns a representation of the created Passlink as a Link.
func (c *Client) InitializePasslink(requestBody *LinkRequest) (response *Link, err *hankoClient.ApiError) {
	response = &Link{}
	requestUrl := c.getUrl(pathPasslinkInitialize)
	err = c.client.Request("initialize passlink", http.MethodPost, requestUrl, requestBody, response)
	return response, err
}

// FinalizePasslink completes a Passlink-based authentication flow using the linkId of a Passlink previously
// confirmed ("clicked") by the user. Only "confirmed" Passlinks can be finalized. Any attempt to finalize a "pending"
// Passlink (i.e. initialized but not confirmed by the user) will result in an API error response.
// On successful finalization the Hanko Authentication API will return a representation of the finalized Passlink as
// a Link. This response indicates that the status of the Passlink is "finished" and the Passlink can no longer be used
// to authenticate.
func (c *Client) FinalizePasslink(linkId string) (response *Link, err *hankoClient.ApiError) {
	response = &Link{}
	requestUrl := fmt.Sprintf(c.getUrl(pathPasslinkFinalize), linkId)
	err = c.client.Request("finalize passlink", http.MethodPatch, requestUrl, nil, response)
	return response, err
}

