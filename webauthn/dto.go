package webauthn

import (
	hankoClient "github.com/teamhanko/hanko-sdk-golang/client"
	"github.com/teamhanko/webauthn/protocol"
	"time"
)

// RegistrationInitializationUser is the user representation used for credential registration initialization.
type RegistrationInitializationUser struct {
	hankoClient.User
}

// NewRegistrationInitializationUser creates a new RegistrationInitializationUser
func NewRegistrationInitializationUser(id string, name string) RegistrationInitializationUser {
	return RegistrationInitializationUser{hankoClient.User{
		ID:          id,
		Name:        name,
		DisplayName: name,
	}}
}

// WithDisplayName can be used to set a different displayName for the user.
//
// It is intended that the user specifies this value. If not specified the displayName defaults to the value of the name
// given on construction.
func (user RegistrationInitializationUser) WithDisplayName(displayName string) RegistrationInitializationUser {
	user.DisplayName = displayName
	return user
}

// AuthenticationInitializationUser is the user representation used for an authentication initialization.
type AuthenticationInitializationUser struct {
	hankoClient.User
}

// NewAuthenticationInitializationUser creates a new AuthenticationInitializationUser.
func NewAuthenticationInitializationUser(id string) AuthenticationInitializationUser {
	return AuthenticationInitializationUser{hankoClient.User{ID: id}}
}

// WithName allows you to set the user name.
func (user AuthenticationInitializationUser) WithName(name string) AuthenticationInitializationUser {
	user.Name = name
	return user
}

// RegistrationInitializationRequest is used to initialize a credential registration operation.
type RegistrationInitializationRequest struct {
	User    hankoClient.User                         `json:"user"`
	Options RegistrationInitializationRequestOptions `json:"options"`
}

// NewRegistrationInitializationRequest creates a new RegistrationInitializationRequest.
func NewRegistrationInitializationRequest(user RegistrationInitializationUser) *RegistrationInitializationRequest {
	return &RegistrationInitializationRequest{User: user.User}
}

// RegistrationInitializationRequestOptions allows you to set additional authenticator attributes for a registration
// initialization request.
type RegistrationInitializationRequestOptions struct {
	AuthenticatorSelection *AuthenticatorSelection `json:"authenticatorSelection"`
	ConveyancePreference   ConveyancePreference    `json:"attestation"`
}

// WithConveyancePreference allows you to set the preferred authenticator attestation ConveyancePreference.
// If not set, the value will default to "none" (PreferNoAttestation) during registration.
func (request *RegistrationInitializationRequest) WithConveyancePreference(conveyancePreference ConveyancePreference) *RegistrationInitializationRequest {
	request.Options.ConveyancePreference = conveyancePreference
	return request
}

// AuthenticatorSelection specifies the requirements regarding authenticator attributes for credential creation.
type AuthenticatorSelection struct {
	AuthenticatorAttachment AuthenticatorAttachment     `json:"authenticatorAttachment,omitempty"`
	RequireResidentKey      bool                        `json:"requireResidentKey,omitempty"`
	UserVerification        UserVerificationRequirement `json:"userVerification,omitempty"`
}

// NewAuthenticatorSelection creates a new AuthenticatorSelection used to initialize a credential registration.
func NewAuthenticatorSelection() *AuthenticatorSelection {
	return &AuthenticatorSelection{}
}

// WithAuthenticatorSelection lets you specify the requirements regarding authenticator attributes.
func (request *RegistrationInitializationRequest) WithAuthenticatorSelection(authenticatorSelection *AuthenticatorSelection) *RegistrationInitializationRequest {
	request.Options.AuthenticatorSelection = authenticatorSelection
	return request
}

// WithAuthenticatorAttachment allows you to set your preferred authenticator attachment modality.
func (authenticatorSelection *AuthenticatorSelection) WithAuthenticatorAttachment(authenticatorAttachment AuthenticatorAttachment) *AuthenticatorSelection {
	authenticatorSelection.AuthenticatorAttachment = authenticatorAttachment
	return authenticatorSelection
}

