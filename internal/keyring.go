package internal

import (
	"fmt"
	"os"
	"path"

	"github.com/99designs/keyring"
)

const keyringServiceName = "summon-identity-provider"
const keyringTenantId = "tenant-id"
const keyringUsername = "username"
const keyringCookie = "cookie"
const keyringPodFqdn = "PodFqdn"
const keyringChallenge0 = "Challenge0"
const keyringChallenge1 = "Challenge1"

func GetKeyring() keyring.Keyring {
	// https://pkg.go.dev/os#UserConfigDir
	// On Unix systems, it returns $XDG_CONFIG_HOME as specified by
	//   https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html if non-empty, else $HOME/.config.
	// On Darwin, it returns $HOME/Library/Application Support.
	// On Windows, it returns %AppData%. On Plan 9, it returns $home/lib.
	ucd, _ := os.UserConfigDir()

	ring, _ := keyring.Open(keyring.Config{
		ServiceName:      keyringServiceName,
		AllowedBackends:  []keyring.BackendType{keyring.FileBackend},
		FileDir:          path.Join(ucd, keyringServiceName),
		FilePasswordFunc: keyring.FixedStringPrompt(keyringServiceName),
	})
	return ring
}

func GetAvaliableBackend() []keyring.BackendType {
	return keyring.AvailableBackends()
}

func GetTenantId(ring keyring.Keyring) string {
	i, _ := ring.Get(keyringTenantId)
	return string(i.Data)
}

func SetTenantId(input string, ring keyring.Keyring) {
	err := ring.Set(keyring.Item{
		Key:  keyringTenantId,
		Data: []byte(input),
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}

func GetUsername(ring keyring.Keyring) string {
	i, _ := ring.Get(keyringUsername)
	return string(i.Data)
}

func SetUsername(input string, ring keyring.Keyring) {
	err := ring.Set(keyring.Item{
		Key:  keyringUsername,
		Data: []byte(input),
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}

func GetCookie(ring keyring.Keyring) string {
	i, _ := ring.Get(keyringCookie)
	return string(i.Data)
}

func SetCookie(input string, ring keyring.Keyring) {
	err := ring.Set(keyring.Item{
		Key:  keyringCookie,
		Data: []byte(input),
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}

func GetPodFqdn(ring keyring.Keyring) string {
	i, _ := ring.Get(keyringPodFqdn)
	return string(i.Data)
}

func SetPodFqdn(input string, ring keyring.Keyring) {
	err := ring.Set(keyring.Item{
		Key:  keyringPodFqdn,
		Data: []byte(input),
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}

func GetChallenge0(ring keyring.Keyring) string {
	i, _ := ring.Get(keyringChallenge0)
	return string(i.Data)
}

func SetChallenge0(input string, ring keyring.Keyring) {
	err := ring.Set(keyring.Item{
		Key:  keyringChallenge0,
		Data: []byte(input),
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}

func GetChallenge1(ring keyring.Keyring) string {
	i, _ := ring.Get(keyringChallenge1)
	return string(i.Data)
}

func SetChallenge1(input string, ring keyring.Keyring) {
	err := ring.Set(keyring.Item{
		Key:  keyringChallenge1,
		Data: []byte(input),
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}
