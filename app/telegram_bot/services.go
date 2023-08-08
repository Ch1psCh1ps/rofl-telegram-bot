package telegram_bot

import (
	"genieMap/app/refactor_xlsx/Condor"
	"genieMap/app/refactor_xlsx/MAg"
	"genieMap/app/refactor_xlsx/Town_x"
	"genieMap/app/refactor_xlsx/al_dar"
	"genieMap/app/refactor_xlsx/azizi"
	"genieMap/app/refactor_xlsx/binghatii"
	"genieMap/app/refactor_xlsx/condor8Cols"
	"genieMap/app/refactor_xlsx/deyaar"
	"genieMap/app/refactor_xlsx/ellingtonProperties"
	"genieMap/app/refactor_xlsx/emaar"
	"genieMap/app/refactor_xlsx/object1"
	"genieMap/app/refactor_xlsx/reportage_properties"
	"genieMap/app/refactor_xlsx/siadah"
	"genieMap/app/refactor_xlsx/sothobys"
	"genieMap/app/refactor_xlsx/swiss"
	"genieMap/app/refactor_xlsx/tiger"
	"genieMap/app/refactor_xlsx/vincitore"
	"genieMap/cmd"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

func getServiceAlDar(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	if message.Document != nil {
		fileID := message.Document.FileID
		fileURL, err := bot.GetFileDirectURL(fileID)
		fileName := message.Document.FileName
		fileNameArray := strings.Split(fileName, ".")
		fileName = fileNameArray[0]
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
		sendCSVFile(bot, message.Chat.ID, xlsxBuffer, fileName)
	} else {
		errMsg(bot, message.Chat.ID)
	}
}

func getServiceTownX(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	if message.Document != nil {
		fileID := message.Document.FileID
		fileURL, err := bot.GetFileDirectURL(fileID)
		fileName := message.Document.FileName
		fileNameArray := strings.Split(fileName, ".")
		fileName = fileNameArray[0]
		if err != nil {
			log.Printf("Ошибка при получении файла: %v", err)
			return
		}

		sendProcessingMessage(bot, message.Chat.ID)

		xlsxBuffer, err := Town_x.DoBookCSV(fileURL)
		if err != nil {
			log.Printf("Ошибка при обработке файла: %v", err)
			return
		}

		sendUpdateMessage(bot, message.Chat.ID)
		sendCSVFile(bot, message.Chat.ID, xlsxBuffer, fileName)
	} else {
		errMsg(bot, message.Chat.ID)
	}
}

func getServiceCondor(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	if message.Document != nil {
		fileID := message.Document.FileID
		fileURL, err := bot.GetFileDirectURL(fileID)
		fileName := message.Document.FileName
		fileNameArray := strings.Split(fileName, ".")
		fileName = fileNameArray[0]
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
		sendCSVFile(bot, message.Chat.ID, xlsxBuffer, fileName)
	} else {
		errMsg(bot, message.Chat.ID)
	}
}

func getServiceCondor8Cols(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	if message.Document != nil {
		fileID := message.Document.FileID
		fileURL, err := bot.GetFileDirectURL(fileID)
		fileName := message.Document.FileName
		fileNameArray := strings.Split(fileName, ".")
		fileName = fileNameArray[0]
		if err != nil {
			log.Printf("Ошибка при получении файла: %v", err)
			return
		}

		sendProcessingMessage(bot, message.Chat.ID)

		xlsxBuffer, err := condor8Cols.DoBookCSV(fileURL)
		if err != nil {
			log.Printf("Ошибка при обработке файла: %v", err)
			return
		}

		sendUpdateMessage(bot, message.Chat.ID)
		sendCSVFile(bot, message.Chat.ID, xlsxBuffer, fileName)
	} else {
		errMsg(bot, message.Chat.ID)
	}
}

func getServiceSiadah(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	if message.Document != nil {
		fileID := message.Document.FileID
		fileURL, err := bot.GetFileDirectURL(fileID)
		fileName := message.Document.FileName
		fileNameArray := strings.Split(fileName, ".")
		fileName = fileNameArray[0]
		if err != nil {
			log.Printf("Ошибка при получении файла: %v", err)
			return
		}

		sendProcessingMessage(bot, message.Chat.ID)

		xlsxBuffer, err := siadah.DoBookCSV(fileURL)
		if err != nil {
			log.Printf("Ошибка при обработке файла: %v", err)
			return
		}

		sendUpdateMessage(bot, message.Chat.ID)
		sendCSVFile(bot, message.Chat.ID, xlsxBuffer, fileName)
	} else {
		errMsg(bot, message.Chat.ID)
	}
}

func getServiceBinghatii(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	if message.Document != nil {
		fileID := message.Document.FileID
		fileURL, err := bot.GetFileDirectURL(fileID)
		fileName := message.Document.FileName
		fileNameArray := strings.Split(fileName, ".")
		fileName = fileNameArray[0]
		if err != nil {
			log.Printf("Ошибка при получении файла: %v", err)
			return
		}

		sendProcessingMessage(bot, message.Chat.ID)

		fileContent, downloadFileErr := cmd.DownloadFile(fileURL)
		if downloadFileErr != nil {
			log.Printf("Ошибка при загрузке файла: %v", downloadFileErr)
		}

		data := cmd.GetData(fileContent)

		sheetList := data.GetSheetList()
		for _, sheetName := range sheetList {
			xlsxBuffer, err4 := binghatii.DoBookCSV(fileURL, sheetName)
			if err4 != nil {
				log.Printf("Ошибка при обработке файла: %v", err)
				return
			}

			//sendUpdateMessage(bot, message.Chat.ID)
			sendCSVFile(bot, message.Chat.ID, xlsxBuffer, sheetName)
		}
		sendUpdateMessage(bot, message.Chat.ID)
		sendAttentionMessage(bot, message.Chat.ID)
	} else {
		errMsg(bot, message.Chat.ID)
	}
}

//func getServiceBinghatii(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
//	if message.Document != nil {
//		fileID := message.Document.FileID
//		fileURL, err := bot.GetFileDirectURL(fileID)
//		fileName := message.Document.FileName
//		fileNameArray := strings.Split(fileName, ".")
//		fileName = fileNameArray[0]
//		if err != nil {
//			log.Printf("Ошибка при получении файла: %v", err)
//			return
//		}
//
//		sendProcessingMessage(bot, message.Chat.ID)
//
//		xlsxBuffer, err := binghatii.DoBookCSV(fileURL)
//		if err != nil {
//			log.Printf("Ошибка при обработке файла: %v", err)
//			return
//		}
//
//		sendUpdateMessage(bot, message.Chat.ID)
//		sendCSVFile(bot, message.Chat.ID, xlsxBuffer, fileName)
//	} else {
//		errMsg(bot, message.Chat.ID)
//	}
//}

func getServiceDeyaar(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	if message.Document != nil {
		fileID := message.Document.FileID
		fileURL, err := bot.GetFileDirectURL(fileID)
		fileName := message.Document.FileName
		fileNameArray := strings.Split(fileName, ".")
		fileName = fileNameArray[0]
		if err != nil {
			log.Printf("Ошибка при получении файла: %v", err)
			return
		}

		sendProcessingMessage(bot, message.Chat.ID)

		xlsxBuffer, err := deyaar.DoBookCSV(fileURL)
		if err != nil {
			log.Printf("Ошибка при обработке файла: %v", err)
			return
		}

		sendUpdateMessage(bot, message.Chat.ID)
		sendCSVFile(bot, message.Chat.ID, xlsxBuffer, fileName)
	} else {
		errMsg(bot, message.Chat.ID)
	}
}

func getServiceEmaar(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	if message.Document != nil {
		fileID := message.Document.FileID
		fileURL, err := bot.GetFileDirectURL(fileID)
		fileName := message.Document.FileName
		fileNameArray := strings.Split(fileName, ".")
		fileName = fileNameArray[0]
		if err != nil {
			log.Printf("Ошибка при получении файла: %v", err)
			return
		}

		sendProcessingMessage(bot, message.Chat.ID)

		xlsxBuffer, err := emaar.DoBookCSV(fileURL)
		if err != nil {
			log.Printf("Ошибка при обработке файла: %v", err)
			return
		}

		sendUpdateMessage(bot, message.Chat.ID)
		sendCSVFile(bot, message.Chat.ID, xlsxBuffer, fileName)
	} else {
		errMsg(bot, message.Chat.ID)
	}
}

func getServiceEllingtonProperties(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	if message.Document != nil {
		fileID := message.Document.FileID
		fileURL, err := bot.GetFileDirectURL(fileID)
		fileName := message.Document.FileName
		fileNameArray := strings.Split(fileName, ".")
		fileName = fileNameArray[0]
		if err != nil {
			log.Printf("Ошибка при получении файла: %v", err)
			return
		}

		sendProcessingMessage(bot, message.Chat.ID)

		fileContent, downloadFileErr := cmd.DownloadFile(fileURL)
		if downloadFileErr != nil {
			log.Printf("Ошибка при загрузке файла: %v", downloadFileErr)
		}

		data := cmd.GetData(fileContent)

		sheetList := data.GetSheetList()
		for _, sheetName := range sheetList {
			xlsxBuffer, err4 := ellingtonProperties.DoBookCSV(fileURL, sheetName)
			if err4 != nil {
				log.Printf("Ошибка при обработке файла: %v", err)
				return
			}

			//sendUpdateMessage(bot, message.Chat.ID)
			if sheetName != "hiddenSheet" {
				sendCSVFile(bot, message.Chat.ID, xlsxBuffer, sheetName)
			}
		}
		sendUpdateMessage(bot, message.Chat.ID)
		sendAttentionMessage(bot, message.Chat.ID)
	} else {
		errMsg(bot, message.Chat.ID)
	}
}

func getServiceAzizi(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	if message.Document != nil {
		fileID := message.Document.FileID
		fileURL, err := bot.GetFileDirectURL(fileID)
		fileName := message.Document.FileName
		fileNameArray := strings.Split(fileName, ".")
		fileName = fileNameArray[0]
		if err != nil {
			log.Printf("Ошибка при получении файла: %v", err)
			return
		}

		sendProcessingMessage(bot, message.Chat.ID)

		fileContent, downloadFileErr := cmd.DownloadFile(fileURL)
		if downloadFileErr != nil {
			log.Printf("Ошибка при загрузке файла: %v", downloadFileErr)
		}

		data := cmd.GetData(fileContent)

		sheetList := data.GetSheetList()
		for _, sheetName := range sheetList {
			xlsxBuffer, err4 := azizi.DoBookCSV(fileURL, sheetName)
			if err4 != nil {
				log.Printf("Ошибка при обработке файла: %v", err)
				return
			}

			//sendUpdateMessage(bot, message.Chat.ID)
			sendCSVFile(bot, message.Chat.ID, xlsxBuffer, sheetName)
		}
		sendUpdateMessage(bot, message.Chat.ID)
		sendAttentionMessage(bot, message.Chat.ID)
	} else {
		errMsg(bot, message.Chat.ID)
	}
}

func getServiceReportageProperties(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	if message.Document != nil {
		fileID := message.Document.FileID
		fileURL, err := bot.GetFileDirectURL(fileID)
		fileName := message.Document.FileName
		fileNameArray := strings.Split(fileName, ".")
		fileName = fileNameArray[0]
		if err != nil {
			log.Printf("Ошибка при получении файла: %v", err)
			return
		}

		sendProcessingMessage(bot, message.Chat.ID)

		xlsxBuffer, err := reportage_properties.DoBookCSV(fileURL)
		if err != nil {
			log.Printf("Ошибка при обработке файла: %v", err)
			return
		}

		sendUpdateMessage(bot, message.Chat.ID)
		sendCSVFile(bot, message.Chat.ID, xlsxBuffer, fileName)
	} else {
		errMsg(bot, message.Chat.ID)
	}
}

func getServiceTiger(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	if message.Document != nil {
		fileID := message.Document.FileID
		fileURL, err := bot.GetFileDirectURL(fileID)
		fileName := message.Document.FileName
		fileNameArray := strings.Split(fileName, ".")
		fileName = fileNameArray[0]
		if err != nil {
			log.Printf("Ошибка при получении файла: %v", err)
			return
		}

		sendProcessingMessage(bot, message.Chat.ID)

		fileContent, downloadFileErr := cmd.DownloadFile(fileURL)
		if downloadFileErr != nil {
			log.Printf("Ошибка при загрузке файла: %v", downloadFileErr)
		}

		data := cmd.GetData(fileContent)

		sheetList := data.GetSheetList()
		for _, sheetName := range sheetList {
			xlsxBuffer, err4 := tiger.DoBookCSV(fileURL, sheetName)
			if err4 != nil {
				log.Printf("Ошибка при обработке файла: %v", err)
				return
			}

			//sendUpdateMessage(bot, message.Chat.ID)
			sendCSVFile(bot, message.Chat.ID, xlsxBuffer, sheetName)
		}
		sendUpdateMessage(bot, message.Chat.ID)
		sendAttentionMessage(bot, message.Chat.ID)
	} else {
		errMsg(bot, message.Chat.ID)
	}
}

func getServiceSothobys(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	if message.Document != nil {
		fileID := message.Document.FileID
		fileURL, err := bot.GetFileDirectURL(fileID)
		fileName := message.Document.FileName
		fileNameArray := strings.Split(fileName, ".")
		fileName = fileNameArray[0]
		if err != nil {
			log.Printf("Ошибка при получении файла: %v", err)
			return
		}

		sendProcessingMessage(bot, message.Chat.ID)

		fileContent, downloadFileErr := cmd.DownloadFile(fileURL)
		if downloadFileErr != nil {
			log.Printf("Ошибка при загрузке файла: %v", downloadFileErr)
		}

		data := cmd.GetData(fileContent)

		sheetList := data.GetSheetList()
		for _, sheetName := range sheetList {
			xlsxBuffer, err4 := sothobys.DoBookCSV(fileURL, sheetName)
			if err4 != nil {
				log.Printf("Ошибка при обработке файла: %v", err)
				return
			}

			sendCSVFile(bot, message.Chat.ID, xlsxBuffer, sheetName)
		}
		sendUpdateMessage(bot, message.Chat.ID)
		sendAttentionMessage(bot, message.Chat.ID)
	} else {
		errMsg(bot, message.Chat.ID)
	}
}

func getServiceObject1(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	if message.Document != nil {
		fileID := message.Document.FileID
		fileURL, err := bot.GetFileDirectURL(fileID)
		fileName := message.Document.FileName
		fileNameArray := strings.Split(fileName, ".")
		fileName = fileNameArray[0]
		if err != nil {
			log.Printf("Ошибка при получении файла: %v", err)
			return
		}

		sendProcessingMessage(bot, message.Chat.ID)

		fileContent, downloadFileErr := cmd.DownloadFile(fileURL)
		if downloadFileErr != nil {
			log.Printf("Ошибка при загрузке файла: %v", downloadFileErr)
		}

		data := cmd.GetData(fileContent)

		sheetList := data.GetSheetList()
		for _, sheetName := range sheetList {
			xlsxBuffer, err4 := object1.DoBookCSV(fileURL, sheetName)
			if err4 != nil {
				log.Printf("Ошибка при обработке файла: %v", err)
				return
			}

			sendCSVFile(bot, message.Chat.ID, xlsxBuffer, sheetName)
		}
		sendUpdateMessage(bot, message.Chat.ID)
		sendAttentionMessage(bot, message.Chat.ID)
	} else {
		errMsg(bot, message.Chat.ID)
	}
}

func getServiceSwiss(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	if message.Document != nil {
		fileID := message.Document.FileID
		fileURL, err := bot.GetFileDirectURL(fileID)
		fileName := message.Document.FileName
		fileNameArray := strings.Split(fileName, ".")
		fileName = fileNameArray[0]
		if err != nil {
			log.Printf("Ошибка при получении файла: %v", err)
			return
		}

		sendProcessingMessage(bot, message.Chat.ID)

		fileContent, downloadFileErr := cmd.DownloadFile(fileURL)
		if downloadFileErr != nil {
			log.Printf("Ошибка при загрузке файла: %v", downloadFileErr)
		}

		data := cmd.GetData(fileContent)

		sheetList := data.GetSheetList()
		for _, sheetName := range sheetList {
			xlsxBuffer, err4 := swiss.DoBookCSV(fileURL, sheetName)
			if err4 != nil {
				log.Printf("Ошибка при обработке файла: %v", err)
				return
			}

			//sendUpdateMessage(bot, message.Chat.ID)
			sendCSVFile(bot, message.Chat.ID, xlsxBuffer, sheetName)
		}
		sendUpdateMessage(bot, message.Chat.ID)
		sendAttentionMessage(bot, message.Chat.ID)
	} else {
		errMsg(bot, message.Chat.ID)
	}
}

func getServiceVincitore(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	if message.Document != nil {
		fileID := message.Document.FileID
		fileURL, err := bot.GetFileDirectURL(fileID)
		fileName := message.Document.FileName
		fileNameArray := strings.Split(fileName, ".")
		fileName = fileNameArray[0]
		if err != nil {
			log.Printf("Ошибка при получении файла: %v", err)
			return
		}

		sendProcessingMessage(bot, message.Chat.ID)

		fileContent, downloadFileErr := cmd.DownloadFile(fileURL)
		if downloadFileErr != nil {
			log.Printf("Ошибка при загрузке файла: %v", downloadFileErr)
		}

		data := cmd.GetData(fileContent)

		sheetList := data.GetSheetList()
		for _, sheetName := range sheetList {
			xlsxBuffer, err4 := vincitore.DoBookCSV(fileURL, sheetName)
			if err4 != nil {
				log.Printf("Ошибка при обработке файла: %v", err)
				return
			}

			//sendUpdateMessage(bot, message.Chat.ID)
			sendCSVFile(bot, message.Chat.ID, xlsxBuffer, sheetName)
		}
		sendUpdateMessage(bot, message.Chat.ID)
		sendAttentionMessage(bot, message.Chat.ID)
	} else {
		errMsg(bot, message.Chat.ID)
	}
}

//func getServiceMag(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
//	if message.Document != nil {
//		fileID := message.Document.FileID
//		fileURL, err := bot.GetFileDirectURL(fileID)
//		fileName := message.Document.FileName
//		fileNameArray := strings.Split(fileName, ".")
//		fileName = fileNameArray[0]
//		if err != nil {
//			log.Printf("Ошибка при получении файла: %v", err)
//			return
//		}
//
//		sendProcessingMessage(bot, message.Chat.ID)
//
//		fileContent, downloadFileErr := cmd.DownloadFile(fileURL)
//		if downloadFileErr != nil {
//			log.Printf("Ошибка при загрузке файла: %v", downloadFileErr)
//		}
//
//		data := cmd.GetData(fileContent)
//
//		sheetList := data.GetSheetList()
//		for _, sheetName := range sheetList {
//			xlsxBuffer, err4 := MAg.DoBookCSV(fileURL, sheetName)
//			if err4 != nil {
//				log.Printf("Ошибка при обработке файла: %v", err)
//				return
//			}
//
//			//sendUpdateMessage(bot, message.Chat.ID)
//			sendCSVFile(bot, message.Chat.ID, xlsxBuffer, sheetName)
//		}
//		sendUpdateMessage(bot, message.Chat.ID)
//		sendAttentionMessage(bot, message.Chat.ID)
//	} else {
//		errMsg(bot, message.Chat.ID)
//	}
//}

func getServiceMag(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	if message.Document != nil {
		fileID := message.Document.FileID
		fileURL, err := bot.GetFileDirectURL(fileID)
		fileName := message.Document.FileName
		fileNameArray := strings.Split(fileName, ".")
		fileName = fileNameArray[0]
		if err != nil {
			log.Printf("Ошибка при получении файла: %v", err)
			return
		}

		sendProcessingMessage(bot, message.Chat.ID)

		xlsxBuffer, err := MAg.DoBookCSV(fileURL)
		if err != nil {
			log.Printf("Ошибка при обработке файла: %v", err)
			return
		}

		sendUpdateMessage(bot, message.Chat.ID)
		sendCSVFile(bot, message.Chat.ID, xlsxBuffer, fileName)
	} else {
		errMsg(bot, message.Chat.ID)
	}
}
