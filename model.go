package hankoApiClient

type Operation string

const (
	AUTH  Operation = "AUTH"
	REG   Operation = "REG"
	DEREG Operation = "DEREG"
)

type ClientData struct {
	RemoteAddress string `json:"remoteAddress"`
	UserAgent     string `json:"userAgent"`
}

type AuthenticatorSelectionCriteria struct {
	UserVerification        string `json:"userVerification"`
	AuthenticatorAttachment string `json:"authenticatorAttachment"`
}

type Request struct {
	Operation                      Operation                       `json:"operation"`
	Username                       string                          `json:"username"`
	UserId                         string                       `json:"userId"`
	ClientData                     *ClientData                     `json:"clientData"`
	DeviceIds                      *[]string                       `json:"deviceIds"`
	AuthenticatorSelectionCriteria *AuthenticatorSelectionCriteria `json:"authenticatorSelectionCriteria"`
}

type RelyingParty struct {
	AppId                       string `json:"appId"`
	AuthenticationTimeoutSecond int    `json:"authenticationTimeoutSeconds"`
	BasicIntegrity              bool   `json:"basicIntegrity"`
	CtsProfileMatch             bool   `json:"ctsProfileMatch"`
	Icon                        string `json:"icon"`
	Id                          string `json:"id"`
	Jailbreak                   bool   `json:"jailbreak"`
	Name                        string `json:"name"`
	RegistrationTimeoutSeconds  int    `json:"registrationTimeoutSeconds"`
	ShowLocation                bool   `json:"showLocation"`
}

type Link struct {
	Href   string `json:"href"`
	Method string `json:"method"`
	Rel    string `json:"rel"`
}

type Response struct {
	Id           string       `json:"id"`
	Operation    Operation    `json:"operation"`
	Username     string       `json:"username"`
	UserId       string       `json:"userId"`
	Status       string       `json:"status"`
	CreatedAt    string       `json:"createdAt"`
	ValidUntil   string       `json:"validUntil"`
	RelyingParty RelyingParty `json:"relyingParty"`
	Request      string       `json:"request"`
	DeviceId     string       `json:"deviceId"`
	Links        []Link       `json:"links"`
}

type AuthenticatorResponse struct {
	AttestationObject string `json:"attestationObject,omitempty"`
	ClientDataJson    string `json:"clientDataJSON"`
	AuthenticatorData string `json:"authenticatorData,omitempty"`
	Signature         string `json:"signature,omitempty"`
	UserHandle        string `json:"userHandle,omitempty"`
}

type PublicKeyCredential struct {
	Id       string                `json:"id"`
	RawId    string                `json:"rawId"`
	Type     string                `json:"type"`
	Response AuthenticatorResponse `json:"response"`
}

type DeviceKeyInfo struct {
	KeyName string `json:"keyName"`
}

type HankoCredentialRequest struct {
	WebAuthnResponse PublicKeyCredential `json:"webAuthnResponse"`
	DeviceKeyInfo    DeviceKeyInfo       `json:"deviceKeyInfo"`
}
