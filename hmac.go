package hankoApiClient

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type HmacMessageData struct {
	apiSecret     string
	apiKeyId      string
	requestMethod string
	requestPath   string
	requestBody   string
}

type HmacJson struct {
	ApiKeyId	string `json:"apiKeyId"`
	Time		string `json:"time"`
	Nonce		string `json:"nonce"`
	Signature	string `json:"signature"`
}

func CalculateHmac(messageData *HmacMessageData) string {
	date := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	nonce := uuid.New().String()

	message := fmt.Sprintf("%s:%s:%s:%s:%s",
		messageData.apiKeyId,
		date,
		messageData.requestMethod,
		messageData.requestPath,
		nonce,
	)

	if messageData.requestBody != "" {
		hashedRequestBody := hashRequestBody(messageData.requestBody)
		message = fmt.Sprintf("%s:%s", message, hashedRequestBody)
	}

	mac := hmac.New(sha256.New, []byte(messageData.apiSecret))
	mac.Write([]byte(message))
	signatureHex := hex.EncodeToString(mac.Sum(nil))

	hmacJson := &HmacJson{
		ApiKeyId:  messageData.apiKeyId,
		Time:      date,
		Nonce:     nonce,
		Signature: signatureHex,
	}

	jsonRepresentation, _ := json.Marshal(hmacJson)

	return base64.StdEncoding.EncodeToString(jsonRepresentation)
}

func hashRequestBody(requestBody string) string {
	digest := sha256.New()
	digest.Write([]byte(requestBody))
	return hex.EncodeToString(digest.Sum(nil))
}