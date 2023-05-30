package telegram_bot

import (
	"genieMap/app/refactor_xlsx/Condor"
	"genieMap/app/refactor_xlsx/al_dar"
	"genieMap/app/refactor_xlsx/luma_22"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func getServiceAlDar(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	if message.Document != nil {
		fileID := message.Document.FileID
		fileURL, err := bot.GetFileDirectURL(fileID)
		if err != nil {
			log.Printf("Ошибка при получении файла: %v", err)
			return
		}

		sendProcessingMessage(bot, message.Chat.ID)

		xlsxBuffer, err := al_dar.DoBookCSV(fileURL)
		if err != nil {
			log.Printf("Ошибка при обработке файла: %v", err)
			return
		}

		sendUpdateMessage(bot, message.Chat.ID)
		sendCSVFile(bot, message.Chat.ID, xlsxBuffer)
	} else {
		errMsg(bot, message.Chat.ID)
	}
}

func getServiceLuna22(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	if message.Document != nil {
		fileID := message.Document.FileID
		fileURL, err := bot.GetFileDirectURL(fileID)
		if err != nil {
			log.Printf("Ошибка при получении файла: %v", err)
			return
		}

		sendProcessingMessage(bot, message.Chat.ID)

		xlsxBuffer, err := luma_22.DoBookCSV(fileURL)
		if err != nil {
			log.Printf("Ошибка при обработке файла: %v", err)
			return
		}

		sendUpdateMessage(bot, message.Chat.ID)
		sendCSVFile(bot, message.Chat.ID, xlsxBuffer)
	} else {
		errMsg(bot, message.Chat.ID)
	}
}

func getServiceCondor(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	if message.Document != nil {
		fileID := message.Document.FileID
		fileURL, err := bot.GetFileDirectURL(fileID)
		if err != nil {
			log.Printf("Ошибка при получении файла: %v", err)
			return
		}

		sendProcessingMessage(bot, message.Chat.ID)

		xlsxBuffer, err := Condor.DoBookCSV(fileURL)
		if err != nil {
			log.Printf("Ошибка при обработке файла: %v", err)
			return
		}

		sendUpdateMessage(bot, message.Chat.ID)
		sendCSVFile(bot, message.Chat.ID, xlsxBuffer)
	} else {
		errMsg(bot, message.Chat.ID)
	}
}
