package hankoApiClient

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

var port = ":9496"
var apiHost = "http://" + port
var apiSecret = "test"
var apiKeyId = "test"

func createTestServer(request interface{}, response interface{}) *httptest.Server {
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

	l, _ := net.Listen("tcp", port)
	ts.Listener = l

	return ts
}

func TestHankoApiClient_RegistrationInitialization(t *testing.T) {
	request := &WebAuthnRegistrationInitializationRequest{}
	response := &WebAuthnRegistrationInitializationResponse{}
	ts := createTestServer(request, response)
	ts.Start()
	defer ts.Close()

	client := NewHankoHmacClient(apiHost, apiSecret, apiKeyId)
	_, err := client.InitWebAuthnRegistration(request)
	if err != nil {
		t.Error(err)
	}
}

func TestHankoApiClient_RegistrationFinalization(t *testing.T) {
	request := &WebAuthnRegistrationFinalizationRequest{}
	response := &WebAuthnRegistrationFinalizationResponse{}
	ts := createTestServer(request, response)
	ts.Start()
	defer ts.Close()

	client := NewHankoHmacClient(apiHost, apiSecret, apiKeyId)
	_, err := client.FinalizeWebAuthnRegistration(request)
	if err != nil {
		t.Error(err)
	}
}

func TestHankoApiClient_AuthenticationInitialization(t *testing.T) {
	request := &WebAuthnAuthenticationInitializationRequest{}
	response := &WebAuthnAuthenticationInitializationResponse{}
	ts := createTestServer(request, response)
	ts.Start()
	defer ts.Close()

	client := NewHankoHmacClient(apiHost, apiSecret, apiKeyId)
	_, err := client.InitWebAuthnAuthentication(request)
	if err != nil {
		t.Error(err)
	}
}

func TestHankoApiClient_AuthenticationFinalization(t *testing.T) {
	request := &WebAuthnAuthenticationFinalizationRequest{}
	response := &WebAuthnAuthenticationFinalizationResponse{}
	ts := createTestServer(request, response)
	ts.Start()
	defer ts.Close()

	client := NewHankoHmacClient(apiHost, apiSecret, apiKeyId)
	_, err := client.FinalizeWebAuthnAuthentication(request)
	if err != nil {
		t.Error(err)
	}
}

func TestHankoApiClient_TransactionInitialization(t *testing.T) {
	request := &WebAuthnTransactionInitializationRequest{}
	response := &WebAuthnTransactionInitializationResponse{}
	ts := createTestServer(request, response)
	ts.Start()
	defer ts.Close()

	client := NewHankoHmacClient(apiHost, apiSecret, apiKeyId)
	_, err := client.InitWebAuthnTransaction(request)
	if err != nil {
		t.Error(err)
	}
}

func TestHankoApiClient_TransactionFinalization(t *testing.T) {
	request := &WebAuthnTransactionFinalizationRequest{}
	response := &WebAuthnTransactionFinalizationResponse{}
	ts := createTestServer(request, response)
	ts.Start()
	defer ts.Close()

	client := NewHankoHmacClient(apiHost, apiSecret, apiKeyId)
	_, err := client.FinalizeWebAuthnTransaction(request)
	if err != nil {
		t.Error(err)
	}
}

func TestHankoApiClient_ListWebAuthnCredentials(t *testing.T) {
	response := &[]WebAuthnCredential{}
	ts := createTestServer(nil, response)
	ts.Start()
	defer ts.Close()

	client := NewHankoHmacClient(apiHost, apiSecret, apiKeyId)
	_, err := client.ListWebAuthnCredentials(&WebAuthnCredentialQuery{UserId: "test", PageSize: nil, Page: 3})
	if err != nil {
		t.Error(err)
	}
}

func TestHankoApiClient_GetWebAuthnCredential(t *testing.T) {
	response := &WebAuthnCredential{}
	ts := createTestServer(nil, response)
	ts.Start()
	defer ts.Close()

	client := NewHankoHmacClient(apiHost, apiSecret, apiKeyId)
	_, err := client.GetWebAuthnCredential("test")
	if err != nil {
		t.Error(err)
	}
}

func TestHankoApiClient_DeleteWebAuthnCredential(t *testing.T) {
	ts := createTestServer(nil, nil)
	ts.Start()
	defer ts.Close()

	client := NewHankoHmacClient(apiHost, apiSecret, apiKeyId)
	err := client.DeleteWebAuthnCredential("test")
	if err != nil {
		t.Error(err)
	}
}

func TestHankoApiClient_UpdateWebAuthnCredential(t *testing.T) {
	request := &WebAuthnCredentialUpdateRequest{}
	response := &WebAuthnCredential{}
	ts := createTestServer(request, response)
	ts.Start()
	defer ts.Close()

	client := NewHankoHmacClient(apiHost, apiSecret, apiKeyId)
	_, err := client.UpdateWebAuthnCredential("test", request)
	if err != nil {
		t.Error(err)
	}
}
