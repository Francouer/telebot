package user

import (
	"context"
	"fmt"
	"log"

	"telebot/telebot/CA/internal/domain/service"

	"github.com/Syfaro/telegram-bot-api"
)

type Handler struct {
	UserService service.UserService
}

func (h Handler) Start(ctx context.Context, updates tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	var message string
	startID := service.CreateUserParams{
		UserID: uint(updates.Message.From.ID),
	}

	if _, err := h.UserService.Create(ctx, startID); err != nil {
		log.Printf("%v", err)
		message = fmt.Sprintf("%v", err)
		return msgSenderToUser(updates, message, bot)

	} else {
		message = "User added."
		log.Printf("User added: userID - %v ", updates.Message.From.ID)
		return msgSenderToUser(updates, message, bot)
	}
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
