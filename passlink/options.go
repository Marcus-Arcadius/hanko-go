package passlink

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// WithHmac sets the hmacApiKeyId. If set, the Client will use HMAC authorization when making a request to the Hanko
// Authentication API. This means an HMAC token will be calculated and used in the HTTP authorization header of each
// request.
//
// You can generate API Keys consisting of an API Key ID (here: hmacApiKeyId), API secret pair
// through the Hanko Console (https://console.hanko.io). To create API keys (an API Key ID, API secret pair) in the
// Console, log in, select your relying party and go to "General Settings" -> "ApiKeys" -> "Add new".
// Make sure you provide an API secret in when constructing the Client via NewClient() accordingly.
func (c *Client) WithHmac(hmacApiKeyId string) *Client {
	c.client.SetHmac(hmacApiKeyId)
	return c
}

// WithHttpClient allows you to set a custom http.Client. This can be useful, for example, if you want to configure
// an HTTP proxy.
func (c *Client) WithHttpClient(httpClient *http.Client) *Client {
	c.client.SetHttpClient(httpClient)
	return c
}

// WithLogger allows you to set your own custom logrus.Logger.
func (c *Client) WithLogger(logger *log.Logger) *Client {
	c.client.SetLogger(logger)
	return c
}

// WithoutLogs allows you to disable logging by setting an io.Writer that discards the logs.
func (c *Client) WithoutLogs() *Client {
	c.client.SetLogWriter(ioutil.Discard)
	return c
}

// WithLogLevel allows you to set the specified logrus.Level. The default level is logrus.InfoLevel.
func (c *Client) WithLogLevel(level log.Level) *Client {
	c.client.SetLogLevel(level)
	return c
}

// WithLogFormatter allows you to set the specified logrus.Formatter. The default formatter is logrus.JSONFormatter.
func (c *Client) WithLogFormatter(formatter log.Formatter) *Client {
	c.client.SetLogFormatter(formatter)
	return c
}
