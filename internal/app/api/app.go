package api

import (
	"fmt"
	"log"

	"telebot/telebot/CA/internal/config"

	user_repository "telebot/telebot/CA/internal/database/user/pgsql"
	user_service "telebot/telebot/CA/internal/service/user"
	user_handler "telebot/telebot/CA/internal/telebot/handlers/user"

	site_repository "telebot/telebot/CA/internal/database/site/pgsql"
	site_service "telebot/telebot/CA/internal/service/site"
	site_handler "telebot/telebot/CA/internal/telebot/handlers/site"

	"telebot/telebot/CA/internal/telebot/router"
	"telebot/telebot/CA/pkg/infrastructure/pgsql"
)

type App struct {
	Router *router.Router
}

func Run() error {
	var app = new(App)
	err := app.Init("/home/stanislau/go/src/telebot/telebot/CA/config.json")
	if err != nil {
		return err
	}
	return app.Run()
}

func (a *App) Init(filePath string) error {
	cfg, err := config.New(filePath)
	if err != nil {
		log.Printf("Что-то не так в (a *App) Init: %v", err)
		return err
	}
	db, err := pgsql.New(cfg.DB)
	if err != nil {
		log.Printf("Что-то не так в (a *App) Init: %v", err)
		return err
	}

	//User connections
	userRepository := &user_repository.Repo{
		DB: db,
	}

	userService := &user_service.Service{
		UserRepository: userRepository,
	}

	userHandler := &user_handler.Handler{
		UserService: userService,
	}

	//Site connections
	siteRepository := &site_repository.Repo{
		DB: db,
	}

	siteService := &site_service.Service{
		SiteRepository: siteRepository,
	}

	siteHandler := &site_handler.Handler{
		SiteService: siteService,
	}

	a.Router = &router.Router{
		SiteHandler: siteHandler,
		UserHandler: userHandler,
	}

	if err := a.Router.Init(cfg.Token); err != nil {
		log.Printf("Что-то не так в (a *App) Init: %v", err)
		return err
	}

	return nil
}

func (a *App) Run() error {
	if err := a.Router.Run(); err != nil {
		return fmt.Errorf("(a *App) Run(): %v", err)
	}
	return a.Router.Run()
}
