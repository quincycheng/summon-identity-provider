package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	confirmation "github.com/erikgeiser/promptkit/confirmation"
	"github.com/erikgeiser/promptkit/selection"
	textinput "github.com/erikgeiser/promptkit/textinput"
	"github.com/quincycheng/summon-identity-provider/internal"
	"github.com/quincycheng/summon-identity-provider/pkg/api"
	"github.com/quincycheng/summon-identity-provider/pkg/controller"
	"github.com/tidwall/gjson"
)

/**************************************************
	Commands: --config
*************************************************/
func StartConfigWizard() {
	// Header
	printHeader("\nConfiguration Wizard for Summon Provider for CyberArk Identity\n")
	fmt.Println("")
	fmt.Println("This wizard will help you to configure this provider")

	// Start Config
	startConfig()

	// Want to login now?
	inputLogin := confirmation.New("Do you want to login now?", confirmation.Yes)
	ready, err := inputLogin.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	if !ready {
		fmt.Println("Please note that login is required for fetching secured item")
		os.Exit(0)
	}
	// Login
	Login()
}

func startConfig() {

	ring := internal.GetKeyring()
	configTenantId := internal.GetTenantId(ring)
	configUsername := internal.GetUsername(ring)

	if configTenantId != "" {
		fmt.Println("")
		fmt.Println("Existing configuration found.  Current settings will be overwritten")
		input := confirmation.New("Proceed?", confirmation.No)
		ready, err := input.RunPrompt()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		if !ready {
			fmt.Println("Bye!")
			os.Exit(0)
		}
		fmt.Println("")
	}

	//Tenant ID
	fmt.Println(highlightStyle.Render("1️. CyberArk Identity Tenant ID"))
	fmt.Printf("If your CyberArk Identity login URL is https://%s.my.idaptive.app, then your tenant ID is %s\n",
		highlightStyle.Render("cyberark"), highlightStyle.Render("cyberark"))

	input := textinput.New("What is your tenant ID?")
	input.InitialValue = configTenantId
	input.Placeholder = "CyberArk Identity Tenant ID cannot be empty"

	inputTenantId, err := input.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(-1)
	}
	fmt.Println("")

	//User Name
	fmt.Println(highlightStyle.Render("2. User Name"))
	fmt.Println("This is the user name for logging into CyberArk Identity")
	fmt.Println("Typical it is your email address")
	fmt.Println("")

	input = textinput.New("What is your user name?")
	input.InitialValue = configUsername
	input.Placeholder = "CyberArk Identity user name cannot be empty"

	inputUsername, err := input.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(-1)
	}
	fmt.Println("")

	// Get PodFQDN
	fmt.Println(highlightStyle.Render("3. Fully qualified domain name"))

	s := spinner.New(spinner.CharSets[35], 100*time.Millisecond)
	s.Suffix = " Verifying tenant URL..."
	s.Start()

	inputFqdn, err := controller.GetPodFqnd(inputTenantId, inputUsername)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(-1)
	}
	s.FinalMSG = fmt.Sprintf("Fully qualified domain name of your tenant is %s\n\n",
		highlightStyle.Render(inputFqdn))

	s.Stop()

	// List Challanages
	fmt.Println(highlightStyle.Render("4. Challenges"))

	s.Suffix = " Getting Challenges..."
	s.FinalMSG = "Choose your challenge mechanism(s)"

	s.Start()

	jsonString, _, err := api.StartAuthnWithFqdn(inputFqdn, inputUsername, "")
	if err != nil {
		fmt.Printf("\n\nError: %v\n", err)
		os.Exit(-1)
	}
	if !gjson.Get(jsonString, "success").Bool() {
		fmt.Printf("Error: API returns %s", jsonString)
		os.Exit(-1)
	}

	s.Stop()

	// Social Login
	if gjson.Get(jsonString, "Result.IdpRedirectUrl").String() != "" {
		fmt.Printf("Error: IDP redirection is not supported.  API returns %s", jsonString)
		os.Exit(-1)
	}

	inputChallenge1 := ""
	inputChallenge0 := ""

	fmt.Println()

	challenges := gjson.Get(jsonString, "Result.Challenges")
	challenges.ForEach(func(key, value gjson.Result) bool {

		sp := selection.New(fmt.Sprintf("Challenge Mechanism %d", key.Int()+1),
			selection.Choices(gjson.Get(value.String(), "Mechanisms.#.PromptSelectMech").Array()))
		sp.PageSize = 5
		sp.FilterPlaceholder = "Use ↑ ↓ keys to navigate, Press ↵ to select"
		sp.FilterPrompt = "Your choice: "

		choice, err := sp.RunPrompt()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		if key.Int() == 0 {
			inputChallenge0 = choice.String
		}
		if key.Int() == 1 {
			inputChallenge1 = choice.String
		}
		return true // keep iterating
	})

	// Confirm to Save
	fmt.Println()
	inputSave := confirmation.New("Save Setttings?", confirmation.No)
	ready, err := inputSave.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	if !ready {
		fmt.Println("Bye!")
		os.Exit(0)
	}
	s.Suffix = " Saving configuration..."
	s.FinalMSG = "Configuration saved"
	s.Start()

	internal.SetTenantId(inputTenantId, ring)
	internal.SetUsername(inputUsername, ring)
	internal.SetPodFqdn(inputFqdn, ring)
	internal.SetChallenge0(inputChallenge0, ring)
	internal.SetChallenge1(inputChallenge1, ring)

	internal.SetCookie("", ring)

	s.Stop()
	fmt.Print("\n\n")
}
