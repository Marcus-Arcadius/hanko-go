package webauthn

import (
	"github.com/teamhanko/hanko-sdk-golang/client"
	"reflect"
	"testing"
)

func TestWebauthn_NewRegistrationInitializationUser(t *testing.T) {
	var tests = []struct {
		name     string
		test     RegistrationInitializationUser
		expected RegistrationInitializationUser
	}{
		{
			name: "init object",
			test: NewRegistrationInitializationUser("id", "name"),
			expected: RegistrationInitializationUser{User: client.User{
				ID:          "id",
				Name:        "name",
				DisplayName: "name",
			}},
		},
		{
			name: "init object with options",
			test: NewRegistrationInitializationUser("id", "name").WithDisplayName("display_name"),
			expected: RegistrationInitializationUser{User: client.User{
				ID:          "id",
				Name:        "name",
				DisplayName: "display_name",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			equal := reflect.DeepEqual(tt.expected, tt.test)
			if !equal {
				t.Errorf("got %+v, want %+v", tt.test, tt.expected)
			}
		})
	}
}

func TestWebauthn_NewAuthenticationInitializationUser(t *testing.T) {
	var tests = []struct {
		name     string
		test     AuthenticationInitializationUser
		expected AuthenticationInitializationUser
	}{
		{
			name: "init object",
			test: NewAuthenticationInitializationUser("id"),
			expected: AuthenticationInitializationUser{User: client.User{
				ID:          "id",
				Name:        "",
				DisplayName: "",
			}},
		},
		{
			name: "init object with options",
			test: NewAuthenticationInitializationUser("id"),
			expected: AuthenticationInitializationUser{User: client.User{
				ID:          "id",
				Name:        "",
				DisplayName: "",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			equal := reflect.DeepEqual(tt.expected, tt.test)
			if !equal {
				t.Errorf("got %+v, want %+v", tt.test, tt.expected)
			}
		})
	}
}

func TestWebauthn_NewRegistrationInitializationRequest(t *testing.T) {
	var tests = []struct {
		name     string
		test     *RegistrationInitializationRequest
		expected *RegistrationInitializationRequest
	}{
		{
			name: "init object",
			test: NewRegistrationInitializationRequest(NewRegistrationInitializationUser("id", "name")),
			expected: &RegistrationInitializationRequest{
				User: client.User{
					ID:          "id",
					Name:        "name",
					DisplayName: "name",
				},
				Options: RegistrationInitializationRequestOptions{
					AuthenticatorSelection: nil,
					ConveyancePreference:   "",
				},
			},
		},
		{
			name: "init object with options",
			test: NewRegistrationInitializationRequest(NewRegistrationInitializationUser("id", "name")).
				WithConveyancePreference(PreferIndirectAttestation).WithAuthenticatorSelection(&AuthenticatorSelection{
				AuthenticatorAttachment: Platform,
				RequireResidentKey:      true,
				UserVerification:        VerificationPreferred,
			}),
			expected: &RegistrationInitializationRequest{
				User: client.User{
					ID:          "id",
					Name:        "name",
					DisplayName: "name",
				},
				Options: RegistrationInitializationRequestOptions{
					AuthenticatorSelection: &AuthenticatorSelection{
						AuthenticatorAttachment: Platform,
						RequireResidentKey:      true,
						UserVerification:        VerificationPreferred,
					},
					ConveyancePreference: PreferIndirectAttestation,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			equal := reflect.DeepEqual(tt.expected, tt.test)
			if !equal {
				t.Errorf("got %+v, want %+v", tt.test, tt.expected)
			}
		})
	}
}

func TestWebauthn_NewAuthenticatorSelection(t *testing.T) {
	var tests = []struct {
		name     string
		test     *AuthenticatorSelection
		expected *AuthenticatorSelection
	}{
		{
			name: "init object",
			test: NewAuthenticatorSelection(),
			expected: &AuthenticatorSelection{
				AuthenticatorAttachment: "",
				RequireResidentKey:      false,
				UserVerification:        "",
			},
		},
		{
			name: "init object with options",
			test: NewAuthenticatorSelection().WithRequireResidentKey(true).
				WithAuthenticatorAttachment(CrossPlatform).WithUserVerification(VerificationRequired),
			expected: &AuthenticatorSelection{
				AuthenticatorAttachment: CrossPlatform,
				RequireResidentKey:      true,
				UserVerification:        VerificationRequired,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			equal := reflect.DeepEqual(tt.expected, tt.test)
			if !equal {
				t.Errorf("got %+v, want %+v", tt.test, tt.expected)
			}
		})
	}
}

func TestWebauthn_NewAuthenticationInitializationRequest(t *testing.T) {
	var tests = []struct {
		name     string
		test     *AuthenticationInitializationRequest
		expected *AuthenticationInitializationRequest
	}{
		{
			name: "init object",
			test: NewAuthenticationInitializationRequest().WithUser(NewAuthenticationInitializationUser("id")),
			expected: &AuthenticationInitializationRequest{
				User: client.User{
					ID:          "id",
					Name:        "",
					DisplayName: "",
				},
				Options: AuthenticationInitializationRequestOptions{
					UserVerification:        "",
					AuthenticatorAttachment: "",
				},
			},
		},
		{
			name: "init object with options",
			test: NewAuthenticationInitializationRequest().WithUser(NewAuthenticationInitializationUser("id")).
				WithUserVerification(VerificationDiscouraged).WithAuthenticatorAttachment(CrossPlatform),
			expected: &AuthenticationInitializationRequest{
				User: client.User{
					ID:          "id",
					Name:        "",
					DisplayName: "",
				},
				Options: AuthenticationInitializationRequestOptions{
					UserVerification:        VerificationDiscouraged,
					AuthenticatorAttachment: CrossPlatform,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			equal := reflect.DeepEqual(tt.expected, tt.test)
			if !equal {
				t.Errorf("got %+v, want %+v", tt.test, tt.expected)
			}
		})
	}
}

func TestWebauthn_NewTransactionInitializationRequest(t *testing.T) {
	var tests = []struct {
		name     string
		test     *TransactionInitializationRequest
		expected *TransactionInitializationRequest
	}{
		{
			name: "init object",
			test: NewTransactionInitializationRequest(NewAuthenticationInitializationUser("id")),
			expected: &TransactionInitializationRequest{
				User:        client.User{
					ID:          "id",
					Name:        "",
					DisplayName: "",
				},
				Options:     AuthenticationInitializationRequestOptions{},
				Transaction: "",
			},
		},
		{
			name: "init object with options",
			test: NewTransactionInitializationRequest(NewAuthenticationInitializationUser("id")).
				WithTransaction("transaction").WithUserVerification(VerificationDiscouraged).
				WithAuthenticatorAttachment(Platform),
			expected: &TransactionInitializationRequest{
				User: client.User{
					ID:          "id",
					Name:        "",
					DisplayName: "",
				},
				Options: AuthenticationInitializationRequestOptions{
					UserVerification:        VerificationDiscouraged,
					AuthenticatorAttachment: Platform,
				},
				Transaction: "transaction",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			equal := reflect.DeepEqual(tt.expected, tt.test)
			if !equal {
				t.Errorf("got %+v, want %+v", tt.test, tt.expected)
			}
		})
	}
}

func TestWebauthn_NewCredentialQuery(t *testing.T) {
	var tests = []struct {
		name     string
		test     *CredentialQuery
		expected *CredentialQuery
	}{
		{
			name: "init object",
			test: NewCredentialQuery(),
			expected: &CredentialQuery{},
		},
		{
			name: "init object with options",
			test: NewCredentialQuery().WithPage(3).WithPageSize(10).WithUserId("user_id"),
			expected: &CredentialQuery{
				UserId:   "user_id",
				PageSize: 10,
				Page:     3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			equal := reflect.DeepEqual(tt.expected, tt.test)
			if !equal {
				t.Errorf("got %+v, want %+v", tt.test, tt.expected)
			}
		})
	}
}

func TestWebauthn_NewCredentialUpdateRequest(t *testing.T) {
	var tests = []struct {
		name     string
		test     *CredentialUpdateRequest
		expected *CredentialUpdateRequest
	}{
		{
			name: "init object",
			test: NewCredentialUpdateRequest(),
			expected: &CredentialUpdateRequest{},
		},
		{
			name: "init object with options",
			test: NewCredentialUpdateRequest().WithName("test"),
			expected: &CredentialUpdateRequest{Name: "test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			equal := reflect.DeepEqual(tt.expected, tt.test)
			if !equal {
				t.Errorf("got %+v, want %+v", tt.test, tt.expected)
			}
		})
	}
}
