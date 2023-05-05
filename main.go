package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"

	"github.com/nickname76/telegrambot"
)

func main() {
	fmt.Println("Chat GPT Bot")
	jsonFile, err := os.Open("config.json")
	if err != nil {
		log.Printf("Error openinig json file: %v", err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]string
	json.Unmarshal([]byte(byteValue), &result)
	var T_TOKEN string = result["TELEGRAM_TOKEN"]
	var GPT_TOKEN string = result["CHAT_GPT_TOKEN"]

	bot := NewTelegramBot(T_TOKEN)

	bot.Api.SetMyCommands(&telegrambot.SetMyCommandsParams{
		Commands: []*telegrambot.BotCommand{
			{
				Command:     "new",
				Description: "Command to start a new session",
			},
		},
	})

	gptClient := NewCli(GPT_TOKEN)

	stop := telegrambot.StartReceivingUpdates(bot.Api, func(update *telegrambot.Update, err error) {
		if err != nil {
			log.Printf("Error: %v", err)
			return
		}
		msg := update.Message
		if msg == nil {
			return
		}
		voice, text := bot.CheckMessageType(msg)
		if voice != nil {
			bot.HandleVoice(msg, gptClient)
		}
		if text != "" {
			bot.HandleText(msg, gptClient)
		}
	})
	log.Printf("Started on %v", bot.User.Username)
	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, os.Interrupt)
	<-exitCh
	// Waits for all updates handling to complete
	stop()
}

func downloadAndConevrt(file, path string) (string, string) {
	//This object represents a file ready to be
	//downloaded.
	//The file can be downloaded via the link https://api.telegram.org/file/bot<token>/<file_path>.//
	responce, err := http.Get(file + path)
	fmt.Println("get: ", file+"getFile")
	if err != nil {
		fmt.Println("Error by downloading", err)
	}

	defer responce.Body.Close()
	r := []rune(path)
	r = r[6 : len(r)-0]
	ogaFile := string(r)

	fileMp, err := os.Create(ogaFile)
	if err != nil {
		fmt.Println("Error by creating", err)
	}
	defer fileMp.Close()
	fmt.Println("Body:", &responce.Body)

	_, err = io.Copy(fileMp, responce.Body)
	if err != nil {
		fmt.Println("Error by copying", err)
	}
	mp3File := toMp3(ogaFile)
	return mp3File, ogaFile
}

func toMp3(path string) string {
	fmp3 := path
	runf := []rune(fmp3)
	runf = runf[:len(runf)-4]
	fmp3 = string(runf)
	s := fmt.Sprintf("ffmpeg -i %s %s.mp3", path, fmp3)
	args := strings.Split(s, " ")
	cmd := exec.Command(args[0], args[1:]...)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error by converting", err)
	}
	return fmp3 + ".mp3"
}

// func openAi(file string) (string, string) {
// 	c := openai.NewClient("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")
// 	ctx := context.Background()

// 	req := openai.AudioRequest{
// 		Model:    openai.Whisper1,
// 		FilePath: file,
// 	}
// 	resp, err := c.CreateTranscription(ctx, req)
// 	if err != nil {
// 		fmt.Printf("Transcription error: %v\n", err)
// 		return "", ""
// 	}
// 	fmt.Println(resp.Text)

// 	chatResp := chat(c, resp.Text)
// 	return resp.Text, chatResp
// }

// func chat(c *openai.Client, text string) string {
// 	ctx := context.Background()
// 	msgs := []openai.ChatCompletionMessage{
// 		{
// 			Role:    "user",
// 			Content: text,
// 		},
// 	}
// 	req := openai.ChatCompletionRequest{
// 		Model:    "gpt-3.5-turbo",
// 		Messages: msgs,
// 	}
// 	response, err := c.CreateChatCompletion(ctx, req)
// 	if err != nil {
// 		fmt.Println("Error getting response from gpt: ", err)
// 	}

//		return response.Choices[0].Message.Content
//	}
// func openAiText(file string) string {
// 	c := openai.NewClient("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")

// 	chatResp := chat(c, file)
// 	return chatResp
// }
