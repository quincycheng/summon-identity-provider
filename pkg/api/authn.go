package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

type AuthInfo struct {
	TenantId string
	User     string
	Version  string
}

type UPInfo struct {
	SessionId   string
	MechanismId string
	Action      string
	Answer      string
}

// Start Authn Process
// https://identity-developer.cyberark.com/docs/starting-the-authentication-process
func StartAuthnWithTenantId(tenantId string, username string, previousAuth string) (string, string, error) {
	fqdn := fmt.Sprintf("%s%s", tenantId, IdentityRootUrlSuffix)
	result, cookie, err := StartAuthnWithFqdn(fqdn, username, previousAuth)
	if err != nil {
		return "", "", err
	}

	//Check if PodFqdn is returned
	PodFqdn := gjson.Get(result, "PodFqdn").String()
	if PodFqdn != "" {
		result, cookie, err = StartAuthnWithFqdn(PodFqdn, username, previousAuth)
		if err != nil {
			return "", "", err
		}
	}
	return result, cookie, nil
}

func StartAuthnWithFqdn(fqdn string, username string, previousAuth string) (string, string, error) {

	authInfo := AuthInfo{
		TenantId: fqdn,
		User:     username,
		Version:  IdentityApiVersion,
	}
	reqBody, _ := json.Marshal(authInfo)

	c := http.Client{Timeout: time.Duration(HttpTimeout) * time.Second}

	req, err := http.NewRequest("POST",
		fmt.Sprintf("https://%s%s", fqdn, IdentityStartAuthnUrl),
		strings.NewReader(string(reqBody)))

	if err != nil {
		panic(err)
	}

	if previousAuth != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", previousAuth))
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-IDAP-NATIVE-CLIENT", "true")
	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("error %s", err)
		return "", "", err
	}
	defer resp.Body.Close()
	responseBody, responseErr := ioutil.ReadAll(resp.Body)

	if responseErr != nil {
		fmt.Printf("response error %s", responseErr)
		return "", "", err
	}

	theCookie := ""

	for _, cookie := range resp.Cookies() {
		if cookie.Name == ".ASPXAUTH" {
			theCookie = cookie.Value
		}
	}
	return string(responseBody), theCookie, nil
}

// Adv Authn - Text
// returns responseBody, cookie, error
func AdvAuthnText(fqdn string, mechanismId string, sessionId string, cookie string, password string) (string, string, error) {

	theBody := UPInfo{
		SessionId:   sessionId,
		MechanismId: mechanismId,
		Action:      "Answer",
		Answer:      password,
	}
	reqBody, _ := json.Marshal(theBody)

	c := http.Client{Timeout: time.Duration(HttpTimeout) * time.Second}

	req, err := http.NewRequest("POST",
		fmt.Sprintf("https://%s%s", fqdn, IdentityAdvAuthnUrl),
		strings.NewReader(string(reqBody)))

	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", cookie))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-IDAP-NATIVE-CLIENT", "true")
	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("error %s", err)
		return "", "", err
	}
	defer resp.Body.Close()
	responseBody, responseErr := ioutil.ReadAll(resp.Body)

	if responseErr != nil {
		fmt.Printf("response error %s", responseErr)
		return "", "", err
	}

	theCookie := ""

	for _, cookie := range resp.Cookies() {
		if cookie.Name == ".ASPXAUTH" {
			theCookie = cookie.Value
		}
	}
	return string(responseBody), theCookie, nil
}

// Adv Authn - Send OOB
// returns responseBody, cookie, error
func AdvAuthnStartOOB(fqdn string, mechanismId string, sessionId string, cookie string) (string, string, error) {

	theBody := UPInfo{
		SessionId:   sessionId,
		MechanismId: mechanismId,
		Action:      "StartOOB",
	}
	reqBody, _ := json.Marshal(theBody)

	c := http.Client{Timeout: time.Duration(HttpTimeout) * time.Second}

	req, err := http.NewRequest("POST",
		fmt.Sprintf("https://%s%s", fqdn, IdentityAdvAuthnUrl),
		strings.NewReader(string(reqBody)))

	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", cookie))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-IDAP-NATIVE-CLIENT", "true")
	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("error %s", err)
		return "", "", err
	}
	defer resp.Body.Close()
	responseBody, responseErr := ioutil.ReadAll(resp.Body)

	if responseErr != nil {
		fmt.Printf("response error %s", responseErr)
		return "", "", err
	}

	theCookie := ""

	for _, cookie := range resp.Cookies() {
		if cookie.Name == ".ASPXAUTH" {
			theCookie = cookie.Value
		}
	}
	return string(responseBody), theCookie, nil
}

// Adv Authn - Send OOB
// returns responseBody, cookie, error
func AdvAuthnPollOOB(fqdn string, mechanismId string, sessionId string, cookie string) (string, string, error) {

	theBody := UPInfo{
		SessionId:   sessionId,
		MechanismId: mechanismId,
		Action:      "Poll",
	}
	reqBody, _ := json.Marshal(theBody)

	c := http.Client{Timeout: time.Duration(HttpTimeout) * time.Second}

	req, err := http.NewRequest("POST",
		fmt.Sprintf("https://%s%s", fqdn, IdentityAdvAuthnUrl),
		strings.NewReader(string(reqBody)))

	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", cookie))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-IDAP-NATIVE-CLIENT", "true")
	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("error %s", err)
		return "", "", err
	}
	defer resp.Body.Close()
	responseBody, responseErr := ioutil.ReadAll(resp.Body)

	if responseErr != nil {
		fmt.Printf("response error %s", responseErr)
		return "", "", err
	}

	theCookie := ""

	for _, cookie := range resp.Cookies() {
		if cookie.Name == ".ASPXAUTH" {
			theCookie = cookie.Value
		}
	}
	return string(responseBody), theCookie, nil
}