// WithRequireResidentKey allows you to specify whether a credential should be registered as a resident credential.
//
// If set to true, the user must use an authenticator device that supports resident keys.
// If not set, the value will default to "false" and the credential will not be registered
// as a resident credential.
func (authenticatorSelection *AuthenticatorSelection) WithRequireResidentKey(requireResidentKey bool) *AuthenticatorSelection {
	authenticatorSelection.RequireResidentKey = requireResidentKey
	return authenticatorSelection
}

// WithUserVerification allows you to set your UserVerificationRequirement for the credential registration.
// If not set, the value will default to "preferred" (VerificationPreferred) during registration.
func (authenticatorSelection *AuthenticatorSelection) WithUserVerification(userVerificationRequirement UserVerificationRequirement) *AuthenticatorSelection {
	authenticatorSelection.UserVerification = userVerificationRequirement
	return authenticatorSelection
}

// RegistrationInitializationResponse contains the representation of CredentialCreationOptions generated by the Hanko
// Authentication API that must be passed to browser's WebAuthn API via navigator.credentials.create() in order
// to create a credential.
//
// See also: https://www.w3.org/TR/webauthn/#sctn-credentialcreationoptions-extension
type RegistrationInitializationResponse struct {
	protocol.CredentialCreation
}

// RegistrationFinalizationRequest contains the representation of a PublicKeyCredential obtained through credential
// creation via the browsers' navigator.credentials.create().
//
// See also: https://www.w3.org/TR/webauthn-2/#publickeycredential
type RegistrationFinalizationRequest struct {
	protocol.CredentialCreationResponse
}

// RegistrationFinalizationResponse is the response when the credential registration was successful.
type RegistrationFinalizationResponse struct {
	User       hankoClient.User `json:"user"`
	Credential *Credential      `json:"credential"`
}

// Authenticator holds information about the authenticator associated with a registered credential.
type Authenticator struct {
	// The "Authenticator Attestation Globally Unique ID" is a 128-bit identifier indicating the type
	// (e.g. make and model) of the authenticator.
	Aaguid string `json:"aaguid,omitempty"`

	// A description of the authenticator.
	Name string `json:"name,omitempty"`

	// The AuthenticatorAttachment modality of the authenticator.
	Attachment string `json:"attachment,omitempty"`
}

// AuthenticationInitializationRequest is used to initialize an authentication operation.
type AuthenticationInitializationRequest struct {
	User    hankoClient.User                           `json:"user"`
	Options AuthenticationInitializationRequestOptions `json:"options"`
}

// NewAuthenticationInitializationRequest creates a new AuthenticationInitializationRequest.
func NewAuthenticationInitializationRequest() (request *AuthenticationInitializationRequest) {
	request = &AuthenticationInitializationRequest{}
	return request
}

// WithUser allows you to set an AuthenticationInitializationUser which is nessaccery to authenticate with non-resident
// keys.
func (request *AuthenticationInitializationRequest) WithUser(user AuthenticationInitializationUser) *AuthenticationInitializationRequest {
	request.User = user.User
	return request
}

// AuthenticationInitializationRequestOptions allows you to set additional authenticator attributes for the
// authentication initialization.
type AuthenticationInitializationRequestOptions struct {
	UserVerification        UserVerificationRequirement `json:"userVerification"`
	AuthenticatorAttachment AuthenticatorAttachment     `json:"authenticatorAttachment"`
}

// WithAuthenticatorAttachment allows you to set the AuthenticatorAttachment modality for the authentication
// initialization.
func (request *AuthenticationInitializationRequest) WithAuthenticatorAttachment(authenticatorAttachment AuthenticatorAttachment) *AuthenticationInitializationRequest {
	request.Options.AuthenticatorAttachment = authenticatorAttachment
	return request
}

// WithUserVerification allows you to set your UserVerification requirements for the authentication.
// If not set, the value default to "preferred" (VerificationPreferred) during authentication.
func (request *AuthenticationInitializationRequest) WithUserVerification(userVerificationRequirement UserVerificationRequirement) *AuthenticationInitializationRequest {
	request.Options.UserVerification = userVerificationRequirement
	return request
}

