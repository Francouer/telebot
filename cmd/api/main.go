package main

import (
	"log"

	"telebot/telebot/CA/internal/app/api"
)

func main() {
	if err := api.Run(); err != nil {
		log.Printf("%v", err)
	}
}
