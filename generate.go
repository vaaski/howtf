package main

import (
	"context"
	"errors"
	"io"
	"log"
	"runtime"

	"github.com/sashabaranov/go-openai"
)

var (
	LLM_GENERATE_MESSAGE = "You solve problems by generating console commands for a developer. Do not explain anything, just provide the command for " + getUserShell() + " on " + runtime.GOOS
	LLM_EXPLAIN_MESSAGE  = "You briefly explain a given console command to a developer. Do not provide any code, just explain the command and its flags in markdown bullet points. The user is using" + getUserShell() + " on " + runtime.GOOS
)

func makeOpenAICall(token string, model string, query string, responseChannel chan queryResponse, systemMessage string) {
	client := openai.NewClient(token)
	ctx := context.Background()

	request := openai.ChatCompletionRequest{
		Model:       model,
		Temperature: 0,
		MaxTokens:   512,
		Stream:      true,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemMessage,
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

func generateCommand(token string, model string, query string, responseChannel chan queryResponse) {
	makeOpenAICall(token, model, query, responseChannel, LLM_GENERATE_MESSAGE)
}

func explainCommand(token string, model string, query string, responseChannel chan queryResponse) {
	makeOpenAICall(token, model, query, responseChannel, LLM_EXPLAIN_MESSAGE)
}
