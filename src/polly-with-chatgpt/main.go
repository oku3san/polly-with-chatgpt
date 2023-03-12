package main

import (
    "context"
    "fmt"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/polly"
    "github.com/aws/aws-sdk-go-v2/service/polly/types"
    openai "github.com/sashabaranov/go-openai"
    "io/ioutil"
    "log"
    "os"
    "os/exec"
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

    msg := resp.Choices[0].Message.Content

    // AWS Config の情報を取得
    ctx := context.Background()
    cfg, err := config.LoadDefaultConfig(ctx)
    if err != nil {
        log.Fatalf("AWS セッションエラー: %v", err)
    }

    pollyClient := polly.NewFromConfig(cfg)

    input := &polly.SynthesizeSpeechInput{
        OutputFormat: types.OutputFormatMp3,
        Text:         aws.String(msg),
        VoiceId:      types.VoiceIdJoanna,
    }

    output, err := pollyClient.SynthesizeSpeech(context.Background(), input)
    if err != nil {
        // 変換に失敗した場合のエラー処理
        fmt.Println(err.Error())
    }

    // 結果をmp3ファイルで出力
    content, err := ioutil.ReadAll(output.AudioStream)
    ioutil.WriteFile("/tmp/gopolly.mp3", content, os.ModePerm)

    // mp3ファイルを再生
    exerr := exec.Command("afplay", "/tmp/gopolly.mp3").Run()
    if exerr != nil {
        fmt.Println(exerr.Error())
        return
    }
}
