package cmd

import (
	"fmt"
	"os"

	"github.com/quincycheng/summon-identity-provider/pkg/controller"
)

func GetSecuredItem(securedItemName string) {
	value, err := controller.GetSecuredItem(securedItemName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(value)
}
