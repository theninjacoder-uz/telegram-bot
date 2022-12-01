package server

import (
	"log"
	"tgbot/configs"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	bot *tgbotapi.BotAPI
)

func InitTelegram() *tgbotapi.BotAPI {
	var err error

	bot, err = tgbotapi.NewBotAPI(configs.Config().TelgramBotToken)
	if err != nil {
		log.Panic(err)
	}

	// this perhaps should be conditional on GetWebhookInfo()
	// only set webhook if it is not set properly
	url := configs.Config().ServerHost + bot.Token
	_, err = tgbotapi.NewWebhook(url)
	if err != nil {
		log.Println(err)
	}

	return bot
}
