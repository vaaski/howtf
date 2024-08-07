package main

import (
	"context"
	"errors"
	"io"
	"log"
	"regexp"
	"runtime"

	"github.com/sashabaranov/go-openai"
)

var LLM_SYSTEM_MESSAGE = "You're a generator for console commands for a developer. Do not explain anything, just provide the command for " + getUserShell() + " on " + runtime.GOOS

const LLM_PRETEXT = "I want to solve the following problem:\n"

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
				Content: LLM_PRETEXT + query,
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
		log.Println(responseAccumulator)
		responseChannel <- queryResponse(extractMarkdownMaybe(responseAccumulator))
	}
}

var MARKDOWN_REGEX = regexp.MustCompile("```(?:.*\n)?(.+)\n?```|`(.+)`")

func extractMarkdownMaybe(s string) string {
	matches := MARKDOWN_REGEX.FindStringSubmatch(s)

	if len(matches) <= 1 {
		return s
	} else if len(matches[1]) > 0 {
		return matches[1]
	} else if len(matches[2]) > 0 {
		return matches[2]
	}

	return s
}
