package cmd

import (
	"fmt"
)

/**************************************************
	Commands: print the default message
*************************************************/
func PrintIntroMessage() {
	printHeader("\nSummon Provider for CyberArk Identity\n")
	fmt.Println("")
	fmt.Print(highlightStyle.Render("  summon"))
	fmt.Print(" is a command-line tool that reads a file in secrets.yml format and ")
	fmt.Print(highlightStyle.Render("injects secrets as environment variables"))
	fmt.Println(" into any process")
	fmt.Println("  Once the process exits, the secrets are gone")
	fmt.Println("  This is the provider for fetching passwords of Secured Items from CyberArk Identity")

	fmt.Println("")
	fmt.Println("To get started, you can configure the provider by executing:")
	fmt.Println(highlightStyle.Render("\tsummon-identity-provider --config"))
	fmt.Println("")
	fmt.Println("To login to CyberArk Identity, execute:")
	fmt.Println(highlightStyle.Render("\tsummon-identity-provider --login"))
	fmt.Println("")
	fmt.Println("Usage:")

	fmt.Printf("\t%s\t\t\t Print this message\n", highlightStyle.Render("summon-identity-provider"))
	fmt.Printf("\t%s\t Fetch the content of secured item from CyberArk Identity\n", highlightStyle.Render("summon-identity-provider <Secured Item Name>"))
	fmt.Printf("\t%s\t\t Start configuration wizard\n", highlightStyle.Render("summon-identity-provider --config"))
	fmt.Printf("\t%s\t\t Login to CyberArk Identity\n", highlightStyle.Render("summon-identity-provider --login"))

	fmt.Println("")
	fmt.Print("For more details: ")
	fmt.Println(highlightStyle.Render("https://github.com/quincycheng/summon-identity-provider\n"))

}
