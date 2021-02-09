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

func WithHttpClient(httpClient *http.Client) Option {
	return func(client *HankoApiClient) *HankoApiClient {
		client.httpClient = httpClient
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
		client.log.Warn("hmac authentication is disabled. " +
			"please provide a valid api key id using the WithHmac() option.")
	}

	client.log.Debugf("Hanko client created (base url: %s)", client.baseUrl)

	return client
}

func (c *HankoApiClient) NewHttpRequest(method string, requestUrl string, requestBody interface{}) (httpRequest *http.Request, err error) {
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

func (c *HankoApiClient) Do(httpRequest *http.Request) (httpResponse *http.Response, err error) {
	httpResponse, err = c.httpClient.Do(httpRequest)
	if err != nil {
		return nil, errors.Wrap(err, "could not do request")
	}

	if !strings.HasPrefix(strconv.Itoa(httpResponse.StatusCode), "2") {
		return httpResponse, errors.Errorf("got status code: %s", httpResponse.Status)
	}

	return httpResponse, nil
}

func (c *HankoApiClient) getAuthorizationHeader(method string, url *url.URL, body *bytes.Buffer) (authHeader string) {
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

func (c *HankoApiClient) decodeHttpResponse(httpResponse *http.Response, responseType interface{}, ctxLogger *log.Entry) (err error) {
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

func (c *HankoApiClient) run(action string, method string, requestUrl string, requestBody interface{}, responseType interface{}) error {
	ctxLogger := c.log.WithFields(log.Fields{
		"action": action,
		"method": method,
		"url":    requestUrl,
	})

	ctxLogger.WithFields(log.Fields{
		"request_type": reflect.TypeOf(requestBody).String(),
		"request":      fmt.Sprintf("%+v", requestBody),
	}).Debug("new http request")

	httpRequest, err := c.NewHttpRequest(method, requestUrl, requestBody)
	if err != nil {
		ctxLogger.WithError(err).Error("failed to create http request")
		return err
	}

	httpResponse, err := c.Do(httpRequest)
	if err != nil {
		if httpResponse != nil {
			apiErr := &Error{}
			if decErr := c.decodeHttpResponse(httpResponse, apiErr, ctxLogger); decErr == nil {
				ctxLogger.WithError(err).Error(apiErr.Message)
				return errors.New(apiErr.Message)
			}
		}
		ctxLogger.WithError(err).Error("hanko api call failed")
		return err
	}

	if responseType != nil {
		err = c.decodeHttpResponse(httpResponse, responseType, ctxLogger)
		if err != nil {
			ctxLogger.WithError(err).Error("failed to decode the hanko api response")
			return err
		}
	}

	ctxLogger.Info("hanko api call succeeded")
	return nil
}
