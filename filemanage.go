package main

import (
	"fmt"
	"os"

	"github.com/nickname76/telegrambot"
)

type File struct {
	FileOga string
	FileMp3 string
}

func NewFile(bot *Bot, msg *telegrambot.Message) (*File, error) {
	f, err := bot.Api.GetFile(&telegrambot.GetFileParams{
		FileID: msg.Voice.FileID,
	})
	if err != nil {
		fmt.Println("Error getting the voice file!!")
		return nil, err
	}
	telegramFileLink := fmt.Sprintf("https://api.telegram.org/file/%s%v/", "bot", bot.token)
	fmt.Println(telegramFileLink)
	mp3File, ogaFile := downloadAndConevrt(telegramFileLink, f.FilePath)

	//text, resp := openAi(mp3File)

	return &File{
		FileOga: ogaFile,
		FileMp3: mp3File,
	}, nil
}

func (f *File) DelOga() error {
	err := os.Remove(f.FileOga) // remove a single file
	if err != nil {
		return err
	}
	return nil
}
func (f *File) DelMp3() error {
	err := os.Remove(f.FileMp3) // remove a single file
	if err != nil {
		return err
	}
	return nil
}
