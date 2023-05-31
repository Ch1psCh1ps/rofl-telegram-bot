package telegram_bot

import (
	"bytes"
	"fmt"
	"genieMap/cmd"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strconv"
	"strings"
)

var serviceNum int

func getServices() []string {
	services := []string{
		"Al Dar",
		"Luna 22",
		"Condor",
		"Siadah",
		"Binghatii",
	}

	return services
}

func GetBot() {
	cmd.LoadEnv()
	apiToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		message := strings.ToLower(update.Message.Text)
		chatID := update.Message.Chat.ID

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				msg := tgbotapi.NewMessage(chatID, "Выберите сервисы")
				bot.Send(msg)

				sendServiceList(bot, chatID)
			}
		}

		if isServiceNumber(message) {
			serviceNumber, _ := strconv.Atoi(message)

			switch serviceNumber {
			case 1:
				msg := tgbotapi.NewMessage(chatID, "Вы выбрали сервис Al Dar")
				bot.Send(msg)
				msg = tgbotapi.NewMessage(chatID, "Пришлите файл в формате CSV")
				bot.Send(msg)

				serviceNum = 1
			case 2:
				msg := tgbotapi.NewMessage(chatID, "Вы выбрали сервис Luma22")
				bot.Send(msg)
				msg = tgbotapi.NewMessage(chatID, "Пришлите файл в формате XLSX")
				bot.Send(msg)

				serviceNum = 2
			case 3:
				msg := tgbotapi.NewMessage(chatID, "Вы выбрали сервис Condor")
				bot.Send(msg)
				msg = tgbotapi.NewMessage(chatID, "Пришлите файл в формате XLSX")
				bot.Send(msg)

				serviceNum = 3
			case 4:
				msg := tgbotapi.NewMessage(chatID, "Вы выбрали сервис Siadah")
				bot.Send(msg)
				msg = tgbotapi.NewMessage(chatID, "Пришлите файл в формате XLSX")
				bot.Send(msg)

				serviceNum = 4
			case 5:
				msg := tgbotapi.NewMessage(chatID, "Вы выбрали сервис Binghatii")
				bot.Send(msg)
				msg = tgbotapi.NewMessage(chatID, "Пришлите файл в формате XLSX")
				bot.Send(msg)

				serviceNum = 5
			default:
				msg := tgbotapi.NewMessage(chatID, "Некорректный выбор сервиса")
				bot.Send(msg)
			}
		}

		if update.Message.Document != nil {

			switch serviceNum {
			case 1:
				getServiceAlDar(update.Message, bot)
			case 2:
				getServiceLuma22(update.Message, bot)
			case 3:
				getServiceCondor(update.Message, bot)
			case 4:
				getServiceSiadah(update.Message, bot)
			case 5:
				getServiceBinghatii(update.Message, bot)
			default:
				msg := tgbotapi.NewMessage(chatID, "Сначала напишите номер сервиса")
				bot.Send(msg)
			}
		}
	}
}

func isServiceNumber(message string) bool {
	serviceNumber, err := strconv.Atoi(message)
	if err != nil {
		return false
	}

	correctNumber := serviceNumber >= 1 && serviceNumber <= len(getServices())

	return correctNumber
}

func sendProcessingMessage(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Сейчас сделаю обновленный файл 🥰")
	bot.Send(msg)
}

func sendUpdateMessage(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Вот обновление 🤩")
	bot.Send(msg)
}

func sendCSVFile(bot *tgbotapi.BotAPI, chatID int64, xlsxBuffer *bytes.Buffer, fileName string) {
	xlsxConfig := tgbotapi.FileBytes{
		Name:  fileName + ".csv",
		Bytes: xlsxBuffer.Bytes(),
	}

	docMsg := tgbotapi.NewDocument(chatID, xlsxConfig)
	bot.Send(docMsg)
}

func errMsg(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Требуется отправить файл в формате CSV")
	bot.Send(msg)
}

func sendServiceList(bot *tgbotapi.BotAPI, chatID int64) {
	services := getServices()

	var numberedServices strings.Builder
	for i, service := range services {
		numberedServices.WriteString(fmt.Sprintf("%d. %s\n", i+1, service))
	}

	servicesMsg := tgbotapi.NewMessage(chatID, numberedServices.String())
	bot.Send(servicesMsg)
}
