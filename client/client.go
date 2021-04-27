// Package client provides a base client for communication with the Hanko Authentication API. Please take a look at our
// protocol-specific client SDKs:
//	- WebAuthn: https://github.com/teamhanko/hanko-sdk-golang/webauthn
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Client is used to communicate with the Hanko Authentication API. It handles authentication, provides logging and
// functionality for decoding the responses.
type Client struct {
	baseUrl      string       // the url of the hanko server, e.g. https://api.hanko.io
	apiVersion   string       // version string to be appended to baseUrl
	secret       string       // required to access the hanko api
	hmacApiKeyId string       // contains the api key id when HMAC is used
	httpClient   *http.Client // for http communication with the hanko server
	log          *log.Logger  // logrus logger
}

// NewClient returns a new basic hanko Client. pass in the base url (e.g. https://api.hanko.io) and your api secret.
func NewClient(baseUrl string, secret string) *Client {
	client := &Client{
		baseUrl:    baseUrl,
		secret:     secret,
		apiVersion: "v1",
		log:        log.New(),
		httpClient: &http.Client{Timeout: time.Second * 10},
	}
	client.log.SetFormatter(&log.JSONFormatter{})
	return client
}

// GetUrl returns a concatenation of the API baseUrl and apiVersion.
func (c *Client) GetUrl() string {
	return fmt.Sprintf("%s/%s", c.baseUrl, c.apiVersion)
}

// SetHmac sets the given hmacApiKeyId to be used while generating the authorization header.
func (c *Client) SetHmac(hmacApiKeyId string) {
	c.hmacApiKeyId = hmacApiKeyId
}

// SetHttpClient sets a different http.Client. To example if you want to use an http proxy you can pass your own
// configured http.Client.
func (c *Client) SetHttpClient(httpClient *http.Client) {
	c.httpClient = httpClient
}

// SetLogger allows you to set a custom logrus.Logger.
func (c *Client) SetLogger(logger *log.Logger) {
	c.log = logger
}

// SetLogWriter allows you to change the log output to the given io.Writer.
func (c *Client) SetLogWriter(out io.Writer) {
	c.log.Out = out
}

// SetLogLevel allows you to change the log level. The default is logrus.InfoLevel.
func (c *Client) SetLogLevel(level log.Level) {
	c.log.SetLevel(level)
}

// SetLogFormatter allows you to set a custom log formatter. The default is logrus.JSONFormatter.
func (c *Client) SetLogFormatter(formatter log.Formatter) {
	c.log.SetFormatter(formatter)
}

// NewHttpRequest encodes the given requestBody, creates a new HTTP request using the http.NewRequest method and
// sets the authorization and content-type headers.
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

// HttpClientDo calls the API with the specified http.Request and returns an error if the status code was not 2xx.
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

// getAuthorizationHeader calculates an HMAC with the specified values for method, url and body. If the client has been
// created without the client.WithHmac option, the api secret will be used for authentication. Returns the
// an HTTP authorization header as a string.
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

// decodeHttpResponse decodes the httpResponse into the given responseType.
func (c *Client) decodeHttpResponse(httpResponse *http.Response, responseType interface{}, ctxLogger *log.Entry) (err error) {
	responseTypeName := reflect.TypeOf(responseType).String()
	defer httpResponse.Body.Close()
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

// Request is used to make a request to the API. The action parameter should contain a string that indicates which
// action is currently performed. Parameters method, requestUrl and requestBody are used to construct the request.
// The body of the API response will be decoded into the given responseType. On error, returns an ApiError.
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
