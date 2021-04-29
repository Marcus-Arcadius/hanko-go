package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

const (
	testPort         = ":9496"
	testBaseUrl      = "http://" + testPort
	testApiSecret    = "test"
	testHmacApiKeyId = "test"
)

func runTestApi(requestType interface{}, response interface{}, responseStatus int) *httptest.Server {
	ts := httptest.NewUnstartedServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			if requestType != nil {
				dec := json.NewDecoder(r.Body)
				dec.DisallowUnknownFields()
				err := dec.Decode(&requestType)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
			}

			w.WriteHeader(responseStatus)

			if response != nil {
				_ = json.NewEncoder(w).Encode(response)
			}

			return
		}),
	)

	l, _ := net.Listen("tcp", testPort)
	ts.Listener = l

	return ts
}

func getTestClient() *Client {
	client := NewClient(testBaseUrl, testApiSecret)
	client.SetLogWriter(ioutil.Discard)
	return client
}

func TestHankoApiClient_NewHttpRequestWithHmac(t *testing.T) {
	client := getTestClient()
	client.SetHmac(testHmacApiKeyId)
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
	if authorizationHeader != fmt.Sprintf("secret %s", testApiSecret) {
		t.Errorf("wrong authorization header, got: %s", authorizationHeader)
		t.Fail()
	}
}

func TestHankoApiClient_Do(t *testing.T) {
	client := getTestClient()
	httpRequest, err := http.NewRequest(http.MethodPost, testBaseUrl, nil)

	ts := runTestApi(nil, nil, http.StatusOK)
	ts.Start()
	_, err = client.HttpClientDo(httpRequest)
	if err != nil {
		t.Error("no error expected")
		t.Fail()
	}
	ts.Close()

	ts = runTestApi(nil, nil, http.StatusCreated)
	ts.Start()
	_, err = client.HttpClientDo(httpRequest)
	if err != nil {
		t.Error("no error expected")
		t.Fail()
	}
	ts.Close()

	ts = runTestApi(nil, nil, http.StatusBadRequest)
	ts.Start()
	_, err = client.HttpClientDo(httpRequest)
	if err == nil {
		t.Error("error expected")
		t.Fail()
	}
	ts.Close()
}
