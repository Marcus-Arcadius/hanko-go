package hankoApiClient

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestHankoApiClient_CreateHttpRequestWithHmac(t *testing.T) {
	client := NewHankoApiClient(testBaseUrl, testApiSecret, WithHmac(testHmacApiKeyId), WithoutLogs())
	requestBody := &struct{ foo string }{"bar"}
	request, err := client.CreateHttpRequest(http.MethodPost, "/test", &requestBody)
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

func TestHankoApiClient_CreateHttpRequestWithoutHmac(t *testing.T) {
	client := NewHankoApiClient(testBaseUrl, testApiSecret, WithoutLogs())
	requestBody := &struct{ foo string }{"bar"}
	request, err := client.CreateHttpRequest(http.MethodPost, "/test", requestBody)
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
