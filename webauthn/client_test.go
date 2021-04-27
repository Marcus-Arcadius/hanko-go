package webauthn

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testPort      = ":9496"
	testBaseUrl   = "http://" + testPort
	testApiSecret = "test"
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

func TestHankoApiClient_RegistrationInitialization(t *testing.T) {
	requestBody := &RegistrationInitializationRequest{}
	responseType := &RegistrationInitializationResponse{}
	ts := runTestApi(requestBody, responseType, http.StatusOK)
	ts.Start()
	defer ts.Close()
	client := NewClient(testBaseUrl, testApiSecret).WithoutLogs()
	_, err := client.InitializeRegistration(requestBody)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestHankoApiClient_RegistrationFinalization(t *testing.T) {
	requestBody := &RegistrationFinalizationRequest{}
	responseType := &RegistrationFinalizationResponse{}
	ts := runTestApi(requestBody, responseType, http.StatusOK)
	ts.Start()
	defer ts.Close()
	client := NewClient(testBaseUrl, testApiSecret).WithoutLogs()
	_, err := client.FinalizeRegistration(requestBody)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestHankoApiClient_AuthenticationInitialization(t *testing.T) {
	requestBody := &AuthenticationInitializationRequest{}
	responseType := &AuthenticationInitializationResponse{}
	ts := runTestApi(requestBody, responseType, http.StatusOK)
	ts.Start()
	defer ts.Close()
	client := NewClient(testBaseUrl, testApiSecret).WithoutLogs()
	_, err := client.InitializeAuthentication(requestBody)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestHankoApiClient_AuthenticationFinalization(t *testing.T) {
	requestBody := &AuthenticationFinalizationRequest{}
	responseType := &AuthenticationFinalizationResponse{}
	ts := runTestApi(requestBody, responseType, http.StatusOK)
	ts.Start()
	defer ts.Close()
	client := NewClient(testBaseUrl, testApiSecret).WithoutLogs()
	_, err := client.FinalizeAuthentication(requestBody)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestHankoApiClient_TransactionInitialization(t *testing.T) {
	requestBody := &TransactionInitializationRequest{}
	responseType := &TransactionInitializationResponse{}
	ts := runTestApi(requestBody, responseType, http.StatusOK)
	ts.Start()
	defer ts.Close()
	client := NewClient(testBaseUrl, testApiSecret).WithoutLogs()
	_, err := client.InitializeTransaction(requestBody)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestHankoApiClient_TransactionFinalization(t *testing.T) {
	requestBody := &TransactionFinalizationRequest{}
	responseType := &TransactionFinalizationResponse{}
	ts := runTestApi(requestBody, responseType, http.StatusOK)
	ts.Start()
	defer ts.Close()
	client := NewClient(testBaseUrl, testApiSecret).WithoutLogs()
	_, err := client.FinalizeTransaction(requestBody)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestHankoApiClient_ListCredentials(t *testing.T) {
	responseType := &[]Credential{}
	ts := runTestApi(nil, responseType, http.StatusOK)
	ts.Start()
	defer ts.Close()
	client := NewClient(testBaseUrl, testApiSecret).WithoutLogs()
	_, err := client.ListCredentials(&CredentialQuery{UserId: "test", Page: 3})
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestHankoApiClient_GetCredential(t *testing.T) {
	responseType := &Credential{}
	ts := runTestApi(nil, responseType, http.StatusOK)
	ts.Start()
	defer ts.Close()
	client := NewClient(testBaseUrl, testApiSecret).WithoutLogs()
	_, err := client.GetCredential("test")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestHankoApiClient_DeleteCredential(t *testing.T) {
	ts := runTestApi(nil, nil, http.StatusOK)
	ts.Start()
	defer ts.Close()
	client := NewClient(testBaseUrl, testApiSecret).WithoutLogs()
	err := client.DeleteCredential("test")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestHankoApiClient_UpdateCredential(t *testing.T) {
	requestBody := &CredentialUpdateRequest{}
	responseType := &Credential{}
	ts := runTestApi(requestBody, responseType, http.StatusOK)
	ts.Start()
	defer ts.Close()
	client := NewClient(testBaseUrl, testApiSecret).WithoutLogs()
	_, err := client.UpdateCredential("test", requestBody)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}
