package magic

import (
	"fmt"
	hankoClient "github.com/teamhanko/hanko-go/client"
	"net/http"
)

type urlPath string

const (
	pathMagicBase       urlPath = "magic_link"
	pathMagicInitialize urlPath = "initialize"
	pathMagicFinalize   urlPath = "%s/finalize"
)

type Client struct {
	client *hankoClient.Client
}

func NewClient(baseUrl string, secret string) *Client {
	return &Client{client: hankoClient.NewClient(baseUrl, secret)}
}

func (c *Client) getUrl(p urlPath) string {
	return fmt.Sprintf("%s/%s/%s", c.client.GetUrl(), pathMagicBase, p)
}

func (c *Client) InitializeMagicLink(requestBody *LinkRequest) (response *Link, err *hankoClient.ApiError) {
	response = &Link{}
	requestUrl := c.getUrl(pathMagicInitialize)
	err = c.client.Request("initialize magic link", http.MethodPost, requestUrl, requestBody, response)
	return response, err
}

func (c *Client) FinalizeMagicLink(linkId string) (response *Link, err *hankoClient.ApiError) {
	response = &Link{}
	requestUrl := fmt.Sprintf(c.getUrl(pathMagicFinalize), linkId)
	err = c.client.Request("finalize magic link", http.MethodPatch, requestUrl, nil, response)
	return response, err
}

