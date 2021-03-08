package webauthn

import (
	hankoClient "github.com/teamhanko/hanko-sdk-golang/client"
	"net/http"
	"testing"
)

func TestHankoApiClient_RegistrationInitialization(t *testing.T) {
	requestBody := &RegistrationInitializationRequest{}
	responseType := &RegistrationInitializationResponse{}
	ts := hankoClient.RunTestApi(requestBody, responseType, http.StatusOK)
	ts.Start()
	defer ts.Close()
	client := NewClient(hankoClient.TestBaseUrl, hankoClient.TestApiSecret, hankoClient.WithoutLogs())
	_, err := client.InitializeRegistration(requestBody)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestHankoApiClient_RegistrationFinalization(t *testing.T) {
	requestBody := &RegistrationFinalizationRequest{}
	responseType := &RegistrationFinalizationResponse{}
	ts := hankoClient.RunTestApi(requestBody, responseType, http.StatusOK)
	ts.Start()
	defer ts.Close()
	client := NewClient(hankoClient.TestBaseUrl, hankoClient.TestApiSecret, hankoClient.WithoutLogs())
	_, err := client.FinalizeRegistration(requestBody)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestHankoApiClient_AuthenticationInitialization(t *testing.T) {
	requestBody := &AuthenticationInitializationRequest{}
	responseType := &AuthenticationInitializationResponse{}
	ts := hankoClient.RunTestApi(requestBody, responseType, http.StatusOK)
	ts.Start()
	defer ts.Close()
	client := NewClient(hankoClient.TestBaseUrl, hankoClient.TestApiSecret, hankoClient.WithoutLogs())
	_, err := client.InitializeAuthentication(requestBody)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestHankoApiClient_AuthenticationFinalization(t *testing.T) {
	requestBody := &AuthenticationFinalizationRequest{}
	responseType := &AuthenticationFinalizationResponse{}
	ts := hankoClient.RunTestApi(requestBody, responseType, http.StatusOK)
	ts.Start()
	defer ts.Close()
	client := NewClient(hankoClient.TestBaseUrl, hankoClient.TestApiSecret, hankoClient.WithoutLogs())
	_, err := client.FinalizeAuthentication(requestBody)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestHankoApiClient_TransactionInitialization(t *testing.T) {
	requestBody := &TransactionInitializationRequest{}
	responseType := &TransactionInitializationResponse{}
	ts := hankoClient.RunTestApi(requestBody, responseType, http.StatusOK)
	ts.Start()
	defer ts.Close()
	client := NewClient(hankoClient.TestBaseUrl, hankoClient.TestApiSecret, hankoClient.WithoutLogs())
	_, err := client.InitializeTransaction(requestBody)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestHankoApiClient_TransactionFinalization(t *testing.T) {
	requestBody := &TransactionFinalizationRequest{}
	responseType := &TransactionFinalizationResponse{}
	ts := hankoClient.RunTestApi(requestBody, responseType, http.StatusOK)
	ts.Start()
	defer ts.Close()
	client := NewClient(hankoClient.TestBaseUrl, hankoClient.TestApiSecret, hankoClient.WithoutLogs())
	_, err := client.FinalizeTransaction(requestBody)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestHankoApiClient_ListCredentials(t *testing.T) {
	responseType := &[]Credential{}
	ts := hankoClient.RunTestApi(nil, responseType, http.StatusOK)
	ts.Start()
	defer ts.Close()
	client := NewClient(hankoClient.TestBaseUrl, hankoClient.TestApiSecret, hankoClient.WithoutLogs())
	_, err := client.ListCredentials(&CredentialQuery{UserId: "test", Page: 3})
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestHankoApiClient_GetCredential(t *testing.T) {
	responseType := &Credential{}
	ts := hankoClient.RunTestApi(nil, responseType, http.StatusOK)
	ts.Start()
	defer ts.Close()
	client := NewClient(hankoClient.TestBaseUrl, hankoClient.TestApiSecret, hankoClient.WithoutLogs())
	_, err := client.GetCredential("test")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestHankoApiClient_DeleteCredential(t *testing.T) {
	ts := hankoClient.RunTestApi(nil, nil, http.StatusOK)
	ts.Start()
	defer ts.Close()
	client := NewClient(hankoClient.TestBaseUrl, hankoClient.TestApiSecret, hankoClient.WithoutLogs())
	err := client.DeleteCredential("test")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestHankoApiClient_UpdateCredential(t *testing.T) {
	requestBody := &CredentialUpdateRequest{}
	responseType := &Credential{}
	ts := hankoClient.RunTestApi(requestBody, responseType, http.StatusOK)
	ts.Start()
	defer ts.Close()
	client := NewClient(hankoClient.TestBaseUrl, hankoClient.TestApiSecret, hankoClient.WithoutLogs())
	_, err := client.UpdateCredential("test", requestBody)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}
