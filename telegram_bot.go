package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	telegramBot *tgbotapi.BotAPI
)

func InitTelegramBot(telegramBotKey string) error {

	var err error
	telegramBot, err = tgbotapi.NewBotAPI(telegramBotKey)
	if err != nil {
		return err
	}

	return nil
}
