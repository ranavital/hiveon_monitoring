package services

import (
	"hiveon_monitoring/config"
	"hiveon_monitoring/logger"
	"net/http"
	"net/url"
)

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
