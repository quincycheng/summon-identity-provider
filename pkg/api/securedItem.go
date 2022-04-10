package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// return list, error
func GetSecuredItemsList(fqdn string, cookie string) (string, error) {

	c := http.Client{Timeout: time.Duration(HttpTimeout) * time.Second}

	req, err := http.NewRequest("POST",
		fmt.Sprintf("https://%s%s", fqdn, IdentityListSecuredItemsUrl), nil)

	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", cookie))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-IDAP-NATIVE-CLIENT", "true")
	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("error %s", err)
		return "", err
	}
	defer resp.Body.Close()
	responseBody, responseErr := ioutil.ReadAll(resp.Body)

	if responseErr != nil {
		fmt.Printf("response error %s", responseErr)
		return "", err
	}

	return string(responseBody), nil
}

// return list, error
func GetCredsForSecuredItem(fqdn string, cookie string, itemKey string) (string, error) {

	c := http.Client{Timeout: time.Duration(HttpTimeout) * time.Second}

	req, err := http.NewRequest("POST",
		fmt.Sprintf("https://%s%s?sItemKey=%s", fqdn, IdentityGetCredsForSecuredItemUrl, itemKey), nil)

	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", cookie))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-IDAP-NATIVE-CLIENT", "true")
	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("error %s", err)
		return "", err
	}
	defer resp.Body.Close()
	responseBody, responseErr := ioutil.ReadAll(resp.Body)

	if responseErr != nil {
		fmt.Printf("response error %s", responseErr)
		return "", err
	}

	return string(responseBody), nil
}
