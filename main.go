// main.go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

func main() {
	// OpenAI APIキーを環境変数から取得
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("OPENAI_API_KEYが設定されていません")
		return
	}

	// OpenAI APIクライアントを初期化
	client := openai.NewClient(apiKey)

	// システムプロンプトを設定
	systemPrompt := `あなたはコーディングエージェントです。以下のツールを使ってタスクを完了してください：

# ListFile
ディレクトリ内のファイル一覧を取得します。
<list_file>
<path>ディレクトリのパス</path>
<recursive>true または false</recursive>
</list_file>

# ReadFile
ファイルの内容を読み取ります。
<read_file>
<path>ファイルのパス</path>
</read_file>

# WriteFile
ファイルに内容を書き込みます。
<write_file>
<path>ファイルのパス</path>
<content>
書き込む内容
</content>
</write_file>

# AskQuestion
ユーザーに質問します。
<ask_question>
<question>質問内容</question>
</ask_question>

# ExecuteCommand
コマンドを実行します。
<execute_command>
<command>実行するコマンド</command>
<requires_approval>true または false</requires_approval>
</execute_command>

# Complete
タスクの完了を示します。
<complete>
<result>タスクの結果や成果物の説明</result>
</complete>

必ず上記のいずれかのツールを使用してください。ツールを使わずに直接回答しないでください。`

	// ユーザーからのタスク入力を受け取る
	fmt.Println("コーディングエージェントにタスクを入力してください:")
	var userTask string
	fmt.Scanln(&userTask)

	// 会話履歴を初期化
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: userTask,
		},
	}

	// メインループ
	isComplete := false
	for !isComplete {
		// LLMにリクエストを送信
		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:    openai.GPT4,
				Messages: messages,
			},
		)

		if err != nil {
			fmt.Printf("エラーが発生しました: %v\n", err)
			return
		}

		// LLMのレスポンスを取得
		assistantResponse := resp.Choices[0].Message.Content

		// レスポンスをパースしてツールを実行
		toolResponse, toolType, complete := ParseAndExecuteTool(assistantResponse)

		// ツールの実行結果をメッセージに追加
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: assistantResponse,
		})

		// ツールの実行結果をユーザーに表示
		if toolType != ToolTypeAskQuestion && toolType != ToolTypeExecuteCommand {
			fmt.Printf("\n[%s] %s\n", toolType, toolResponse.Message)
		}

		// ツールの実行結果をメッセージに追加
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf("[%s Result] %s", toolType, toolResponse.Message),
		})

		// Completeツールが実行された場合はループを終了
		if complete {
			isComplete = true
		}
	}
}
