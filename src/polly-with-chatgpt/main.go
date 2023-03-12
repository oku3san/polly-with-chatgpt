package main

import (
    "context"
    "fmt"
    openai "github.com/sashabaranov/go-openai"
    "os"
)

func main() {
    token := os.Getenv("OPENAI_API_KEY")
    client := openai.NewClient(token)
    resp, err := client.CreateChatCompletion(
        context.Background(),
        openai.ChatCompletionRequest{
            Model: openai.GPT3Dot5Turbo,
            Messages: []openai.ChatCompletionMessage{
                {
                    Role:    openai.ChatMessageRoleUser,
                    Content: "Hello!",
                },
            },
        },
    )

    if err != nil {
        fmt.Printf("ChatCompletion error: %v\n", err)
        return
    }

    fmt.Println(resp.Choices[0].Message.Content)
}
