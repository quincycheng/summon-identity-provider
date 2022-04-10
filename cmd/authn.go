package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	textinput "github.com/erikgeiser/promptkit/textinput"
	"github.com/quincycheng/summon-identity-provider/internal"
	"github.com/quincycheng/summon-identity-provider/pkg/api"
	"github.com/quincycheng/summon-identity-provider/pkg/controller"

	"github.com/tidwall/gjson"
)

const msgLoginSuccessful = "✔️  Login Successful"

/**************************************************
	Commands: --config
*************************************************/
func StartLogin() {
	// Header
	printHeader("\nSummon Provider for CyberArk Identity\n")
	Login()
}

func Login() {
	// Check if authn token (cookie) works
	ring := internal.GetKeyring()
	cookie := internal.GetCookie(ring)
	fqdn := internal.GetPodFqdn(ring)

	result, newCookie, err := api.StartAuthnWithFqdn(fqdn, internal.GetUsername(ring), cookie)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	isAuthn := gjson.Get(result, "Result.Auth")
	if isAuthn.Exists() {
		// Login Successful
		fmt.Println(msgLoginSuccessful)
		return
	}

	sessionId := gjson.Get(result, "Result.SessionId").String()
	// Prompt for Challenges
	challenge0 := internal.GetChallenge0(ring)
	challenge1 := internal.GetChallenge1(ring)

	if challenge0 != "" {
		fmt.Println("")
		mech0, err := controller.GetMechanismTypeByValue(gjson.Get(result, "Result.Challenges.0.Mechanisms"), challenge0)

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		isAuthnOk, newCookie, err := continueAuth(fqdn, sessionId, newCookie, mech0)

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if isAuthnOk {
			internal.SetCookie(newCookie, ring)
			fmt.Println(msgLoginSuccessful)
			return
		}
	}
	if challenge1 != "" {

		mech1, err := controller.GetMechanismTypeByValue(gjson.Get(result, "Result.Challenges.1.Mechanisms"), challenge1)

		if err != nil {
			fmt.Printf("Error: %v\n", err.Error())
			os.Exit(1)
		}

		fmt.Println()
		isAuthnOk, newCookie, err := continueAuth(fqdn, sessionId, newCookie, mech1)

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if isAuthnOk {
			internal.SetCookie(newCookie, ring)
			fmt.Println(msgLoginSuccessful)
			return
		}
	}

	// Uncaught conditions
	fmt.Printf("Error: Login failed.   Please make sure your input info & configuration are correct")
	os.Exit(1)
}

// returns isAuthnOK, new ocokie, error
func continueAuth(fqdn string, sessionId string, cookie string, mech controller.MechType) (bool, string, error) {

	fmt.Println(highlightStyle.Render(mech.PromptSelectMech))
	theCookie := cookie

	switch mech.AnswerType {
	case "Text":
		response, theCookie, err := api.AdvAuthnText(fqdn,
			mech.MechanismId, sessionId, cookie,
			inputAnswerText(mech.PromptMechChosen))

		if err != nil {
			fmt.Println(err)
			return false, theCookie, err
		}
		if !gjson.Get(response, "success").Bool() {
			return false, theCookie, errors.New(gjson.Get(response, "Message").String())
		}

		return gjson.Get(response, "Result.Auth").Exists(), theCookie, err

	case "StartTextOob":
		_, theCookie, err := api.AdvAuthnStartOOB(fqdn, mech.MechanismId, sessionId, cookie)
		if err != nil {
			return false, theCookie, err
		}
		theCode := inputAnswerOOB(mech.PromptMechChosen)

		if theCode != "" {
			response, theCookie, err := api.AdvAuthnText(fqdn, mech.MechanismId, sessionId, cookie, theCode)

			if err != nil {
				return false, theCookie, err
			}
			return gjson.Get(response, "Result.Auth").Exists(), theCookie, err

		} else {
			s := spinner.New(spinner.CharSets[35], 100*time.Millisecond)
			s.Suffix = " Please follow the instructions from the link"
			s.FinalMSG = "Login request approved"
			s.Start()

			for {
				response, theCookie, err := api.AdvAuthnPollOOB(fqdn, mech.MechanismId, sessionId, cookie)

				if gjson.Get(response, "Result.Summary").String() == "LoginSuccess" {
					s.Stop()
					fmt.Println()
					return gjson.Get(response, "Result.Auth").Exists(), theCookie, err
				} else {
					time.Sleep(1 * time.Second)
				}
			}

		}

	case "StartOob":
		_, theCookie, err := api.AdvAuthnStartOOB(fqdn, mech.MechanismId, sessionId, cookie)
		if err != nil {
			return false, theCookie, err
		}

		for {
			response, theCookie, err := api.AdvAuthnPollOOB(fqdn, mech.MechanismId, sessionId, cookie)
			fmt.Println(response)

			if gjson.Get(response, "Summary").String() == "LoginSuccess" {
				return gjson.Get(response, "Result.Auth").Exists(), theCookie, err
			} else {
				time.Sleep(1 * time.Second)
			}
		}
	}

	return false, theCookie, nil
}

func inputAnswerText(prompt string) string {
	input := textinput.New(fmt.Sprintf("%s:", prompt))
	input.Placeholder = "Password will NOT be stored"
	input.Validate = func(s string) bool { return len(s) > 0 } // nolint:gomnd
	input.Hidden = true

	password, err := input.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	return password
}

func inputAnswerOOB(prompt string) string {

	fmt.Println(prompt)
	input := textinput.New("Code: ")
	input.Placeholder = "If you have clicked the link in the message, press ↵ "
	input.Validate = func(s string) bool { return len(s) >= 0 } // nolint:gomnd

	input.Hidden = true

	theCode, err := input.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	return theCode
}
