package telegram_bot

import (
	"genieMap/app/refactor_xlsx"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"log"
)

//func GoBot() {
//	bot, err := tgbotapi.NewBotAPI("6094698233:AAFjy8J16Pd2Uv90W9VJ-UNqr9gZaBrWeqU")
//	if err != nil {
//		log.Panic(err)
//	}
//
//	bot.Debug = true
//
//	//log.Printf("Authorized on account %s", bot.Self.UserName)
//
//	u := tgbotapi.NewUpdate(0)
//	u.Timeout = 60
//
//	updates := bot.GetUpdatesChan(u)
//
//	//myPhotoPath := "/Users/anastasyaplotnikova/Desktop/2023-05-26 22.02.49.jpg"
//
//	refactor_xlsx.DoBookXlsx()
//	myXlsxFile := "refactor_available_units.xlsx"
//
//	for update := range updates {
//		if update.Message != nil {
//			//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
//
//			helloMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Сейчас сделаем обновленный файл")
//			bot.Send(helloMsg)
//
//			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "By the way: nice ass, Dude. Awesome balls, bro")
//			bot.Send(msg)
//
//			xlsxBytes, errphotoBytes := ioutil.ReadFile(myXlsxFile)
//			if errphotoBytes != nil {
//				log.Println("Failed to read photo:", err)
//				continue
//			}
//
//			xlsxConfig := tgbotapi.FileBytes{
//				Name:  myXlsxFile,
//				Bytes: xlsxBytes,
//			}
//
//			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Вот обнволение!")
//			bot.Send(msg)
//
//			docMsg := tgbotapi.NewDocument(update.Message.Chat.ID, xlsxConfig)
//			docMsg.ReplyToMessageID = update.Message.MessageID
//
//			bot.Send(docMsg)
//		}
//	}
//}

func StartBot() {
	bot, err := tgbotapi.NewBotAPI("6094698233:AAFjy8J16Pd2Uv90W9VJ-UNqr9gZaBrWeqU")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	myPhotoPath := "/Users/anastasyaplotnikova/Desktop/gachi.jpeg"

	for update := range updates {
		if update.Message != nil {

			////////////////////////////////////РАЗРАБОТКА////////////////////////////////////
			if update.Message.Text == "/start" {
				buttonMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Чем могу помочь?")

				// Создание кнопки
				buttonText := "Рефактор xlsx"
				button := tgbotapi.NewKeyboardButton(buttonText)

				// Создание клавиатуры и добавление кнопки
				keyboard := tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(button),
				)

				// Прикрепление клавиатуры к сообщению
				buttonMsg.ReplyMarkup = keyboard

				// Отправка сообщения с кнопкой
				bot.Send(buttonMsg)
			}
			////////////////////////////////////РАЗРАБОТКА////////////////////////////////////

			if update.Message.Document != nil {
				if isXLSXFile(update.Message.Document.FileName) {
					fileID := update.Message.Document.FileID
					fileUrl, GetFileDirectURLErr := bot.GetFileDirectURL(fileID)
					if GetFileDirectURLErr != nil {
						refactor_xlsx.LogError("Ошибка при получении файла: %v", GetFileDirectURLErr)
						errMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Упс... 🙊 Что-то пошло не так 😓\nОшибка при получении файла 📁 Уже бежим исправлять 🛠")
						bot.Send(errMsg)
						continue
					}

					helloMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Сейчас сделаю обновленный файл 🥰")
					bot.Send(helloMsg)

					myXlsxBuffer, refactorErr := refactor_xlsx.DoBookXlsx(fileUrl)
					if refactorErr != nil {
						refactor_xlsx.LogError("Ошибка при обработке файла: %v", refactorErr)
						errMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Упс... 🙊 Что-то пошло не так 😓\nОшибка при обработке файла 📁 Уже бежим исправлять 🛠")
						bot.Send(errMsg)
						continue
					}

					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вот обновление 🤩")
					bot.Send(msg)

					xlsxConfig := tgbotapi.FileBytes{
						Name:  "refactor_available_units.xlsx",
						Bytes: myXlsxBuffer.Bytes(),
					}

					docMsg := tgbotapi.NewDocument(update.Message.Chat.ID, xlsxConfig)
					docMsg.ReplyToMessageID = update.Message.MessageID
					bot.Send(docMsg)

					myPhotoByte, errToBytes := ioutil.ReadFile(myPhotoPath)
					if errToBytes == nil {
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, "By the way: nice ass, Dude 😏\nAwesome balls, bro 😎")
						bot.Send(msg)

						myPhoto := tgbotapi.FileBytes{
							Name:  myPhotoPath,
							Bytes: myPhotoByte,
						}

						PhotoMsg := tgbotapi.NewPhoto(update.Message.Chat.ID, myPhoto)
						PhotoMsg.Caption = "Ну что? Разомнемся?"
						bot.Send(PhotoMsg)
					}

					if errToBytes != nil {
						refactor_xlsx.LogError("Failed to read photo: %v", errToBytes)
						continue
					}

				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Присланный документ не имеет расширение .xlsx 🤕\nПожалуйста, пришлите файл в формате .xlsx")
					bot.Send(msg)
				}
			} else if update.Message.Text == "Рефактор xlsx" {
				// Если нажата кнопка "Рефактор xlsx", отправляем следующее сообщение
				replyMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пришлите файл")
				bot.Send(replyMsg)
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Я могу принимать файлы, но не могу с вами общаться 😔\nМогу вам помочь с рефактором вашего файла? 🥺")
				bot.Send(msg)
			}
		}
	}
}
