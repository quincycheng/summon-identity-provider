package main

import (
	"os"

	"github.com/quincycheng/summon-identity-provider/cmd"
)

/**
summon-identity-provider
summon-identity-provider <secured item name>
summon-identity-provider --config
summon-identity-provider --auth
**/

func main() {

	// Parsing Commands
	if len(os.Args) != 2 {
		cmd.PrintIntroMessage()
	} else {
		switch os.Args[1] {
		case "--config":
			cmd.StartConfigWizard()
		case "--login":
			cmd.StartLogin()
		default:
			cmd.GetSecuredItem(os.Args[1])
		}
	}

}
