package main

import (
	"context"
	"fmt"

	"github.com/nickname76/telegrambot"
	"github.com/sashabaranov/go-openai"
)

type Msg struct {
	message []*telegrambot.Message
}

type Client struct {
	c *openai.Client
}

func NewCli(t string) *Client {
	return &Client{
		c: openai.NewClient(t),
	}
}

func NewMsg(request string) *Msg {
	return &Msg{
		message: make([]*telegrambot.Message, 0),
	}
}

func (m *Msg) getMsgs() []*telegrambot.Message {
	return m.message
}

func (c *Client) speechToText(f *File) (string, string) {
	text, responce := c.openAi(f.FileMp3)
	return text, responce
}

func (c *Client) openAi(file string) (string, string) {

	ctx := context.Background()

	req := openai.AudioRequest{
		Model:    openai.Whisper1,
		FilePath: file,
	}
	resp, err := c.c.CreateTranscription(ctx, req)
	if err != nil {
		fmt.Printf("Transcription error: %v\n", err)
		return "", ""
	}
	fmt.Println(resp.Text)

	chatResp := c.chat(resp.Text)
	return resp.Text, chatResp
}
func (c *Client) chat(text string) string {
	ctx := context.Background()
	msgs := []openai.ChatCompletionMessage{
		{
			Role:    "user",
			Content: text,
		},
	}
	req := openai.ChatCompletionRequest{
		Model:    "gpt-3.5-turbo",
		Messages: msgs,
	}
	response, err := c.c.CreateChatCompletion(ctx, req)
	if err != nil {
		fmt.Println("Error getting response from gpt: ", err)
	}

	return response.Choices[0].Message.Content
}

func (c *Client) openAiText(file string) string {
	chatResp := c.chat(file)
	return chatResp
}
