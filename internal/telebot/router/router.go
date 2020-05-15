package router

import (
	"context"
	"log"
	"time"

	"telebot/telebot/CA/internal/telebot/handlers/site"
	"telebot/telebot/CA/internal/telebot/handlers/user"

	"github.com/Syfaro/telegram-bot-api"
)

type Router struct {
	Bot         *tgbotapi.BotAPI
	UserHandler *user.Handler
	SiteHandler *site.Handler
}

func (r *Router) Init(token string) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Printf("Что-то пошло не так в (r *Router) Init: %v", err)
		return err
	}
	r.Bot = bot
	return nil
}

func (r *Router) Run() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := r.Bot.GetUpdatesChan(u)
	if err != nil {
		log.Printf("Check GetUpdatesChan Run for %v", err)
		return err
	}
	for update := range updates {
		ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
		if update.Message == nil {
			continue
		}
		switch update.Message.Command() {
		case "start":
			err := r.UserHandler.Start(ctx, update, r.Bot)
			if err != nil {
				log.Printf("Check UserHandler.Start /start for %v", err)
				return err
			}
		case "site_add":
			err := r.SiteHandler.SaveByURLorNameChooser(ctx, update, r.Bot)
			if err != nil {
				log.Printf("Check SaveByURLorNameChooser /site_add for %v", err)
				return err
			}
		case "site_delete":
			err := r.SiteHandler.DeleteByURLorNameChooser(ctx, update, r.Bot)
			if err != nil {
				log.Printf("Check DeleteByURLorNameChooser /site_delete for %v", err)
				return err
			}
		case "site_list":
			err := r.SiteHandler.FindIdUrlNameDescByUserID(ctx, update, r.Bot)
			if err != nil {
				log.Printf("Check FindIdUrlNameDescByUserID /site_list for %v", err)
				return err
			}
		case "site_update":
			err := r.SiteHandler.UpdateSiteRuleChooser(ctx, update, r.Bot)
			if err != nil {
				log.Printf("Check FindIdUrlNameDescByUserID /site_update for %v", err)
				return err
			}
		case "rule_delete":
			err := r.SiteHandler.UpdateSiteRuleChooser(ctx, update, r.Bot)
			if err != nil {
				log.Printf("Check /rule_delete's UpdateSiteRuleChooser for %v", err)
				return err
			}
		case "rule_list":
			err := r.SiteHandler.FindRuleByID(ctx, update, r.Bot)
			if err != nil {
				log.Printf("Check /rule_list's CallRuleAndSend for %v", err)
				return err
			}
		default:
			log.Printf("Switch default")
		}
	}
	return nil
}
