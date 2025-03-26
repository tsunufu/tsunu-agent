// tool.go
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ListFileParams はlist_fileツールのパラメータ
type ListFileParams struct {
	Path      string `xml:"path"`
	Recursive string `xml:"recursive"`
}

// ReadFileParams はread_fileツールのパラメータ
type ReadFileParams struct {
	Path string `xml:"path"`
}

// WriteFileParams はwrite_fileツールのパラメータ
type WriteFileParams struct {
	Path    string `xml:"path"`
	Content string `xml:"content"`
}

// AskQuestionParams はask_questionツールのパラメータ
type AskQuestionParams struct {
	Question string `xml:"question"`
}

// ExecuteCommandParams はexecute_commandツールのパラメータ
type ExecuteCommandParams struct {
	Command          string `xml:"command"`
	RequiresApproval string `xml:"requires_approval"`
}

// CompleteParams はcompleteツールのパラメータ
type CompleteParams struct {
	Result string `xml:"result"`
}

// ToolResponse はツールの実行結果
type ToolResponse struct {
	Success bool
	Message string
}

// 1. ListFile - ディレクトリ内のファイル一覧を取得
func ListFile(params ListFileParams) ToolResponse {
	path := params.Path
	recursive := strings.ToLower(params.Recursive) == "true"

	var files []string
	var err error

	if recursive {
		err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			files = append(files, path)
			return nil
		})
	} else {
		entries, err := ioutil.ReadDir(path)
		if err == nil {
			for _, entry := range entries {
				files = append(files, filepath.Join(path, entry.Name()))
			}
		}
	}

	if err != nil {
		return ToolResponse{
			Success: false,
			Message: fmt.Sprintf("ディレクトリの読み取りに失敗しました: %v", err),
		}
	}

	result := fmt.Sprintf("ディレクトリ %s のファイル一覧:\n", path)
	for _, file := range files {
		result += fmt.Sprintf("- %s\n", file)
	}

	return ToolResponse{
		Success: true,
		Message: result,
	}
}

// 2. ReadFile - ファイルの内容を読み取る
func ReadFile(params ReadFileParams) ToolResponse {
	content, err := ioutil.ReadFile(params.Path)
	if err != nil {
		return ToolResponse{
			Success: false,
			Message: fmt.Sprintf("ファイルの読み取りに失敗しました: %v", err),
		}
	}

	return ToolResponse{
		Success: true,
		Message: string(content),
	}
}

// 3. WriteFile - ファイルに内容を書き込む
func WriteFile(params WriteFileParams) ToolResponse {
	dir := filepath.Dir(params.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return ToolResponse{
			Success: false,
			Message: fmt.Sprintf("ディレクトリの作成に失敗しました: %v", err),
		}
	}

	err := ioutil.WriteFile(params.Path, []byte(params.Content), 0644)
	if err != nil {
		return ToolResponse{
			Success: false,
			Message: fmt.Sprintf("ファイルの書き込みに失敗しました: %v", err),
		}
	}

	return ToolResponse{
		Success: true,
		Message: fmt.Sprintf("ファイル %s に書き込みました", params.Path),
	}
}

// 4. AskQuestion - ユーザーに質問する
func AskQuestion(params AskQuestionParams) ToolResponse {
	fmt.Printf("\n質問: %s\n回答: ", params.Question)

	var answer string
	fmt.Scanln(&answer)

	return ToolResponse{
		Success: true,
		Message: fmt.Sprintf("ユーザーの回答: %s", answer),
	}
}

// 5. ExecuteCommand - コマンドを実行する
func ExecuteCommand(params ExecuteCommandParams) ToolResponse {
	requiresApproval := strings.ToLower(params.RequiresApproval) == "true"

	if requiresApproval {
		fmt.Printf("\n以下のコマンドを実行しますか？\n%s\n[y/n]: ", params.Command)

		var answer string
		fmt.Scanln(&answer)

		if strings.ToLower(answer) != "y" {
			return ToolResponse{
				Success: false,
				Message: "コマンドの実行がキャンセルされました",
			}
		}
	}

	cmd := exec.Command("sh", "-c", params.Command)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return ToolResponse{
			Success: false,
			Message: fmt.Sprintf("コマンドの実行に失敗しました: %v\n出力: %s", err, string(output)),
		}
	}

	return ToolResponse{
		Success: true,
		Message: fmt.Sprintf("コマンドの実行結果:\n%s", string(output)),
	}
}

// 6. Complete - タスクの完了を示す
func Complete(params CompleteParams) ToolResponse {
	return ToolResponse{
		Success: true,
		Message: fmt.Sprintf("タスク完了: %s", params.Result),
	}
}
