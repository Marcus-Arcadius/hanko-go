package webauthn

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func (c *Client) WithHmac(hmacApiKeyId string) *Client {
	c.client.SetHmac(hmacApiKeyId)
	return c
}

func (c *Client) WithHttpClient(httpClient *http.Client) *Client {
	c.client.SetHttpClient(httpClient)
	return c
}

func (c *Client) WithLogger(logger *log.Logger) *Client {
	c.client.SetLogger(logger)
	return c
}

func (c *Client) WithoutLogs() *Client {
	c.client.SetLogWriter(ioutil.Discard)
	return c
}

func (c *Client) WithLogLevel(level log.Level) *Client {
	c.client.SetLogLevel(level)
	return c
}

func (c *Client) WithLogFormatter(formatter log.Formatter) *Client {
	c.client.SetLogFormatter(formatter)
	return c
}
