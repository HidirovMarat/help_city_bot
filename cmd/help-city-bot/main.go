package main

import (
	"fmt"
	"log"
	"time"

	"help-city-bot/internal/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"help-city-bot/internal/http/request"
)

var gBot *tgbotapi.BotAPI
var gToken string
var gChatId int64

func main() {
	con := config.MustLoad()

	if gToken = con.SigningKey; gToken == "" {
		panic(fmt.Errorf("Нет конфига и токена %s", TOKEN_NAME_IN_OS))
	}

	var err error

	if gBot, err = tgbotapi.NewBotAPI(gToken); err != nil {
		log.Panic(err)
	}

	gBot.Debug = true

	log.Printf("Authorized on account %s", gBot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = UPDATE_CONFIG_TIMEOUT

	for update := range gBot.GetUpdatesChan(updateConfig) {
		if isCallbackQuery(&update) {
			updateProcessing(&update)
		}

		if isStartMessage(&update) {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			gChatId = update.Message.Chat.ID
			askToPrintIntro()
		}

	}
}

func updateProcessing(update *tgbotapi.Update) {
	choiceCode := update.CallbackQuery.Data
	log.Printf("[%T] %s", time.Now(), choiceCode)

	switch choiceCode {
	case BUTTON_CODE_PRINT_INTRO:
		printIntro(update)
		showMenu(update)
	case BUTTON_CODE_SKIP_INTRO:
		showMenu(update)
	case BUTTON_CODE_GARBAGE, BUTTON_CODE_WATER, BUTTON_CODE_LIGHT:
		complaintAbout(choiceCode, update)
		askToMessage("Спасибо за жалобу")
	}
}

func complaintAbout(typeComplaint string, update *tgbotapi.Update) {
	askToMessage("Введите имя")

	firstName := update.Message.Text

	askToMessage("Введите фамилию")

	lastName := update.Message.Text

	askToMessage("Опишете жалобу")

	complaint := update.Message.Text

	request.SendComplaint(typeComplaint, firstName, lastName, complaint)
}

func showMenu(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(gChatId, "Выбери один из вариантов")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		getKeyboardRow(BUTTON_TEXT_GARBAGE, BUTTON_CODE_GARBAGE),
		getKeyboardRow(BUTTON_TEXT_WATER, BUTTON_CODE_WATER),
	)

	gBot.Send(msg)
}

func printIntro(update *tgbotapi.Update) {

	message := TEXT_INTRO

	askToMessage(message)
}

func isCallbackQuery(update *tgbotapi.Update) bool {
	return update.CallbackQuery != nil && update.CallbackQuery.Data != ""
}

func isStartMessage(update *tgbotapi.Update) bool {
	return update.Message != nil && update.Message.Text == "/start"
}

func getKeyboardRow(buttonText, buttonCode string) []tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(buttonText, buttonCode))
}

func askToPrintIntro() {
	msg := tgbotapi.NewMessage(gChatId, "Во вступительных сообщениях вы можете ознакомиться с назначением этого бота. Что вы думаете?")

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		getKeyboardRow(BUTTON_TEXT_PRINT_INTRO, BUTTON_CODE_PRINT_INTRO),
		getKeyboardRow(BUTTON_TEXT_SKIP_INTRO, BUTTON_CODE_SKIP_INTRO),
	)

	gBot.Send(msg)
}

func askToMessage(message string) {
	msg := tgbotapi.NewMessage(gChatId, message)

	gBot.Send(msg)
}
