package hankoApiClient

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

// HankoApiClient Provides Methods for interacting with the Hanko API
type HankoApiClient struct {
	baseUrl      string
	apiVersion   string
	secret       string
	hmacApiKeyId string
	httpClient   *http.Client
	log          *log.Logger
}

type Option func(client *HankoApiClient) *HankoApiClient

func WithHmac(hmacApiKeyId string) Option {
	return func(client *HankoApiClient) *HankoApiClient {
		client.hmacApiKeyId = hmacApiKeyId
		return client
	}
}

func WithLogger(logger *log.Logger) Option {
	return func(client *HankoApiClient) *HankoApiClient {
		client.log = logger
		return client
	}
}

func WithoutLogs() Option {
	return func(client *HankoApiClient) *HankoApiClient {
		client.log.Out = ioutil.Discard
		return client
	}
}

func WithLogLevel(level log.Level) Option {
	return func(client *HankoApiClient) *HankoApiClient {
		client.log.SetLevel(level)
		return client
	}
}

func WithLogFormatter(formatter log.Formatter) Option {
	return func(client *HankoApiClient) *HankoApiClient {
		client.log.SetFormatter(formatter)
		return client
	}
}

// Returns a HankoApiClient give it the base url e.g. https://api.hanko.io and your API Secret
func NewHankoApiClient(baseUrl string, secret string, opts ...Option) (client *HankoApiClient) {
	client = &HankoApiClient{
		baseUrl:    baseUrl,
		secret:     secret,
		apiVersion: "v1",
		log:        log.New(),
		httpClient: &http.Client{},
	}

	client.log.SetFormatter(&log.JSONFormatter{})

	for _, opt := range opts {
		client = opt(client)
	}

	if client.hmacApiKeyId == "" {
		client.log.Warn("HMAC authentication is disabled. " +
			"Please provide a valid API Key ID using the WithHMAC() option.")
	}

	client.log.Debugf("Hanko client created (base url: %s)", client.baseUrl)

	return client
}

func (c *HankoApiClient) NewHttpRequest(method string, requestUrl string, requestBody interface{}) (httpRequest *http.Request, err error) {
	parsedRequestUrl, err := url.Parse(requestUrl)
	if err != nil {
		return nil, errors.Errorf("Failed to parse URL: '%s'", requestUrl)
	}

	buf := new(bytes.Buffer)

	if requestBody != nil {
		err = json.NewEncoder(buf).Encode(requestBody)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to encode the requestBody")
		}
	}

	httpRequest, err = http.NewRequest(method, parsedRequestUrl.String(), buf)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create a new HTTP request")
	}

	httpRequest.Header.Add("Authorization", c.getAuthorizationHeader(method, parsedRequestUrl, buf))
	httpRequest.Header.Add("Content-Type", "application/json")

	return httpRequest, nil
}

func (c *HankoApiClient) getAuthorizationHeader(method string, url *url.URL, buf *bytes.Buffer) (authHeader string) {
	if c.hmacApiKeyId != "" {
		hmac := &HmacMessageData{
			apiSecret:     c.secret,
			apiKeyId:      c.hmacApiKeyId,
			requestMethod: method,
			requestPath:   url.Path,
			requestBody:   buf.String(),
		}
		authHeader = fmt.Sprintf("hanko %s", CalculateHmac(hmac))
	} else {
		authHeader = fmt.Sprintf("secret %s", c.secret)
	}
	return authHeader
}

func (c *HankoApiClient) decodeHttpResponse(httpResponse *http.Response, responseType interface{}, ctxLogger *log.Entry) (err error) {
	responseTypeName := reflect.TypeOf(responseType).String()
	body, err := ioutil.ReadAll(httpResponse.Body)
	ctxLogger.WithField("raw_response", string(body)).Debug("Response body read")
	err = json.Unmarshal(body, responseType)
	if err != nil {
		return errors.Wrapf(err, "Failed to decode Hanko API response (%s)", responseTypeName)
	}
	ctxLogger.Debugf("Response body (%s) decoded successfully", responseTypeName)
	return nil
}

func (c *HankoApiClient) run(action string, method string, requestUrl string, requestBody interface{}, responseType interface{}) error {
	ctxLogger := c.log.WithFields(log.Fields{
		"action": action,
		"method": method,
		"url":    requestUrl,
	})

	ctxLogger.Debug("Creating Hanko API response")

	httpRequest, err := c.NewHttpRequest(method, requestUrl, requestBody)
	if err != nil {
		ctxLogger.WithError(err).Error("HTTP request creation failed")
		return err
	}

	httpResponse, err := c.Do(httpRequest)
	if err != nil {
		ctxLogger.WithError(err).Error("Hanko API request failed")
		return err
	}

	if responseType != nil {
		err = c.decodeHttpResponse(httpResponse, responseType, ctxLogger)
		if err != nil {
			ctxLogger.WithError(err).Error("Failed to decode Hanko API response")
			return err
		}
	}

	ctxLogger.Info("Hanko API request succeeded")

	return nil
}

func (c *HankoApiClient) Do(httpRequest *http.Request) (httpResponse *http.Response, err error) {
	httpResponse, err = c.httpClient.Do(httpRequest)
	if err != nil {
		return nil, errors.Wrap(err, "Could not do response")
	}

	if !strings.HasPrefix(strconv.Itoa(httpResponse.StatusCode), "2") {
		return nil, errors.Errorf("Response status not ok (got code: %s)", httpResponse.Status)
	}

	return httpResponse, nil
}
