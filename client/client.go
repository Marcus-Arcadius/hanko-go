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
)

type Client struct {
	baseUrl      string
	apiVersion   string
	secret       string
	hmacApiKeyId string
	httpClient   *http.Client
	log          *log.Logger
}

func NewClient(baseUrl string, secret string) *Client {
	client := &Client{
		baseUrl:    baseUrl,
		secret:     secret,
		apiVersion: "v1",
		log:        log.New(),
		httpClient: &http.Client{},
	}

	client.log.SetFormatter(&log.JSONFormatter{})

	if client.hmacApiKeyId == "" {
		client.log.Warn("hmac authentication is disabled. " +
			"please provide a valid api key id using the WithHmac() option.")
	}

	client.log.Debugf("Hanko client created (base url: %s)", client.baseUrl)

	return client
}

func (c *Client) GetUrl() string {
	return fmt.Sprintf("%s/%s", c.baseUrl, c.apiVersion)
}

func (c *Client) SetHmac(hmacApiKeyId string) {
	c.hmacApiKeyId = hmacApiKeyId
}

func (c *Client) SetHttpClient(httpClient *http.Client) {
	c.httpClient = httpClient
}

func (c *Client) SetLogger(logger *log.Logger) {
	c.log = logger
}

func (c *Client) SetLogWriter(out io.Writer) {
	c.log.Out = out
}

func (c *Client) SetLogLevel(level log.Level) {
	c.log.SetLevel(level)
}

func (c *Client) SetLogFormatter(formatter log.Formatter) {
	c.log.SetFormatter(formatter)
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
