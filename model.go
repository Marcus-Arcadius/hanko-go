package hankoApiClient

import (
	"gitlab.com/hanko/webauthn/credential"
	"gitlab.com/hanko/webauthn/protocol"
	"time"
)

// Misc

type OperationStatus string

type Error struct {
	Message string `json:"message"`
	Code    uint   `json:"code"`
}

const (
	Ok     OperationStatus = "ok"
	Failed                 = "failed"
)

type User struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	IconUrl     string `json:"icon"`
}

// WebAuthn Registration - Initialization

type WebAuthnRegistrationInitializationRequest struct {
	User    User                                             `json:"user"`
	Options WebAuthnRegistrationInitializationRequestOptions `json:"options"`
}

type WebAuthnRegistrationInitializationRequestOptions struct {
	AuthenticatorSelection protocol.AuthenticatorSelection `json:"authenticatorSelection"`
	ConveyancePreference   protocol.ConveyancePreference   `json:"conveyancePreference"`
}

type WebAuthnRegistrationInitializationResponse struct {
	protocol.PublicKeyCredentialCreationOptions
}

// WebAuthn Registration - Finalization

type WebAuthnRegistrationFinalizationRequest struct {
	protocol.CredentialCreationResponse
}

type WebAuthnRegistrationFinalizationResponse struct {
	Status     OperationStatus       `json:"status"`
	User       User                  `json:"user"`
	Credential credential.Credential `json:"credential"`
	Error      *Error                `json:"error,omitempty"`
}

// WebAuthn Authentication - Initialization

type WebAuthnAuthenticationInitializationRequestOptions struct {
	UserVerification        protocol.UserVerificationRequirement `json:"userVerification"`
	AuthenticatorAttachment protocol.AuthenticatorAttachment     `json:"authenticatorAttachment"`
}

type WebAuthnAuthenticationInitializationRequest struct {
	User    User                                               `json:"user"`
	Options WebAuthnAuthenticationInitializationRequestOptions `json:"options"`
}

type WebAuthnAuthenticationInitializationResponse struct {
	protocol.PublicKeyCredentialRequestOptions
}

// WebAuthn Authentication - Finalization
type WebAuthnAuthenticationFinalizationRequest struct {
	protocol.CredentialAssertionResponse
}

type WebAuthnAuthenticationFinalizationResponse struct {
	Status OperationStatus `json:"status"`
	User   User            `json:"user"`
	Error  *Error          `json:"error,omitempty"`
}

// Transactions - Initialization

type WebAuthnTransactionInitializationRequest struct {
	WebAuthnAuthenticationInitializationRequest
	Transaction string `json:"transaction"`
}

type WebAuthnTransactionInitializationResponse struct {
	WebAuthnAuthenticationInitializationResponse
}

// Transactions - Finalization

type WebAuthnTransactionFinalizationRequest struct {
	WebAuthnAuthenticationFinalizationRequest
}

type WebAuthnTransactionFinalizationResponse struct {
	WebAuthnAuthenticationFinalizationResponse
}

// Credentials

type WebAuthnCredentialQuery struct {
	UserId   string `json:"user_id" url:"user_id"`
	PageSize uint   `json:"page_size" url:"page_size"`
	Page     uint   `json:"page" url:"page"`
}

type WebAuthnAuthenticator struct {
	AaGuid     string `json:"aa_guid"`
	Name       string `json:"name"`
	Attachment string `json:"attachment"`
}

type WebAuthnCredential struct {
	Id            string                `json:"id"`
	CreatedAt     time.Time             `json:"created_at"`
	LastUsed      time.Time             `json:"last_used"`
	Name          string                `json:"name"`
	Authenticator WebAuthnAuthenticator `json:"authenticator"`
}

type WebAuthnCredentialUpdateRequest struct {
	Name string `json:"name"`
}

