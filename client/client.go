package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// Client Provides Methods for interacting with the Hanko API
type Client struct {
	BaseUrl      string
	ApiVersion   string
	secret       string
	hmacApiKeyId string
	httpClient   *http.Client
	log          *log.Logger
}

type Option func(client *Client) *Client

func WithHmac(hmacApiKeyId string) Option {
	return func(client *Client) *Client {
		client.hmacApiKeyId = hmacApiKeyId
		return client
	}
}

func WithHttpClient(httpClient *http.Client) Option {
	return func(client *Client) *Client {
		client.httpClient = httpClient
		return client
	}
}

func WithLogger(logger *log.Logger) Option {
	return func(client *Client) *Client {
		client.log = logger
		return client
	}
}

func WithoutLogs() Option {
	return func(client *Client) *Client {
		client.log.Out = ioutil.Discard
		return client
	}
}

func WithLogLevel(level log.Level) Option {
	return func(client *Client) *Client {
		client.log.SetLevel(level)
		return client
	}
}

func WithLogFormatter(formatter log.Formatter) Option {
	return func(client *Client) *Client {
		client.log.SetFormatter(formatter)
		return client
	}
}

// Returns a Client give it the base url e.g. https://api.hanko.io and your API Secret
func NewClient(baseUrl string, secret string, opts ...Option) (client *Client) {
	client = &Client{
		BaseUrl:    baseUrl,
		secret:     secret,
		ApiVersion: "v1",
		log:        log.New(),
		httpClient: &http.Client{},
	}

	client.log.SetFormatter(&log.JSONFormatter{})

	for _, opt := range opts {
		client = opt(client)
	}

	if client.hmacApiKeyId == "" {
		client.log.Warn("hmac authentication is disabled. " +
			"please provide a valid api key id using the WithHmac() option.")
	}

	client.log.Debugf("Hanko client created (base url: %s)", client.BaseUrl)

	return client
}

func (c *Client) NewHttpRequest(method string, requestUrl string, requestBody interface{}) (httpRequest *http.Request, err error) {
	parsedRequestUrl, err := url.Parse(requestUrl)
	if err != nil {
		return nil, errors.Errorf("failed to parse url: '%s'", requestUrl)
	}

	encodedRequestBody := new(bytes.Buffer)
	if requestBody != nil {
		err = json.NewEncoder(encodedRequestBody).Encode(requestBody)
		if err != nil {
			return nil, errors.Wrap(err, "failed to encode http request body")
		}
	}

	httpRequest, err = http.NewRequest(method, parsedRequestUrl.String(), encodedRequestBody)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a new http request")
	}

	authorizationHeader := c.getAuthorizationHeader(method, parsedRequestUrl, encodedRequestBody)
	httpRequest.Header.Add("Authorization", authorizationHeader)
	httpRequest.Header.Add("Content-Type", "application/json")

	return httpRequest, nil
}

func (c *Client) HttpClientDo(httpRequest *http.Request) (httpResponse *http.Response, err error) {
	httpResponse, err = c.httpClient.Do(httpRequest)
	if err != nil {
		return nil, errors.Wrap(err, "could not do request")
	}

	if !strings.HasPrefix(strconv.Itoa(httpResponse.StatusCode), "2") {
		return httpResponse, errors.Errorf("got status code: %s", httpResponse.Status)
	}

	return httpResponse, nil
}

func (c *Client) getAuthorizationHeader(method string, url *url.URL, body *bytes.Buffer) (authHeader string) {
	if c.hmacApiKeyId != "" {
		hmac := &HmacMessageData{
			apiSecret:     c.secret,
			apiKeyId:      c.hmacApiKeyId,
			requestMethod: method,
			requestPath:   url.Path,
			requestBody:   body.String(),
		}
		authHeader = fmt.Sprintf("hanko %s", CalculateHmac(hmac))
	} else {
		authHeader = fmt.Sprintf("secret %s", c.secret)
	}
	return authHeader
}

func (c *Client) decodeHttpResponse(httpResponse *http.Response, responseType interface{}, ctxLogger *log.Entry) (err error) {
	responseTypeName := reflect.TypeOf(responseType).String()
	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read http response body")
	}
	ctxLogger.WithField("raw_response", string(body)).Debug("http response body read")
	err = json.Unmarshal(body, responseType)
	if err != nil {
		ctxLogger.WithField("response_type", responseTypeName).Error("failed to decode http response")
		return errors.Wrap(err, "failed to decode http response")
	}
	ctxLogger.WithFields(log.Fields{
		"decoded_response": fmt.Sprintf("%+v", responseType),
		"response_type":    responseTypeName,
	}).Debug("http response body decoded")
	return nil
}

func (c *Client) Request(action string, method string, requestUrl string, requestBody interface{}, responseType interface{}) *ApiError {
	ctxLogger := c.log.WithFields(log.Fields{
		"action": action,
		"method": method,
		"url":    requestUrl,
	})

	ctxLogger.Debug("new http request")

	if requestBody != nil {
		ctxLogger.WithFields(log.Fields{
			"request_type": reflect.TypeOf(requestBody).String(),
			"request":      fmt.Sprintf("%+v", requestBody),
		}).Debug("got request body")
	}

	httpRequest, err := c.NewHttpRequest(method, requestUrl, requestBody)
	if err != nil {
		ctxLogger.WithError(err).Error("failed to create http request")
		return WrapError(err)
	}

	httpResponse, err := c.HttpClientDo(httpRequest)
	if err != nil {
		if httpResponse != nil {
			apiErr := &ApiError{}
			decErr := c.decodeHttpResponse(httpResponse, apiErr, ctxLogger)
			if decErr == nil {
				ctxLogger.WithError(err).WithFields(log.Fields{
					"debug_message": apiErr.DebugMessage,
					"details":       apiErr.Details,
				}).Error(apiErr.Message)
				return apiErr
			}
		}
		ctxLogger.WithError(err).Error("hanko api call failed")
		return WrapError(err)
	}

	if responseType != nil {
		err = c.decodeHttpResponse(httpResponse, responseType, ctxLogger)
		if err != nil {
			ctxLogger.WithError(err).Error("failed to decode the hanko api response")
			return WrapError(err)
		}
	}

	ctxLogger.Info("hanko api call succeeded")
	return nil
}
