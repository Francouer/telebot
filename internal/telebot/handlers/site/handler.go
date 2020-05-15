package site

import (
	"context"
	"fmt"
	"log"

	"telebot/telebot/CA/internal/domain/service"

	"github.com/Syfaro/telegram-bot-api"
)

type Handler struct {
	SiteService service.SiteService
}

var message string

//SaveByURLorNameChooser - choose the function to
//for handle incoming save option from user.
//Choose between:
// addSiteToDatabaseByURL and addSiteToDatabaseByName.
func (h *Handler) SaveByURLorNameChooser(ctx context.Context, updates tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	incData, err := NewCreateSiteParamsModel(updates)
	if err != nil {
		log.Printf("Check SaveByURLorNameChooser NewCreateSiteParamsModel for %v", err)
		message = fmt.Sprintf("Something went wrong while we was searching for your site: %v", err)
		return msgSenderToUser(updates, message, bot)
	}
	s := service.CreateSiteParams{
		UserID: incData.UserID,
		Name:   incData.Name,
	}
	sp := service.SiteParams{
		URL:            &incData.URL,
		RequestTimeout: &incData.RequestTimeout,
		ResponseStatus: &incData.ResponseStatus,
		Description:    &incData.Description,
	}
	if *sp.URL != "" {
		err = addUserIdSiteToDatabaseByURL(ctx, h, s, updates, bot)
		return err
	} else {
		err = addUserIdSiteToDatabaseByName(ctx, h, s, updates, bot)
		return err
	}

}

//DeleteByURLorNameChooser - choose the function to
//for handle incoming save option from user.
//Choose between:
// deleteSiteFromDatabaseByURL and deleteSiteFromDatabaseByName.
func (h *Handler) DeleteByURLorNameChooser(ctx context.Context, updates tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	incData, err := CreateDeletingSiteParamsModel(updates)
	if err != nil {
		log.Printf("Check DeleteByURLorNameChooser NewCreateSiteParamsModel for %v", err)
		message = fmt.Sprintf("Something went wrong while we was searching for your site: %v", err)
		return msgSenderToUser(updates, message, bot)
	}

	s := service.CreateSiteParams{
		UserID: incData.UserID,
		Name:   incData.Name,
	}
	sp := service.SiteParams{
		URL: &incData.URL,
	}
	if *sp.URL != "" {
		err = deleteUserIdSiteFromDatabaseByURL(ctx, h, s, updates, bot)
		return err
	} else {
		err = deleteUserIdSiteFromDatabaseByName(ctx, h, s, updates, bot)
		return err
	}
}

//FindIdUrlNameDescByUserID - call to find by UserID
// all start registrations site params
func (h *Handler) FindIdUrlNameDescByUserID(ctx context.Context, updates tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	var iSlice []string
	siteParams := service.UpdateSiteParams{
		UserID: uint(updates.Message.From.ID),
	}

	site, err := h.SiteService.GetByUserID(ctx, siteParams.UserID)
	if err != nil {
		message = fmt.Sprintf("Can't find your list of sites: %v", err)
		log.Printf("Check FindIdUrlNameDescByUserID GetByIDAndUserID for %v", err)
		return msgSenderToUser(updates, message, bot)
	}
	for _, v := range site {
		info := fmt.Sprintf("%s    :    %s     :    %s", v.URL, v.Name, v.Description)
		iSlice = append(iSlice, info)
	}

	//TODO:telegram bot buttons for user

	message = `Your monitoring list:
		<URL> : <Name> : <Description>
         %s`

	finish := fmt.Sprintf(message, iSlice)

	return msgSenderToUser(updates, finish, bot)
}

//Rule changers part

func (h *Handler) UpdateSiteRuleChooser(ctx context.Context, updates tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	incData, err := NewCreateSiteParamsModel(updates)
	if err != nil {
		log.Printf("Check SaveByURLorNameChooser NewCreateSiteParamsModel for %v", err)
		message = fmt.Sprintf("Something went wrong while we was searching for your site: %v", err)
		return msgSenderToUser(updates, message, bot)
	}
	siteParams := service.UpdateSiteParams{
		UserID: incData.UserID,
		Name:   &incData.Name,
	}

	sp := service.SiteParams{
		URL:            &incData.URL,
		ResponseStatus: &incData.ResponseStatus,
		RequestTimeout: &incData.RequestTimeout,
		Description:    &incData.Description, //TODO: default value at start
	}
	if *sp.URL != "" {
		err = updateSiteRuleFindByUrlUserId(ctx, h, siteParams, updates, bot)
		return err
	} else {
		err = updateSiteRuleFindByNameUserId(ctx, h, siteParams, updates, bot)
		return err
	}
}

//CallRuleAndSend - translate incoming message from bot to the service and
func (h *Handler) FindRuleByID(ctx context.Context, updates tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	var message string
	var sMsg []string
	var info string
	siteParams := service.UpdateSiteParams{
		UserID: uint(updates.Message.From.ID),
	}

	rule, err := h.SiteService.GetByIDAndUserID(ctx, siteParams.ID, siteParams.UserID)
	if err != nil {
		message = fmt.Sprintf("Can't call for rule information %s", err)
		log.Printf("Check CallRuleAndSend's CallRuleAndSendResult for %s Url", err)
		return msgSenderToUser(updates, message, bot)
	}
	log.Printf("%v", rule)
	message = `List of rules for your site: 
            <Status> : <Timeout> : <Description>
               `
	err = msgSenderToUser(updates, message, bot)
	for i := 0; i == len(sMsg); i++ {
		info := fmt.Sprintf("\n%v, %v, %s\n", rule.ResponseStatus, rule.RequestTimeout, rule.Description)
		sMsg = append(sMsg, info)

	}
	for _, v := range sMsg {
		info = v
		_ = msgSenderToUser(updates, info, bot)
	}
	return err
}

