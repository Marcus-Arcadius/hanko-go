package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"testing"
)

func getTestClient() *Client {
	client := NewClient(TestBaseUrl, TestApiSecret)
	client.SetLogWriter(ioutil.Discard)
	return client
}

func TestHankoApiClient_NewHttpRequestWithHmac(t *testing.T) {
	client := getTestClient()
	client.SetHmac(TestHmacApiKeyId)
	requestBody := &struct{ foo string }{"bar"}
	request, err := client.NewHttpRequest(http.MethodPost, "/test", &requestBody)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	// the authorization header must contain the word "hanko" followed by space followed by a base64 encoded string.
	// the string must be at least 240 chars long to pass the test
	authorizationHeader := request.Header.Get("Authorization")
	matched, _ := regexp.MatchString(`^hanko [A-Za-z0-9+/]{240,}`, authorizationHeader)
	if !matched {
		t.Errorf("wrong authorization header, got: %s", authorizationHeader)
		t.Fail()
	}
}

func TestHankoApiClient_NewHttpRequestWithoutHmac(t *testing.T) {
	client := getTestClient()
	requestBody := &struct{ foo string }{"bar"}
	request, err := client.NewHttpRequest(http.MethodPost, "/test", requestBody)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	authorizationHeader := request.Header.Get("Authorization")
	if authorizationHeader != fmt.Sprintf("secret %s", TestApiSecret) {
		t.Errorf("wrong authorization header, got: %s", authorizationHeader)
		t.Fail()
	}
}

func TestHankoApiClient_Do(t *testing.T) {
	client := getTestClient()
	httpRequest, err := http.NewRequest(http.MethodPost, TestBaseUrl, nil)

	ts := RunTestApi(nil, nil, http.StatusOK)
	ts.Start()
	_, err = client.HttpClientDo(httpRequest)
	if err != nil {
		t.Error("no error expected")
		t.Fail()
	}
	ts.Close()

	ts = RunTestApi(nil, nil, http.StatusCreated)
	ts.Start()
	_, err = client.HttpClientDo(httpRequest)
	if err != nil {
		t.Error("no error expected")
		t.Fail()
	}
	ts.Close()

	ts = RunTestApi(nil, nil, http.StatusBadRequest)
	ts.Start()
	_, err = client.HttpClientDo(httpRequest)
	if err == nil {
		t.Error("error expected")
		t.Fail()
	}
	ts.Close()
}
