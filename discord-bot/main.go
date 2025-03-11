package main

import (
	"log"
	"vivalchemy/discord-bot/bot"
	"vivalchemy/discord-bot/config"
)

func main() {
	err := config.ReadFile("./.apiConfig.json")
	if err != nil {
		log.Println(err)
	}

	bot.Start()
	<-make(chan struct{})
}
