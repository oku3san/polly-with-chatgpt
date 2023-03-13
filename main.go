package main

import (
    "bufio"
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
    reader := bufio.NewReader(os.Stdin)
    question, _ := reader.ReadString('\n')

    answer, err := askChatGpt(question)
    if err != nil {
        log.Fatalf("ChatCompletion error: %v\n", err)
    }

    generateMp3WithPolly(answer.Choices[0].Message.Content)

    // mp3ファイルを再生
    excerr := exec.Command("afplay", "/tmp/gopolly.mp3").Run()
    if excerr != nil {
        fmt.Println(excerr.Error())
        return
    }
}

func askChatGpt(question string) (openai.ChatCompletionResponse, error) {
    token := os.Getenv("OPENAI_API_KEY")
    client := openai.NewClient(token)
    response, err := client.CreateChatCompletion(
        context.Background(),
        openai.ChatCompletionRequest{
            Model: openai.GPT3Dot5Turbo,
            Messages: []openai.ChatCompletionMessage{
                {
                    Role:    openai.ChatMessageRoleUser,
                    Content: question,
                },
            },
        },
    )
    if err != nil {
        return openai.ChatCompletionResponse{}, err
    }

    //msg := resp.Choices[0].Message.Content
    return response, nil
}

func generateMp3WithPolly(text string) {
    // AWS Config の情報を取得
    ctx := context.Background()
    cfg, err := config.LoadDefaultConfig(ctx)
    if err != nil {
        log.Fatalf("AWS セッションエラー: %v", err)
    }

    pollyClient := polly.NewFromConfig(cfg)

    input := &polly.SynthesizeSpeechInput{
        OutputFormat: types.OutputFormatMp3,
        Text:         aws.String(text),
        VoiceId:      types.VoiceIdMizuki,
    }

    output, err := pollyClient.SynthesizeSpeech(context.Background(), input)
    if err != nil {
        // 変換に失敗した場合のエラー処理
        fmt.Println(err.Error())
    }

    // 結果をmp3ファイルで出力
    content, err := ioutil.ReadAll(output.AudioStream)
    ioutil.WriteFile("/tmp/gopolly.mp3", content, os.ModePerm)

    return
}
