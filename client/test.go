package client

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
)

var TestPort = ":9496"
var TestBaseUrl = "http://" + TestPort
var TestApiSecret = "test"
var TestHmacApiKeyId = "test"

func RunTestApi(requestType interface{}, response interface{}, responseStatus int) *httptest.Server {
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

	l, _ := net.Listen("tcp", TestPort)
	ts.Listener = l

	return ts
}
