package main

import (
	"context"
	"errors"
	"io"
	"log"
	"runtime"

	"github.com/sashabaranov/go-openai"
)

var LLM_SYSTEM_MESSAGE = "You solve problems by generating console commands for a developer. Do not explain anything, just provide the command for " + getUserShell() + " on " + runtime.GOOS

func generateGPT(token string, query string, responseChannel chan queryResponse) {
	client := openai.NewClient(token)
	ctx := context.Background()

	request := openai.ChatCompletionRequest{
		Model:       openai.GPT4o,
		Temperature: 0,
		MaxTokens:   512,
		Stream:      true,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: LLM_SYSTEM_MESSAGE,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: query,
			},
		},
	}

	stream, err := client.CreateChatCompletionStream(ctx, request)
	if err != nil {
		panic(err)
	}
	defer stream.Close()

	var responseAccumulator string
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			panic(err)
		}

		responseAccumulator += response.Choices[0].Delta.Content
		log.Println("responseAccumulator", responseAccumulator)
		responseChannel <- queryResponse(responseAccumulator)
	}
}
