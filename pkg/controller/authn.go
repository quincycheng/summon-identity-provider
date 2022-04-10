package controller

import (
	"errors"
	"fmt"

	"github.com/quincycheng/summon-identity-provider/pkg/api"
	"github.com/tidwall/gjson"
)

func GetPodFqnd(tenantId string, username string) (string, error) {
	fqdnFromTenantId := fmt.Sprintf("%s%s", tenantId, api.IdentityRootUrlSuffix)
	result, _, err := api.StartAuthnWithFqdn(fqdnFromTenantId, username, "")
	if err != nil {
		return "", err
	}

	//Check if PodFqdn is returned
	PodFqdnFromResponse := gjson.Get(result, "Result.PodFqdn").String()
	if PodFqdnFromResponse != "" {
		return PodFqdnFromResponse, nil
	} else {
		return fqdnFromTenantId, nil
	}
}

type MechType struct {
	RawJson              string
	AnswerType           string
	Name                 string
	PromptMechChosen     string
	PromptSelectMech     string
	MechanismId          string
	PartialAddress       string
	MaskedEmailAddress   string
	PartialDeviceAddress string
	ThirdPartyMfaRequest string
	PartialPhoneNumber   string
}

func GetMechanismTypeByValue(rawMechanisms gjson.Result, challengeValue string) (MechType, error) {
	// initialize result
	theType := MechType{
		RawJson:              "",
		AnswerType:           "",
		Name:                 "",
		PromptMechChosen:     "",
		PromptSelectMech:     "",
		MechanismId:          "",
		PartialAddress:       "", // Email
		MaskedEmailAddress:   "", // Email
		PartialDeviceAddress: "", // SMS
		ThirdPartyMfaRequest: "", // Duo
		PartialPhoneNumber:   "", // PF
	}

	for _, mech := range rawMechanisms.Array() {
		if mech.Get("PromptSelectMech").String() == challengeValue {
			theType = MechType{
				RawJson:              mech.String(),
				AnswerType:           mech.Get("AnswerType").String(),
				Name:                 mech.Get("Name").String(),
				PromptMechChosen:     mech.Get("PromptMechChosen").String(),
				PromptSelectMech:     mech.Get("PromptSelectMech").String(),
				MechanismId:          mech.Get("MechanismId").String(),
				PartialAddress:       mech.Get("PartialAddress").String(),       // Email
				MaskedEmailAddress:   mech.Get("MaskedEmailAddress").String(),   // Email
				PartialDeviceAddress: mech.Get("PartialDeviceAddress").String(), // SMS
				ThirdPartyMfaRequest: mech.Get("ThirdPartyMfaRequest").String(), // Duo
				PartialPhoneNumber:   mech.Get("PartialPhoneNumber").String(),   // PF
			}
		}
	}

	if theType.MechanismId == "" {
		return theType, errors.New("mechanism not found")
	}
	return theType, nil
}
