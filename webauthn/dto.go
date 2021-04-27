package webauthn

import (
	hankoClient "github.com/teamhanko/hanko-sdk-golang/client"
	"github.com/teamhanko/webauthn/protocol"
	"time"
)

type RegistrationInitializationUser struct {
	hankoClient.User
}

func NewRegistrationInitializationUser(id string, name string) RegistrationInitializationUser {
	return RegistrationInitializationUser{hankoClient.User{
		ID:          id,
		Name:        name,
		DisplayName: name,
	}}
}

func (user RegistrationInitializationUser) WithDisplayName(displayName string) RegistrationInitializationUser {
	user.DisplayName = displayName
	return user
}

func (user RegistrationInitializationUser) WithIconUrl(iconUrl string) RegistrationInitializationUser {
	user.IconUrl = iconUrl
	return user
}

type AuthenticationInitializationUser struct {
	hankoClient.User
}

func NewAuthenticationInitializationUser(id string) AuthenticationInitializationUser {
	return AuthenticationInitializationUser{hankoClient.User{ID: id}}
}

func (user AuthenticationInitializationUser) WithName(name string) AuthenticationInitializationUser {
	user.Name = name
	return user
}

type RegistrationInitializationRequest struct {
	User    hankoClient.User                         `json:"user"`
	Options RegistrationInitializationRequestOptions `json:"options"`
}

func NewRegistrationInitializationRequest(user RegistrationInitializationUser) *RegistrationInitializationRequest {
	return &RegistrationInitializationRequest{User: user.User}
}

type RegistrationInitializationRequestOptions struct {
	AuthenticatorSelection *AuthenticatorSelection `json:"authenticatorSelection"`
	ConveyancePreference   ConveyancePreference    `json:"attestation"`
}

func (request *RegistrationInitializationRequest) WithConveyancePreference(conveyancePreference ConveyancePreference) *RegistrationInitializationRequest {
	request.Options.ConveyancePreference = conveyancePreference
	return request
}

type AuthenticatorSelection struct {
	AuthenticatorAttachment AuthenticatorAttachment     `json:"authenticatorAttachment,omitempty"`
	RequireResidentKey      bool                        `json:"requireResidentKey,omitempty"`
	UserVerification        UserVerificationRequirement `json:"userVerification,omitempty"`
}

func NewAuthenticatorSelection() *AuthenticatorSelection {
	return &AuthenticatorSelection{}
}

func (request *RegistrationInitializationRequest) WithAuthenticatorSelection(authenticatorSelection *AuthenticatorSelection) *RegistrationInitializationRequest {
	request.Options.AuthenticatorSelection = authenticatorSelection
	return request
}

func (authenticatorSelection *AuthenticatorSelection) WithAuthenticatorAttachment(authenticatorAttachment AuthenticatorAttachment) *AuthenticatorSelection {
	authenticatorSelection.AuthenticatorAttachment = authenticatorAttachment
	return authenticatorSelection
}

func (authenticatorSelection *AuthenticatorSelection) WithRequireResidentKey(requireResidentKey bool) *AuthenticatorSelection {
	authenticatorSelection.RequireResidentKey = requireResidentKey
	return authenticatorSelection
}

func (authenticatorSelection *AuthenticatorSelection) WithUserVerification(userVerificationRequirement UserVerificationRequirement) *AuthenticatorSelection {
	authenticatorSelection.UserVerification = userVerificationRequirement
	return authenticatorSelection
}

type RegistrationInitializationResponse struct {
	protocol.PublicKeyCredentialCreationOptions
}

type RegistrationFinalizationRequest struct {
	protocol.CredentialCreationResponse
}

type RegistrationFinalizationResponse struct {
	User       hankoClient.User `json:"user"`
	Credential *Credential      `json:"credential"`
}

type Authenticator struct {
	Aaguid     string `json:"aaguid,omitempty"`
	Name       string `json:"name,omitempty"`
	Attachment string `json:"attachment,omitempty"`
}

type AuthenticationInitializationRequest struct {
	User        hankoClient.User                           `json:"user"`
	Options     AuthenticationInitializationRequestOptions `json:"options"`
}

