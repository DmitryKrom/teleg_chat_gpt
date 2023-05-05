package main

import (
	"fmt"
	"log"

	"github.com/nickname76/telegrambot"
)

type Bot struct {
	Api   *telegrambot.API
	User  *telegrambot.User
	token string
	err   error
}

func NewTelegramBot(token string) *Bot {
	fmt.Println("Chat GPT Bot")
	api, me, err := telegrambot.NewAPI(token)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Println("token: ", me.ID)

	return &Bot{
		Api:   api,
		User:  me,
		token: token,
		err:   err,
	}

}

func (b *Bot) CheckMessageType(msg *telegrambot.Message) (*telegrambot.Voice, string) {

	if msg.Voice != nil {
		return msg.Voice, ""
	}
	if msg.Text != "" {
		return nil, msg.Text
	}

	return nil, ""
}

func (b *Bot) HandleVoice(msg *telegrambot.Message, gptClient *Client) {

	files, err := NewFile(b, msg)
	if err != nil {
		log.Printf("Error creating voice files: %v", err)
	}

	text, resp := gptClient.speechToText(files)

	_, err = b.Api.SendMessage(&telegrambot.SendMessageParams{
		ChatID: msg.Chat.ID,
		Text:   fmt.Sprintf("Hi %v,\nDeine Frage: %s", msg.From.FirstName, text),
		ReplyMarkup: &telegrambot.ReplyKeyboardMarkup{
			Keyboard: [][]*telegrambot.KeyboardButton{{
				{
					Text: "Hallo",
				},
			}},
			ResizeKeyboard:  true,
			OneTimeKeyboard: false,
		},
	})
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	_, err = b.Api.SendMessage(&telegrambot.SendMessageParams{
		ChatID: msg.Chat.ID,
		Text:   fmt.Sprintf("Meine Antwort:\n %s", resp),
		ReplyMarkup: &telegrambot.ReplyKeyboardMarkup{
			Keyboard: [][]*telegrambot.KeyboardButton{{
				{
					Text: "Hallo",
				},
			}},
			ResizeKeyboard:  true,
			OneTimeKeyboard: false,
		},
	})
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	files.DelOga()
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	files.DelMp3()
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
}

func (b *Bot) HandleText(text *telegrambot.Message, gptClient *Client) {
	fmt.Println("text!!!!!", b.User.FirstName, b.User.Username)
	antwort := gptClient.openAiText(text.Text)
	_, err := b.Api.SendMessage(&telegrambot.SendMessageParams{
		ChatID: text.Chat.ID,
		Text:   fmt.Sprintf("Hi %v, I am %v and %v\nDein Text: %s", text.From.FirstName, b.User.FirstName, b.User.ID, text.Text),
		ReplyMarkup: &telegrambot.ReplyKeyboardMarkup{
			Keyboard: [][]*telegrambot.KeyboardButton{{
				{
					Text: "Hallo",
				},
			}},
			ResizeKeyboard:  true,
			OneTimeKeyboard: false,
		},
	})
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	_, err = b.Api.SendMessage(&telegrambot.SendMessageParams{
		ChatID: text.Chat.ID,
		Text:   fmt.Sprintf("Meine Antwort: %s", antwort),
		ReplyMarkup: &telegrambot.ReplyKeyboardMarkup{
			Keyboard: [][]*telegrambot.KeyboardButton{{
				{
					Text: "Hallo",
				},
			}},
			ResizeKeyboard:  true,
			OneTimeKeyboard: false,
		},
	})
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
}
