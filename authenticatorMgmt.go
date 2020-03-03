package hankoApiClient

import (
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"time"
)

// AuthenticatorDevice holds information about an Authenticator Device
type AuthenticatorDevice struct {
	DeviceId          string    `json:"deviceId"`
	KeyName           string    `json:"keyName"`
	AuthenticatorType string    `json:"authenticatorType"`
	AuthenticatorAttachment string `json:"authenticatorAttachment"`
	LastUsage         time.Time `json:"lastUsage"`
	CreatedAt         time.Time `json:"createdAt"`
}

// AuthenticatorDevices Slice of AuthenticatorDevice
type AuthenticatorDevices []AuthenticatorDevice

// AuthenticatorDevices Implements sort.Interface by LastUsage
func (a AuthenticatorDevices) Len() int           { return len(a) }
func (a AuthenticatorDevices) Less(i, j int) bool { return a[i].LastUsage.After(a[j].LastUsage) }
func (a AuthenticatorDevices) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// FilterDuplicateTypes Filters out duplicate types e.g. WebAuthn or UAF
func (a AuthenticatorDevices) FilterDuplicateTypes() AuthenticatorDevices {
	types := map[string]map[string]bool{}
	n := 0
	var aa AuthenticatorDevices
	for _, v := range a {
		if !types[v.AuthenticatorType][v.AuthenticatorAttachment] {
			aa = append(aa, v)
			types[v.AuthenticatorType] = map[string]bool{}
			types[v.AuthenticatorType][v.AuthenticatorAttachment] = true
			n++
		}
	}
	return aa[:n]
}

// AuthenticatorRename struct for Authenticator renaming
type AuthenticatorRename struct {
	NewName string `json:"newName"`
}

// GetAuthenticators returns a list of Authenticators for a given user
func (client *HankoApiClient) GetAuthenticators(userId string) (*AuthenticatorDevices, error) {
	resp, err := client.doRequest(http.MethodGet, "/mgmt/v1/registrations/"+userId, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Could not do Request to get Authenticator Devices.")
	}
	apiResp := &AuthenticatorDevices{}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(apiResp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode the response")
	}

	return apiResp, nil
}

// RenameAuthenticator renames an Authenticator with the given Id
func (client *HankoApiClient) RenameAuthenticator(deviceId string, rename AuthenticatorRename) (*AuthenticatorRename, error) {
	resp, err := client.doRequest(http.MethodPost, "/mgmt/v1/registrations/rename/"+deviceId, rename)

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("hankoApi returned status code: %d", resp.StatusCode)
	}
	dec := json.NewDecoder(resp.Body)
	body := &AuthenticatorRename{}
	err = dec.Decode(body)
	if err != nil {
		log.Println(err)
		return nil, errors.Wrap(err, "Could not Decode Response")
	}
	return body, nil
}
