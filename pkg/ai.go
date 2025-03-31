package pkg

import (
	"context"
	"log"
	"os"

	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

func GetAiResponse(prompt string) (string, error) {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	type CommitOutput struct {
		CommitMessage string `json:"commit_message"`
	}
	var result CommitOutput
	schema, err := jsonschema.GenerateSchemaForType(result)

	if err != nil {
		return "", err
	}

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "Você é um assistente que recebe uma DIFF ou um novo arquivo e retorna um texto para commit, seguindo o padrão de Conventional Commits",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
				JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
					Name:   "commit_message",
					Schema: schema,
					Strict: true,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	if err != nil {
		log.Fatalf("CreateChatCompletion error: %v", err)
	}
	err = schema.Unmarshal(resp.Choices[0].Message.Content, &result)
	if err != nil {
		log.Fatalf("Unmarshal schema error: %v", err)
	}

	return result.CommitMessage, nil
}
