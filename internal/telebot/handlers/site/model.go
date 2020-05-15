package site

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/Syfaro/telegram-bot-api"
)

type SiteParamsModel struct {
	UserID         uint
	Name           string
	URL            string
	RequestTimeout int64
	ResponseStatus int64
	Description    string
}

func NewCreateSiteParamsModel(updates tgbotapi.Update) (SiteParamsModel, error) {
	var newSiteParamsModel = SiteParamsModel{}
	sSplit := strings.Split(updates.Message.CommandArguments(), " ")
	valid, err := isValidUrl(sSplit[0])
	if err != nil {
		log.Printf("Not an <URL> isValidUrl: %v", err)
		return newSiteParamsModel, fmt.Errorf("Not an <URL> entered at the first state: %v", err)
	}
	if valid {
		if len(sSplit) <= 1 {
			newSiteParamsModel = SiteParamsModel{
				UserID:         uint(updates.Message.From.ID),
				URL:            sSplit[0],
				Name:           "Default name",
				ResponseStatus: 200, //Default
				RequestTimeout: 120, //Default
			}
		} else if len(sSplit) <= 2 {
			newSiteParamsModel = SiteParamsModel{
				UserID:         uint(updates.Message.From.ID),
				URL:            sSplit[0],
				Name:           sSplit[1],
				ResponseStatus: 200, //Default
				RequestTimeout: 120, //Default
			}
		} else {
			return newSiteParamsModel, fmt.Errorf("/site_<add> <URL> <Name>. <Name> will set as 'Default name' if user will not name it by himself.")
		}
	}
	return newSiteParamsModel, nil
}

func CreateDeletingSiteParamsModel(updates tgbotapi.Update) (SiteParamsModel, error) {
	var newSiteParamsModel = SiteParamsModel{}
	sSplit := strings.Split(updates.Message.CommandArguments(), " ")
	valid, _ := isValidUrl(sSplit[0])
	if valid {
		if len(sSplit) <= 1 {
			newSiteParamsModel = SiteParamsModel{
				UserID: uint(updates.Message.From.ID),
				URL:    sSplit[0],
			}
		} else {
			return newSiteParamsModel, fmt.Errorf("/site_delete <URL> or <Name>.")
		}
	} else {
		if len(sSplit) <= 1 {
			newSiteParamsModel = SiteParamsModel{
				UserID: uint(updates.Message.From.ID),
				Name:   sSplit[0],
			}
		} else {
			return newSiteParamsModel, fmt.Errorf("/site_delete <URL> or <Name>.")
		}
	}
	return newSiteParamsModel, nil
}

//CreateRuleParamsModel - checking incoming/creating params for /rule_add
func CreateRuleParamsModel(updates tgbotapi.Update) (SiteParamsModel, error) {
	var newRuleParamsModel = SiteParamsModel{}
	sSplit := strings.Split(updates.Message.CommandArguments(), " ")
	valid, err := isValidUrl(sSplit[0])
	if err != nil {
		log.Printf("CreateRuleParamsModel err := isValidURL(sSplit[0]): %v", err)
		return newRuleParamsModel, err
	}
	if valid {
		if len(sSplit) < 3 {
			log.Printf("/rule_add - количество аргументов меньше трех: %v", len(sSplit))
			return newRuleParamsModel, fmt.Errorf("/rule_add <URL> <Response status> <Request timeout>.\nКоличесвтво аргументов меньше 3.")
		} else {
			newRuleParamsModel = SiteParamsModel{
				UserID:         uint(updates.Message.From.ID),
				URL:            sSplit[0],
				ResponseStatus: NewInt(sSplit[1]),
				RequestTimeout: NewInt(sSplit[2]),
			}
		}
	} else {
		if len(sSplit) < 3 {
			log.Printf("/rule_delete - количество аргументов меньше трех: %v", len(sSplit))
			return newRuleParamsModel, fmt.Errorf("/rule_add <Name> <Response status> <Request timeout>.\nКоличесвтво аргументов меньше 3.")
		} else {
			newRuleParamsModel = SiteParamsModel{
				UserID:         uint(updates.Message.From.ID),
				Name:           sSplit[0],
				ResponseStatus: NewInt(sSplit[1]),
				RequestTimeout: NewInt(sSplit[2]),
			}
		}
	}
	return newRuleParamsModel, nil
}

//CreateRuleDeleteParamsModel - checking incoming/creating  params for  /rule_delete
func CreateRuleDeleteParamsModel(updates tgbotapi.Update) (SiteParamsModel, error) {
	var newRuleParamsModel = SiteParamsModel{}
	sSplit := strings.Split(updates.Message.CommandArguments(), " ")
	valid, err := isValidUrl(sSplit[0])
	if err != nil {
		log.Printf("CreateRuleParamsModel err := isValidURL(sSplit[0]): %v", err)
	}
	if valid {
		if len(sSplit) > 1 {
			log.Printf("/rule_delete - количество аргументов меньше трех: %v", len(sSplit))
			return newRuleParamsModel, fmt.Errorf("/rule_delete <URL>.\nКоличесвтво аргументов меньше 3.")
		} else {
			newRuleParamsModel = SiteParamsModel{
				UserID:         uint(updates.Message.From.ID),
				URL:            sSplit[0],
				ResponseStatus: NewInt(sSplit[1]),
				RequestTimeout: NewInt(sSplit[2]),
			}
		}
	} else {
		if len(sSplit) < 3 {
			log.Printf("/rule_delete - количество аргументов меньше трех: %v", len(sSplit))
			return newRuleParamsModel, fmt.Errorf("/rule_delete <Name>.\nКоличесвтво аргументов меньше 3.")
		} else {
			newRuleParamsModel = SiteParamsModel{
				UserID:         uint(updates.Message.From.ID),
				Name:           sSplit[0],
				ResponseStatus: NewInt(sSplit[1]),
				RequestTimeout: NewInt(sSplit[2]),
			}
		}
	}
	return newRuleParamsModel, nil
}

//isValidUrl tests a string to determine
//if it is a url or not.
func isValidUrl(s string) (bool, error) {
	_, err := url.ParseRequestURI(s)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func NewInt(s string) int64 {
	i, _ := strconv.Atoi(s)
	return int64(i)
}