// AuthenticationInitializationResponse contains the representation of CredentialRequestOptions generated by the Hanko
// Authentication API that must be passed to browser's WebAuthn API via navigator.credentials.get() in order
// to authenticate with a credential/create an assertion.
//
// See also: https://www.w3.org/TR/webauthn-2/#sctn-credentialrequestoptions-extension
type AuthenticationInitializationResponse struct {
	protocol.CredentialAssertion
}

// AuthenticationFinalizationRequest contains the representation of a PublicKeyCredential obtained through assertion
// generation via the browsers' navigator.credentials.get().
//
// See also: https://www.w3.org/TR/webauthn-2/#publickeycredential
type AuthenticationFinalizationRequest struct {
	protocol.CredentialAssertionResponse
}

// AuthenticationFinalizationResponse is the response when the authentication was successful.
type AuthenticationFinalizationResponse struct {
	User hankoClient.User `json:"user"`
}

// TransactionInitializationRequest is used to initialize a transaction operation.
type TransactionInitializationRequest struct {
	User        hankoClient.User                           `json:"user"`
	Options     AuthenticationInitializationRequestOptions `json:"options"`
	Transaction string                                     `json:"transaction"`
}

// NewTransactionInitializationRequest creates a new TransactionInitializationRequest.
func NewTransactionInitializationRequest(user AuthenticationInitializationUser) (request *TransactionInitializationRequest) {
	request = &TransactionInitializationRequest{
		User: user.User,
	}
	return request
}

// WithTransaction specifies the transaction text to be verified by the user.
func (request *TransactionInitializationRequest) WithTransaction(txt string) *TransactionInitializationRequest {
	request.Transaction = txt
	return request
}

// WithAuthenticatorAttachment allows you to set your preferred authenticator attachment modality for the transaction
// initialization.
func (request *TransactionInitializationRequest) WithAuthenticatorAttachment(authenticatorAttachment AuthenticatorAttachment) *TransactionInitializationRequest {
	request.Options.AuthenticatorAttachment = authenticatorAttachment
	return request
}

// WithUserVerification allows you to set your UserVerification requirements for the transaction.
// If not set, the value default to "preferred" (VerificationPreferred) during authentication.
func (request *TransactionInitializationRequest) WithUserVerification(userVerificationRequirement UserVerificationRequirement) *TransactionInitializationRequest {
	request.Options.UserVerification = userVerificationRequirement
	return request
}

// TransactionInitializationResponse contains the representation of CredentialRequestOptions generated by the Hanko
// Authentication API that must be passed to browser's WebAuthn API via navigator.credentials.get() in order
// to authenticate with a credential/create an assertion.
//
// See also: https://www.w3.org/TR/webauthn-2/#sctn-credentialrequestoptions-extension
type TransactionInitializationResponse struct {
	AuthenticationInitializationResponse
}

// TransactionFinalizationResponse is the response when the transaction was successful.
type TransactionFinalizationResponse struct {
	AuthenticationFinalizationResponse
}

// TransactionFinalizationRequest contains the representation of a PublicKeyCredential obtained through assertion
// generation via the browsers' navigator.credentials.get().
//
// See also: https://www.w3.org/TR/webauthn-2/#publickeycredential
type TransactionFinalizationRequest struct {
	AuthenticationFinalizationRequest
}

// CredentialQuery is used to search for credentials.
type CredentialQuery struct {
	UserId   string `json:"user_id" url:"user_id"`     // The user ID to filter by.
	PageSize uint   `json:"page_size" url:"page_size"` // The page size of the returned result set.
	Page     uint   `json:"page" url:"page"`           // The desired page to return from the result set.
}

// NewCredentialQuery creates a new CredentialQuery that can be used to filter and paginate results when
// retrieving credentials.
func NewCredentialQuery() *CredentialQuery {
	return &CredentialQuery{}
}

