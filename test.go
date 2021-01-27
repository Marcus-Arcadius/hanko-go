package hankoApiClient

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
)

var testPort = ":9496"
var testBaseUrl = "http://" + testPort
var testApiSecret = "test"
var testHmacApiKeyId = "test"

func runTestApi(request interface{}, response interface{}) *httptest.Server {
	ts := httptest.NewUnstartedServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			if request != nil {
				err := json.NewDecoder(r.Body).Decode(&request)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
			}

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
