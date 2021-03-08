package webauthn

import (
	hankoClient "github.com/teamhanko/hanko-sdk-golang/client"
	"github.com/teamhanko/webauthn/protocol"
	"time"
)

//  Registration - Initialization

type RegistrationInitializationRequest struct {
	User    hankoClient.User                         `json:"user"`
	Options RegistrationInitializationRequestOptions `json:"options"`
}

type RegistrationInitializationRequestOptions struct {
	AuthenticatorSelection protocol.AuthenticatorSelection `json:"authenticatorSelection"`
	ConveyancePreference   protocol.ConveyancePreference   `json:"attestation"`
}

type RegistrationInitializationResponse struct {
	protocol.PublicKeyCredentialCreationOptions
}

//  Registration - Finalization

type RegistrationFinalizationRequest struct {
	protocol.CredentialCreationResponse
}

type RegistrationFinalizationResponse struct {
	Status     hankoClient.OperationStatus `json:"status"`
	User       hankoClient.User            `json:"user"`
	Credential *Credential                 `json:"credential"`
	Error      *hankoClient.Error          `json:"error,omitempty"`
}

type Authenticator struct {
	Aaguid     string `json:"aaguid,omitempty"`
	Name       string `json:"name,omitempty"`
	Attachment string `json:"attachment,omitempty"`
}

//  Authentication - Initialization

type AuthenticationInitializationRequestOptions struct {
	UserVerification        protocol.UserVerificationRequirement `json:"userVerification"`
	AuthenticatorAttachment protocol.AuthenticatorAttachment     `json:"authenticatorAttachment"`
}

type AuthenticationInitializationRequest struct {
	User    hankoClient.User                           `json:"user"`
	Options AuthenticationInitializationRequestOptions `json:"options"`
}

type AuthenticationInitializationResponse struct {
	protocol.PublicKeyCredentialRequestOptions
}

//  Authentication - Finalization
type AuthenticationFinalizationRequest struct {
	protocol.CredentialAssertionResponse
}

type AuthenticationFinalizationResponse struct {
	Status hankoClient.OperationStatus `json:"status"`
	User   hankoClient.User            `json:"user"`
	Error  *hankoClient.Error          `json:"error,omitempty"`
}

//  Transactions - Initialization

type TransactionInitializationRequest struct {
	AuthenticationInitializationRequest
	Transaction string `json:"transaction"`
}

type TransactionInitializationResponse struct {
	AuthenticationInitializationResponse
}

//  Transactions - Finalization

type TransactionFinalizationRequest struct {
	AuthenticationFinalizationRequest
}

type TransactionFinalizationResponse struct {
	AuthenticationFinalizationResponse
}

//  Credentials

type CredentialQuery struct {
	UserId   string `json:"user_id" url:"user_id"`
	PageSize uint   `json:"page_size" url:"page_size"`
	Page     uint   `json:"page" url:"page"`
}

type Credential struct {
	Id               string         `json:"id"`
	CreatedAt        time.Time      `json:"createdAt"`
	LastUsed         time.Time      `json:"lastUsed"`
	Name             string         `json:"name"`
	UserVerification bool           `json:"userVerification"`
	IsResidentKey    bool           `json:"isResidentKey"`
	Authenticator    *Authenticator `json:"authenticator,omitempty"`
}

type CredentialUpdateRequest struct {
	Name string `json:"name"`
}
