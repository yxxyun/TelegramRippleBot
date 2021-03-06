package main

import (
	"log"
	"telegram-bot-api"
	"strings"
)

func sendMessage(chatId int64, text string, keyboard interface{}) {
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.DisableWebPagePreview = true
	_, ok := keyboard.(tgbotapi.ReplyKeyboardMarkup)
	if ok {
		msg.ReplyMarkup = keyboard
	} else {
		_, ok = keyboard.(tgbotapi.InlineKeyboardMarkup)
		if ok {
			msg.ReplyMarkup = &keyboard
		} else {
			msg.ReplyMarkup = nil
		}
	}

	_, err := bot.Send(msg)
	if err != nil {
		log.Print(err)
		if strings.Contains(err.Error(), " can't parse entities"){
			msg.ParseMode = ""
			bot.Send(msg)
		}
	}else {
		log.Printf("[Bot] SENT %s TO %v", msg.Text, msg.ChatID)
	}

}

func sendPhoto(photo tgbotapi.PhotoConfig) string {
	response, err := bot.Send(photo)
	if err != nil {
		log.Print(err)
		return ""
	}
	log.Printf("[Bot] PHOTO %s TO %v", photo.FileID, photo.ChatID)
	return (*(response.Photo))[0].FileID
}

func setUploadingPhoto(id int64) {
	_, err := bot.Send(tgbotapi.NewChatAction(id, tgbotapi.ChatUploadPhoto))
	if err != nil {
		log.Print(err)
	}
}