func NewAuthenticationInitializationRequest(user AuthenticationInitializationUser) (request *AuthenticationInitializationRequest) {
	request = &AuthenticationInitializationRequest{User: user.User}
	return request
}

type AuthenticationInitializationRequestOptions struct {
	UserVerification        UserVerificationRequirement `json:"userVerification"`
	AuthenticatorAttachment AuthenticatorAttachment     `json:"authenticatorAttachment"`
}

func (request *AuthenticationInitializationRequest) WithAuthenticatorAttachment(authenticatorAttachment AuthenticatorAttachment) *AuthenticationInitializationRequest {
	request.Options.AuthenticatorAttachment = authenticatorAttachment
	return request
}

func (request *AuthenticationInitializationRequest) WithUserVerification(userVerificationRequirement UserVerificationRequirement) *AuthenticationInitializationRequest {
	request.Options.UserVerification = userVerificationRequirement
	return request
}

type AuthenticationInitializationResponse struct {
	protocol.PublicKeyCredentialRequestOptions
}

type AuthenticationFinalizationRequest struct {
	protocol.CredentialAssertionResponse
}

type AuthenticationFinalizationResponse struct {
	User hankoClient.User `json:"user"`
}

type TransactionInitializationRequest struct {
	User        hankoClient.User                           `json:"user"`
	Options     AuthenticationInitializationRequestOptions `json:"options"`
	Transaction string `json:"transaction"`
}

func NewTransactionInitializationRequest(user AuthenticationInitializationUser) (request *TransactionInitializationRequest) {
	request = &TransactionInitializationRequest{
		User: user.User,
	}
	return request
}

func (request *TransactionInitializationRequest) WithTransaction(txt string) *TransactionInitializationRequest {
	request.Transaction = txt
	return request
}

func (request *TransactionInitializationRequest) WithAuthenticatorAttachment(authenticatorAttachment AuthenticatorAttachment) *TransactionInitializationRequest {
	request.Options.AuthenticatorAttachment = authenticatorAttachment
	return request
}

func (request *TransactionInitializationRequest) WithUserVerification(userVerificationRequirement UserVerificationRequirement) *TransactionInitializationRequest {
	request.Options.UserVerification = userVerificationRequirement
	return request
}

type TransactionInitializationResponse struct {
	AuthenticationInitializationResponse
}

type TransactionFinalizationResponse struct {
	AuthenticationFinalizationResponse
}

type TransactionFinalizationRequest struct {
	AuthenticationFinalizationRequest
}

type CredentialQuery struct {
	UserId   string `json:"user_id" url:"user_id"`
	PageSize uint   `json:"page_size" url:"page_size"`
	Page     uint   `json:"page" url:"page"`
}

func NewCredentialQuery() *CredentialQuery {
	return &CredentialQuery{}
}

func (c *CredentialQuery) WithUserId(id string) *CredentialQuery {
	c.UserId = id
	return c
}

func (c *CredentialQuery) WithPageSize(pageSize uint) *CredentialQuery {
	c.PageSize = pageSize
	return c
}

func (c *CredentialQuery) WithPage(page uint) *CredentialQuery {
	c.Page = page
	return c
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

func NewCredentialUpdateRequest() *CredentialUpdateRequest {
	return &CredentialUpdateRequest{}
}

func (c *CredentialUpdateRequest) WithName(name string) *CredentialUpdateRequest {
	c.Name = name
	return c
}

type AuthenticatorAttachment string

const (
	Platform      AuthenticatorAttachment = "platform"
	CrossPlatform AuthenticatorAttachment = "cross-platform"
)

type UserVerificationRequirement string

const (
	VerificationRequired    UserVerificationRequirement = "required"
	VerificationPreferred   UserVerificationRequirement = "preferred"
	VerificationDiscouraged UserVerificationRequirement = "discouraged"
)

type ConveyancePreference string

const (
	PreferNoAttestation       ConveyancePreference = "none"
	PreferIndirectAttestation ConveyancePreference = "indirect"
	PreferDirectAttestation   ConveyancePreference = "direct"
)
