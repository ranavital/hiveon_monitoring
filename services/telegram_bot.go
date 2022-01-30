package services

import (
	"hiveon_monitoring/config"
	"hiveon_monitoring/logger"
	"net/http"
	"net/url"
	// tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// var TgBot *tgbotapi.BotAPI

// func InitTgBot() error {
// 	bot, err := tgbotapi.NewBotAPI(config.AppConfig.TgToken)
// 	if err != nil {
// 		return err
// 	}

// 	TgBot = bot
// 	logger.Logging.Error("[InitTgBot]: telegram bot has been created successfully")
// 	return nil
// }

func SendTelegramAlert(alertMsg string) error {

	var telegramApi string = "https://api.telegram.org/bot" + config.AppConfig.TgToken + "/sendMessage"
	if _, err := http.PostForm(
		telegramApi,
		url.Values{
			"chat_id": {config.AppConfig.TgChatId},
			"text":    {alertMsg},
		},
	); err != nil {
		logger.Logging.Error("[SendTelegramAlert]: failed to send telegram bot message: %s", err)
	}

	logger.Logging.Info("[SendTelegramAlert]: successfully sent alert to telegram channel: %s", alertMsg)
	return nil
}
