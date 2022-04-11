package controller

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	"github.com/quincycheng/summon-identity-provider/internal"
	"github.com/quincycheng/summon-identity-provider/pkg/api"
	"github.com/tidwall/gjson"
)

const secureNoteType = "SecureNote"
const securePasswordType = "Password"

// GetSecuredItem
// return value, error
func GetSecuredItem(securedItemName string) (string, error) {

	ring := internal.GetKeyring()
	cookie := internal.GetCookie(ring)
	fqdn := internal.GetPodFqdn(ring)
	itemKey := ""
	itemType := ""

	theList, err := api.GetSecuredItemsList(fqdn, cookie)
	if err != nil {
		return "", err
	}

	for _, item := range gjson.Get(theList, "Result.SecuredItems").Array() {
		if item.Get("DisplayName").String() == securedItemName {
			itemKey = item.Get("ItemKey").String()
			itemType = item.Get("SecuredItemType").String()
		}
	}

	if itemKey == "" {
		return "", fmt.Errorf("item not found: %s", securedItemName)
	}

	if itemType == secureNoteType {
		theItem, err := api.GetCredsForSecuredItem(fqdn, cookie, itemKey)
		if err != nil {
			return "", err
		}
		return gjson.Get(theItem, "Result.n").String(), nil
	}

	if itemType == securePasswordType {

		// RSA-OAEP SHA-256 2048 bit
		privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
		pubk := privateKey.PublicKey
		pubkDer, _ := x509.MarshalPKIXPublicKey(&pubk)
		pubkBlock := pem.Block{
			Type:    "PUBLIC KEY",
			Headers: nil,
			Bytes:   pubkDer,
		}
		pubkPem := pem.EncodeToMemory(&pubkBlock)

		theItem, err := api.GetCredsForSecuredPassword(fqdn, cookie, itemKey, string(pubkPem))
		if err != nil {
			return "", err
		}

		encryptedBytes, err := base64.StdEncoding.DecodeString(gjson.Get(theItem, "Result.e").String())
		if err != nil {
			panic(err)
		}
		decryptedBytes, err := privateKey.Decrypt(nil, encryptedBytes, &rsa.OAEPOptions{Hash: crypto.SHA256})
		if err != nil {
			panic(err)
		}
		return gjson.Get(string(decryptedBytes), "Password").String(), nil
	}

	// Unknown type
	return "", fmt.Errorf("unsupported secured item Type: %s", itemType)

}
