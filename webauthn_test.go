package hankoApiClient

import (
	"testing"
)

func TestHankoApiClient_RegistrationInitialization(t *testing.T) {
	requestBody := &WebAuthnRegistrationInitializationRequest{}
	responseType := &WebAuthnRegistrationInitializationResponse{}
	ts := runTestApi(requestBody, responseType)
	ts.Start()
	defer ts.Close()
	client := NewHankoApiClient(testBaseUrl, testApiSecret, WithoutLogs())
	_, err := client.InitializeWebAuthnRegistration(requestBody)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestHankoApiClient_RegistrationFinalization(t *testing.T) {
	requestBody := &WebAuthnRegistrationFinalizationRequest{}
	responseType := &WebAuthnRegistrationFinalizationResponse{}
	ts := runTestApi(requestBody, responseType)
	ts.Start()
	defer ts.Close()
	client := NewHankoApiClient(testBaseUrl, testApiSecret, WithoutLogs())
	_, err := client.FinalizeWebAuthnRegistration(requestBody)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestHankoApiClient_AuthenticationInitialization(t *testing.T) {
	requestBody := &WebAuthnAuthenticationInitializationRequest{}
	responseType := &WebAuthnAuthenticationInitializationResponse{}
	ts := runTestApi(requestBody, responseType)
	ts.Start()
	defer ts.Close()
	client := NewHankoApiClient(testBaseUrl, testApiSecret, WithoutLogs())
	_, err := client.InitializeWebAuthnAuthentication(requestBody)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestHankoApiClient_AuthenticationFinalization(t *testing.T) {
	requestBody := &WebAuthnAuthenticationFinalizationRequest{}
	responseType := &WebAuthnAuthenticationFinalizationResponse{}
	ts := runTestApi(requestBody, responseType)
	ts.Start()
	defer ts.Close()
	client := NewHankoApiClient(testBaseUrl, testApiSecret, WithoutLogs())
	_, err := client.FinalizeWebAuthnAuthentication(requestBody)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestHankoApiClient_TransactionInitialization(t *testing.T) {
	requestBody := &WebAuthnTransactionInitializationRequest{}
	responseType := &WebAuthnTransactionInitializationResponse{}
	ts := runTestApi(requestBody, responseType)
	ts.Start()
	defer ts.Close()
	client := NewHankoApiClient(testBaseUrl, testApiSecret, WithoutLogs())
	_, err := client.InitializeWebAuthnTransaction(requestBody)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestHankoApiClient_TransactionFinalization(t *testing.T) {
	requestBody := &WebAuthnTransactionFinalizationRequest{}
	responseType := &WebAuthnTransactionFinalizationResponse{}
	ts := runTestApi(requestBody, responseType)
	ts.Start()
	defer ts.Close()
	client := NewHankoApiClient(testBaseUrl, testApiSecret, WithoutLogs())
	_, err := client.FinalizeWebAuthnTransaction(requestBody)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestHankoApiClient_ListWebAuthnCredentials(t *testing.T) {
	responseType := &[]WebAuthnCredential{}
	ts := runTestApi(nil, responseType)
	ts.Start()
	defer ts.Close()
	client := NewHankoApiClient(testBaseUrl, testApiSecret, WithoutLogs())
	_, err := client.ListWebAuthnCredentials(&WebAuthnCredentialQuery{UserId: "test", Page: 3})
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestHankoApiClient_GetWebAuthnCredential(t *testing.T) {
	responseType := &WebAuthnCredential{}
	ts := runTestApi(nil, responseType)
	ts.Start()
	defer ts.Close()
	client := NewHankoApiClient(testBaseUrl, testApiSecret, WithoutLogs())
	_, err := client.GetWebAuthnCredential("test")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestHankoApiClient_DeleteWebAuthnCredential(t *testing.T) {
	ts := runTestApi(nil, nil)
	ts.Start()
	defer ts.Close()
	client := NewHankoApiClient(testBaseUrl, testApiSecret, WithoutLogs())
	err := client.DeleteWebAuthnCredential("test")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestHankoApiClient_UpdateWebAuthnCredential(t *testing.T) {
	requestBody := &WebAuthnCredentialUpdateRequest{}
	responseType := &WebAuthnCredential{}
	ts := runTestApi(requestBody, responseType)
	ts.Start()
	defer ts.Close()
	client := NewHankoApiClient(testBaseUrl, testApiSecret, WithoutLogs())
	_, err := client.UpdateWebAuthnCredential("test", requestBody)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}
