package hankoApiClient

import "time"

type UserEntity struct {
	// TODO
}

type WebAuthnAuthenticator struct {
	AaGuid     string
	Name       string
	Attachment string
}

type AuthenticatorAttestationResponse struct {
	ClientDataJSON    string
	AttestationObject string
}

type AuthenticatorAssertionResponse struct {
	ClientDataJSON    string
	AttestationObject string
	AuthenticatorData string
	UserHandle        string
}

// WebAuthn Registration
// - Initialization

type WebAuthnRegistrationInitializationRequest struct {
	UserEntity UserEntity
	Options    WebAuthnRegistrationInitializationRequestOptions
}

type WebAuthnAuthenticatorSelectionCriteria struct {
	UserVerification        string `json:"userVerification"`
	AuthenticatorAttachment string `json:"authenticatorAttachment"`
}

type WebAuthnRegistrationInitializationRequestOptions struct {
	AuthenticatorSelection WebAuthnAuthenticatorSelectionCriteria
	Attestation            string
	// TODO: Extensions
}

type WebAuthnRegistrationInitializationResponse struct {
	Challenge              string
	RpId                   string
	UserEntity             UserEntity
	PubKeyCredParams       string
	AuthenticatorSelection string
	Timeout                string
	ExcludeCredentials     string
	// TODO: Extensions
	Attestation string
}

// - Finalization

type WebAuthnRegistrationFinalizationRequest struct {
	Id       string
	RawId    string
	Type     string
	Response AuthenticatorAttestationResponse
}

type WebAuthnRegistrationFinalizationResponse struct {
	Status     string
	UserEntity UserEntity
	Credential WebAuthnCredential
	ErrorMsg   string
}

// WebAuthn Authentication
// - Initialization

type WebAuthnAuthenticationInitializationRequestOptions struct {
	UserVerification        string
	AuthenticatorAttachment string
	// TODO: Extension
}

type WebAuthnAuthenticationInitializationRequest struct {
	UserEntity UserEntity
	Options    WebAuthnAuthenticationInitializationRequestOptions
}

type WebAuthnAuthenticationInitializationResponse struct {
	Challenge string
	RpId      string
	Timeout   string
	// TODO: AllowCredentials
	// TODO: Extensions
}

// - Finalization
type WebAuthnAuthenticationFinalizationRequest struct {
	Id       string
	RawId    string
	Type     string
	Response AuthenticatorAssertionResponse
}

type WebAuthnAuthenticationFinalizationResponse struct {
	UserEntity UserEntity
	Status     string
	ErrorMsg   string
}

// Transactions
// - Initialization

type WebAuthnTransactionInitializationRequest struct {
	WebAuthnAuthenticationInitializationRequest
	Transaction string
}

type WebAuthnTransactionInitializationResponse struct {
	WebAuthnAuthenticationInitializationResponse
}

// - Finalization
type WebAuthnTransactionFinalizationRequest struct {
	WebAuthnAuthenticationFinalizationRequest
}

type WebAuthnTransactionFinalizationResponse struct {
	WebAuthnAuthenticationFinalizationResponse
}

type WebAuthnCredentialQuery struct {
	UserId   interface{} `url:"user_id"`
	PageSize interface{} `url:"page_size"`
	Page     interface{} `url:"page"`
}

type WebAuthnCredential struct {
	Id            string
	CreatedAt     time.Time
	LastUsed      time.Time
	Name          string
	Authenticator WebAuthnAuthenticator
}

type WebAuthnCredentialUpdateRequest struct {
	Name string
}

// Error
type ErrorResponse struct {
	Code    uint
	Message string
}
