package controller

import (
	"fmt"

	"github.com/quincycheng/summon-identity-provider/internal"
	"github.com/quincycheng/summon-identity-provider/pkg/api"
	"github.com/tidwall/gjson"
)

// return value, error
func GetSecuredItem(securedItemName string) (string, error) {

	ring := internal.GetKeyring()
	cookie := internal.GetCookie(ring)
	fqdn := internal.GetPodFqdn(ring)
	theKey := ""

	theList, err := api.GetSecuredItemsList(fqdn, cookie)
	if err != nil {
		return "", err
	}

	for _, item := range gjson.Get(theList, "Result.SecuredItems").Array() {
		if item.Get("DisplayName").String() == securedItemName {
			theKey = item.Get("ItemKey").String()
		}
	}

	if theKey == "" {
		return "", fmt.Errorf("item not found: %s", securedItemName)
	}

	theItem, err := api.GetCredsForSecuredItem(fqdn, cookie, theKey)
	if err != nil {
		return "", err
	}

	return gjson.Get(theItem, "Result.n").String(), nil
}
