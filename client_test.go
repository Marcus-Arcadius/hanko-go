package hankoApiClient

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestHankoApiClient_NewHttpRequestWithHmac(t *testing.T) {
	client := NewHankoApiClient(testBaseUrl, testApiSecret, WithHmac(testHmacApiKeyId), WithoutLogs())
	requestBody := &struct{ foo string }{"bar"}
	request, err := client.NewHttpRequest(http.MethodPost, "/test", &requestBody)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	authorizationHeader := request.Header.Get("Authorization")
	if !strings.HasPrefix(authorizationHeader, "hanko eyJobWFjQXBpS2V5SWQ") {
		t.Errorf("wrong authorization header, got: %s", authorizationHeader)
		t.Fail()
	}
}

func TestHankoApiClient_NewHttpRequestWithoutHmac(t *testing.T) {
	client := NewHankoApiClient(testBaseUrl, testApiSecret, WithoutLogs())
	requestBody := &struct{ foo string }{"bar"}
	request, err := client.NewHttpRequest(http.MethodPost, "/test", requestBody)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	authorizationHeader := request.Header.Get("Authorization")
	if authorizationHeader != fmt.Sprintf("secret %s", testApiSecret) {
		t.Errorf("wrong authorization header, got: %s", authorizationHeader)
		t.Fail()
	}
}

func TestHankoApiClient_Do(t *testing.T) {
	client := NewHankoApiClient(testBaseUrl, testApiSecret, WithoutLogs())
	httpRequest, err := http.NewRequest(http.MethodPost, testBaseUrl, nil)

	ts := runTestApi(nil, nil, http.StatusOK)
	ts.Start()
	_, err = client.Do(httpRequest)
	if err != nil {
		t.Error("no error expected")
		t.Fail()
	}
	ts.Close()

	ts = runTestApi(nil, nil, http.StatusCreated)
	ts.Start()
	_, err = client.Do(httpRequest)
	if err != nil {
		t.Error("no error expected")
		t.Fail()
	}
	ts.Close()

	ts = runTestApi(nil, nil, http.StatusBadRequest)
	ts.Start()
	_, err = client.Do(httpRequest)
	if err == nil {
		t.Error("error expected")
		t.Fail()
	}
	ts.Close()
}
