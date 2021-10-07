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

type Client struct {
	client *hankoClient.Client
}

func NewClient(baseUrl string, secret string) *Client {
	return &Client{client: hankoClient.NewClient(baseUrl, secret)}
}

func (c *Client) getUrl(p urlPath) string {
	return fmt.Sprintf("%s/%s/%s", c.client.GetUrl(), pathPasslinkBase, p)
}

func (c *Client) InitializePasslink(requestBody *LinkRequest) (response *Link, err *hankoClient.ApiError) {
	response = &Link{}
	requestUrl := c.getUrl(pathPasslinkInitialize)
	err = c.client.Request("initialize passlink", http.MethodPost, requestUrl, requestBody, response)
	return response, err
}

func (c *Client) FinalizePasslink(linkId string) (response *Link, err *hankoClient.ApiError) {
	response = &Link{}
	requestUrl := fmt.Sprintf(c.getUrl(pathPasslinkFinalize), linkId)
	err = c.client.Request("finalize passlink", http.MethodPatch, requestUrl, nil, response)
	return response, err
}

