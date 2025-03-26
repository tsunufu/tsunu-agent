// parser.go
package main

import (
	"encoding/xml"
	"fmt"
	"regexp"
)

// ツールの種類を表す定数
const (
	ToolTypeListFile       = "list_file"
	ToolTypeReadFile       = "read_file"
	ToolTypeWriteFile      = "write_file"
	ToolTypeAskQuestion    = "ask_question"
	ToolTypeExecuteCommand = "execute_command"
	ToolTypeComplete       = "complete"
)

// ParseAndExecuteTool はLLMのレスポンスをパースしてツールを実行する
func ParseAndExecuteTool(response string) (ToolResponse, string, bool) {
	// XMLタグを抽出する正規表現
	re := regexp.MustCompile(`<([a-z_]+)>([\s\S]*?)</\1>`)
	match := re.FindStringSubmatch(response)

	if len(match) < 3 {
		return ToolResponse{
			Success: false,
			Message: "有効なツールが見つかりませんでした",
		}, "", false
	}

	toolType := match[1]
	toolContent := match[2]

	switch toolType {
	case ToolTypeListFile:
		var params ListFileParams
		if err := xml.Unmarshal([]byte(fmt.Sprintf("<%s>%s</%s>", toolType, toolContent, toolType)), &params); err != nil {
			return ToolResponse{
				Success: false,
				Message: fmt.Sprintf("パラメータのパースに失敗しました: %v", err),
			}, toolType, false
		}
		return ListFile(params), toolType, false

	case ToolTypeReadFile:
		var params ReadFileParams
		if err := xml.Unmarshal([]byte(fmt.Sprintf("<%s>%s</%s>", toolType, toolContent, toolType)), &params); err != nil {
			return ToolResponse{
				Success: false,
				Message: fmt.Sprintf("パラメータのパースに失敗しました: %v", err),
			}, toolType, false
		}
		return ReadFile(params), toolType, false

	case ToolTypeWriteFile:
		var params WriteFileParams
		if err := xml.Unmarshal([]byte(fmt.Sprintf("<%s>%s</%s>", toolType, toolContent, toolType)), &params); err != nil {
			return ToolResponse{
				Success: false,
				Message: fmt.Sprintf("パラメータのパースに失敗しました: %v", err),
			}, toolType, false
		}
		return WriteFile(params), toolType, false

	case ToolTypeAskQuestion:
		var params AskQuestionParams
		if err := xml.Unmarshal([]byte(fmt.Sprintf("<%s>%s</%s>", toolType, toolContent, toolType)), &params); err != nil {
			return ToolResponse{
				Success: false,
				Message: fmt.Sprintf("パラメータのパースに失敗しました: %v", err),
			}, toolType, false
		}
		return AskQuestion(params), toolType, false

	case ToolTypeExecuteCommand:
		var params ExecuteCommandParams
		if err := xml.Unmarshal([]byte(fmt.Sprintf("<%s>%s</%s>", toolType, toolContent, toolType)), &params); err != nil {
			return ToolResponse{
				Success: false,
				Message: fmt.Sprintf("パラメータのパースに失敗しました: %v", err),
			}, toolType, false
		}
		return ExecuteCommand(params), toolType, false

	case ToolTypeComplete:
		var params CompleteParams
		if err := xml.Unmarshal([]byte(fmt.Sprintf("<%s>%s</%s>", toolType, toolContent, toolType)), &params); err != nil {
			return ToolResponse{
				Success: false,
				Message: fmt.Sprintf("パラメータのパースに失敗しました: %v", err),
			}, toolType, false
		}
		return Complete(params), toolType, true

	default:
		return ToolResponse{
			Success: false,
			Message: fmt.Sprintf("未知のツールタイプ: %s", toolType),
		}, toolType, false
	}
}
