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

var waitingForFile map[int64]bool // Map для отслеживания ожидания файла в чате (chatID -> bool)
var serviceNum int

func getServices() []string {
	services := []string{
		"Al Dar",
		"Luna 22",
		"Condor",
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

	waitingForFile = make(map[int64]bool)

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

				//waitingForFile[chatID] = true // Установка флага ожидания файла
				serviceNum = 1
			case 2:
				msg := tgbotapi.NewMessage(chatID, "Вы выбрали сервис Luna")
				bot.Send(msg)
				msg = tgbotapi.NewMessage(chatID, "Пришлите файл в формате XLSX")
				bot.Send(msg)

				//waitingForFile[chatID] = true // Установка флага ожидания файла
				serviceNum = 2
			case 3:
				msg := tgbotapi.NewMessage(chatID, "Вы выбрали сервис Condor")
				bot.Send(msg)
				msg = tgbotapi.NewMessage(chatID, "Пришлите файл в формате XLSX")
				bot.Send(msg)

				//waitingForFile[chatID] = true // Установка флага ожидания файла
				serviceNum = 3
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
				getServiceLuna22(update.Message, bot)
			case 3:
				getServiceCondor(update.Message, bot)
			default:
				msg := tgbotapi.NewMessage(chatID, "Упс")
				bot.Send(msg)
			}
			//waitingForFile[chatID] = false // Сброс флага ожидания файла
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

func getServiceResponse(serviceNumber int) string {
	switch serviceNumber {
	case 1:
		return "Вы выбрали сервис №1"
	case 2:
		return "Вы выбрали сервис №2"
	case 3:
		return "Вы выбрали сервис №3"
	default:
		return "Некорректный выбор сервиса"
	}
}

func sendProcessingMessage(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Сейчас сделаю обновленный файл 🥰")
	bot.Send(msg)
}

func sendUpdateMessage(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Вот обновление 🤩")
	bot.Send(msg)
}

func sendCSVFile(bot *tgbotapi.BotAPI, chatID int64, xlsxBuffer *bytes.Buffer) {
	xlsxConfig := tgbotapi.FileBytes{
		Name:  "refactor_available_units.csv",
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