//Functions add/delete/find site
//addSiteToDatabaseFindByURL - adding non- URL to the database's table "sites"
func addUserIdSiteToDatabaseByURL(ctx context.Context, h *Handler, siteParams service.CreateSiteParams, updates tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	_, err := h.SiteService.Create(ctx, siteParams)
	if err != nil {
		log.Printf("Check FindByURLorSave) addSiteToDatabaseFindByURL for %v", err)
		message = fmt.Sprintf("Can't find site with this URL in the DB: %v", err)
		return msgSenderToUser(updates, message, bot)
	}
	message = "Site added to the database."
	log.Printf("Site added to the database URL: %s", updates.Message.CommandArguments())
	return msgSenderToUser(updates, message, bot)
}

//addSiteToDatabaseFindByName - adding non- Name to the database's table "sites"
func addUserIdSiteToDatabaseByName(ctx context.Context, h *Handler, siteParams service.CreateSiteParams, updates tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	_, err := h.SiteService.Create(ctx, siteParams)
	if err != nil {
		log.Printf("Check addSiteToDatabaseByName FindByURLorSave for %v", err)
		message = fmt.Sprintf("Can't find site with this name in the DB: %v", err)
		return msgSenderToUser(updates, message, bot)
	}
	message = "Site added to the database"
	log.Printf("Site added to the database Name: %s", updates.Message.CommandArguments())
	return msgSenderToUser(updates, message, bot)
}

//deleteSiteFromDatabaseByName -
func deleteUserIdSiteFromDatabaseByName(ctx context.Context, h *Handler, csp service.CreateSiteParams, updates tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	byName, err := h.SiteService.GetIdByUserIdAndName(ctx, csp.Name, csp.UserID)
	if err != nil {
		return err
	}
	err = h.SiteService.DeleteByIDAndUserID(ctx, byName.ID, byName.UserID)
	if err != nil {
		log.Printf("Check deleteSiteFromDatabaseByName FindAndDeleteByName for %v Url", err)
		message = fmt.Sprintf("Can't delete site from DB by this Name: %v", err)
		return msgSenderToUser(updates, message, bot)
	}
	message = "Site deleted."
	log.Printf("Site deleted by Name: %s", updates.Message.CommandArguments())
	return msgSenderToUser(updates, message, bot)
}

//deleteSiteFromDatabaseByURL -
func deleteUserIdSiteFromDatabaseByURL(ctx context.Context, h *Handler, csp service.CreateSiteParams, updates tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	byName, err := h.SiteService.GetIdByUserIdAndUrl(ctx, csp.Name, csp.UserID)
	if err != nil {
		return err
	}
	err = h.SiteService.DeleteByIDAndUserID(ctx, byName.ID, byName.UserID)
	if err != nil {
		log.Printf("Check deleteSiteFromDatabaseByName FindAndDeleteByName for %v Url", err)
		message = fmt.Sprintf("Can't delete site from DB by this Name: %v", err)
		return msgSenderToUser(updates, message, bot)
	}
	message = "Site added to the database"
	log.Printf("Site added to the database Name: %s", updates.Message.CommandArguments())
	return msgSenderToUser(updates, message, bot)
}

//Functions - site's rule update/list
//updateSiteRuleFindByUrlUserId
func updateSiteRuleFindByUrlUserId(ctx context.Context, h *Handler, csp service.UpdateSiteParams, updates tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	err := h.SiteService.Update(ctx, csp)
	if err != nil {
		log.Printf("Check FindByURLorSave updateSiteRuleFindByUrlUserId for %v", err)
		message = fmt.Sprintf("Can't find site with this URL in the DB: %v", err)
		return msgSenderToUser(updates, message, bot)
	}
	message = "New site's rule added to the database."
	log.Printf("New site's rule added to the database by URL: %s", updates.Message.CommandArguments())
	return msgSenderToUser(updates, message, bot)
}

func updateSiteRuleFindByNameUserId(ctx context.Context, h *Handler, csp service.UpdateSiteParams, updates tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	err := h.SiteService.Update(ctx, csp)
	if err != nil {
		log.Printf("Check FindByURLorSave addSiteToDatabaseFindByURL for %v", err)
		message = fmt.Sprintf("Can't find site with this URL in the DB: %v", err)
		return msgSenderToUser(updates, message, bot)
	}
	message = "New site's rule added to the database."
	log.Printf("New site's rule added to the database by Name: %s", updates.Message.CommandArguments())
	return msgSenderToUser(updates, message, bot)
}

//msgSenderToUser - uses tgbotapi for sending message from bot to the user.
func msgSenderToUser(updates tgbotapi.Update, message string, bot *tgbotapi.BotAPI) error {
	msg := tgbotapi.NewMessage(updates.Message.Chat.ID, message)
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Check msgSenderToUser for %s", err)
		return err
	}
	return nil
}
