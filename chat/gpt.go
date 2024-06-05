package chat

import (
	"context"
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
		returncontent := ""
		str := msg[3:]
		url, err := db.ChatDbInstance.GetVideoValue(str)
		if url == "" || err != nil {
			returncontent += "请尝试输入如：tzs哈尔滨一九四四 （哈尔滨一九四四 可更换 资源名字，如没有直接联系vx15210187668反映）\n"
		} else {
			returncontent += "!!!复制到浏览器打开链接，不用使用微信打开链接!!!，下方第三个按钮 可以选集" + url
		}
		return returncontent
	} else if strings.ToUpper(msg) == "SJ" || strings.ToUpper(msg) == "圣经" {
		return "链接: https://pan.baidu.com/s/1iroQ_a-cPXf1the2Wz2NRw?pwd=w1ag 提取码: w1ag"
	}
	cfg := openai.DefaultConfig(s.token)

	client := openai.NewClientWithConfig(cfg)

	var msgs = GetMsgListWithDb(config.Bot_Type_Gpt, userID, openai.ChatCompletionMessage{Role: openai.ChatMessageRoleUser, Content: msg}, s.toDbMsg, s.toChatMsg)

	resp, err := client.CreateChatCompletion(context.Background(),
		openai.ChatCompletionRequest{
			Model:    s.getModel(),
			Messages: msgs,
		})

	if err != nil {
		return err.Error()
	}

	content := resp.Choices[0].Message.Content

	// msgs = append(msgs, openai.ChatCompletionMessage{Role: openai.ChatMessageRoleAssistant, Content: content})
	msgs = append(msgs, openai.ChatCompletionMessage{Role: openai.ChatMessageRoleAssistant, Content: ""})
	SaveMsgListWithDb(config.Bot_Type_Gpt, userID, msgs, s.toDbMsg)
	return content + "\n需要看电视，电影视频资源，输入如：tzs哈尔滨一九四四\n"

}

func (s *SimpleGptChat) Chat(userID string, msg string) string {
	return WithTimeChat(userID, msg, s.chat)
}
