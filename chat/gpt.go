package chat

import (
	"context"
	"fmt"
	"strings"

	"os"

	// Videourl "github.com/pwh-pwh/aiwechat-vercel/chat/videourl"

	"github.com/pwh-pwh/aiwechat-vercel/config"
	"github.com/pwh-pwh/aiwechat-vercel/db"
	"github.com/sashabaranov/go-openai"
)

type SimpleGptChat struct {
	token string
	url   string
	SimpleChat
	maxTokens int
	BaseChat  SimpleChat
}

func (s *SimpleGptChat) toDbMsg(msg openai.ChatCompletionMessage) db.Msg {
	return db.Msg{
		Role: msg.Role,
		Msg:  msg.Content,
	}
}

func (s *SimpleGptChat) toChatMsg(msg db.Msg) openai.ChatCompletionMessage {
	return openai.ChatCompletionMessage{
		Role:    msg.Role,
		Content: msg.Msg,
	}
}

func (s *SimpleGptChat) getModel() string {
	model := os.Getenv("gptModel")
	if model == "" {
		model = "gpt-3.5-turbo"
	}
	return model
}

func (s *SimpleGptChat) chat(userID, msg string) string {
	if strings.Contains(msg, "tzs") {

		return "需要看电视，电影视频资源，输入如：tzs哈尔滨一九四四\n，\n"
	}
	cfg := openai.DefaultConfig(s.token)

	client := openai.NewClientWithConfig(cfg)

	var msgs = GetMsgListWithDb(config.Bot_Type_Gpt, userID, openai.ChatCompletionMessage{Role: openai.ChatMessageRoleUser, Content: msg}, s.toDbMsg, s.toChatMsg)

	resp, err := client.CreateChatCompletion(context.Background(),
		openai.ChatCompletionRequest{
			Model:    s.getModel(),
			Messages: msgs,
		})
	fmt.Println("content$$$$2$$$before:")
	if err != nil {
		return err.Error()
	}

	content := resp.Choices[0].Message.Content
	fmt.Println("content$$$$$$$:", content)
	// msgs = append(msgs, openai.ChatCompletionMessage{Role: openai.ChatMessageRoleAssistant, Content: content})
	msgs = append(msgs, openai.ChatCompletionMessage{Role: openai.ChatMessageRoleAssistant, Content: ""})
	SaveMsgListWithDb(config.Bot_Type_Gpt, userID, msgs, s.toDbMsg)
	return content + "\n需要看电视，电影视频资源，输入如：tzs哈尔滨一九四四\n，\n"

}

func (s *SimpleGptChat) Chat(userID string, msg string) string {
	return WithTimeChat(userID, msg, s.chat)
}
