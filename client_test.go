package hankoApiClient

import (
	"github.com/google/uuid"
	"testing"
)

var apiHost = "https://api.dev.hanko.io/"
var apiSecret = "17a1b9585cc92782d6017324c77887b283427e8076a2e775dbd7570"

func TestHankoApiClient_Request(t *testing.T) {
	client := NewHankoApiClient(apiHost, apiSecret)
	res, err := client.InitWebauthnRegistration(uuid.New().String(), "testuser@hanko.io")
	if err != nil {
		t.Fail()
	}
	if res.Status != "PENDING" {
		t.Fail()
	}
}