// WithUserId allows you to filter credentials by the specified user id.
func (c *CredentialQuery) WithUserId(id string) *CredentialQuery {
	c.UserId = id
	return c
}

// WithPageSize allows you to specify the page size of the result set.
func (c *CredentialQuery) WithPageSize(pageSize uint) *CredentialQuery {
	c.PageSize = pageSize
	return c
}

// WithPage allows you to specify which page of the result set to return.
func (c *CredentialQuery) WithPage(page uint) *CredentialQuery {
	c.Page = page
	return c
}

// Credential represents a credential.
type Credential struct {
	// ID of the credential.
	Id string `json:"id"`

	// Time of credential creation in date-time notation as defined by RFC 3339, section 5.6.
	CreatedAt time.Time `json:"createdAt"`

	// Last time this credential was used for authentication in date-time notation as defined by RFC 3339, section 5.6.
	LastUsed time.Time `json:"lastUsed"`

	// A human-palatable name for the credential that has a default value on credential creation but can be changed later
	// through the Update credential operation.
	Name string `json:"name"`

	// Indicates whether this credential was registered with a successful user verification process.
	UserVerification bool `json:"userVerification"`

	// Indicates whether this credential was registered as a resident credential/client-side discoverable credential.
	IsResidentKey bool `json:"isResidentKey"`

	// Representation of the authenticator used for registering the credential.
	Authenticator *Authenticator `json:"authenticator,omitempty"`
}

// CredentialUpdateRequest is used to update an existing credential.
type CredentialUpdateRequest struct {
	Name string `json:"name"`
}

// NewCredentialUpdateRequest creates a CredentialUpdateRequest for updating the credential. Currently, it is only
// possible to update the credential name.
func NewCredentialUpdateRequest() *CredentialUpdateRequest {
	return &CredentialUpdateRequest{}
}

// WithName allows you to specify the new name of credential to update.
func (c *CredentialUpdateRequest) WithName(name string) *CredentialUpdateRequest {
	c.Name = name
	return c
}

// WebAuthn relying parties use this to express a preferred authenticator attachment modality when calling
// navigator.credentials.create() to create a credential.
//
// See also: https://www.w3.org/TR/webauthn/#enum-attachment
type AuthenticatorAttachment string

const (
	// Indicates that the authenticator should be a platform authenticator
	Platform AuthenticatorAttachment = "platform"

	// Indicates that the authenticator should be a roaming authenticator
	CrossPlatform AuthenticatorAttachment = "cross-platform"
)

// A WebAuthn relying party may require user verification for some of its operations but not for others, and may use
// this type to express its needs.
//
// See also: https://www.w3.org/TR/webauthn/#enum-userVerificationRequirement
type UserVerificationRequirement string

const (
	// Indicates that user verification is required for the operation. the operation will fail if the
	// authenticator response does not have the user verification flag set.
	VerificationRequired UserVerificationRequirement = "required"

	// Indicates that the relying party prefers user verification for the operation if possible, but will
	//not fail the operation if the response does not have the user verification flag set.
	VerificationPreferred UserVerificationRequirement = "preferred"

	// Indicates that the relying party does not want user verification employed during the operation.
	VerificationDiscouraged UserVerificationRequirement = "discouraged"
)

// WebAuthn relying parties may use ConveyancePreference to specify their preference regarding attestation
// conveyance during credential generation.
//
// See also: https://www.w3.org/TR/webauthn/#enum-attestation-convey
type ConveyancePreference string

const (
	// Indicates that the relying party is not interested in authenticator attestation.
	PreferNoAttestation ConveyancePreference = "none"

	// Indicates that the relying party prefers an attestation conveyance yielding verifiable attestation
	// statements, but allows the client to decide how to obtain such attestation statements. Note: There is
	// no guarantee that the RP will obtain a verifiable attestation statement in this case. For example, in
	// the case that the authenticator employs self attestation.
	PreferIndirectAttestation ConveyancePreference = "indirect"

	// Indicates that the relying party wants to receive the attestation statement as generated by the
	// authenticator.
	PreferDirectAttestation ConveyancePreference = "direct"
)
