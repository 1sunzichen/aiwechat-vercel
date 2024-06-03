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
	if msg == "tzw" || strings.Contains(msg, "tzs") || msg == "tza" || msg == "tzm" {
		returncontent := "您当前输入的是 " + msg + " \n " +
			msg
		if msg == "tzw" || msg == "tzm" {

			returncontent += "需要下载：谷歌浏览器和插件，插件下载后 需要解压\n" +
				"插件链接：链接: https://pan.baidu.com/s/1zNa5gnm9TbYH4WEf_OTNRw?pwd=y2g4 提取码: y2g4\n" +
				"打开谷歌浏览器 设置-扩展程序-打开开发者模式-加载已解压的扩展程序-选择下载的插件文件夹-确定\n" +
				"插件安装好,获取新脚本,在新脚本网站中搜索页搜索脚本\n" +
				"输入：某某视频脚本，b站，爱奇艺d等等\n" +
				"最后去网页 爱奇艺网站等需要看的 影片资源刷新。\n" +
				"如需帮助，请输入“会员帮助”\n"
		} else if msg == "tza" {
			returncontent += "链接: https://pan.baidu.com/s/1-G83nFLDw7k_89KFaFLghw?pwd=3b1i 提取码: 3b1i"
		} else if strings.Contains(msg, "tzs") {
			str := msg[3:]
			url, err := db.ChatDbInstance.GetVideoValue(str)
			if url == "" || err != nil {
				returncontent += "请尝试输入如：tzs哈尔滨一九四四 （哈尔滨一九四四 可更换 资源名字，如没有直接联系vx15210187668反映）\n"
			} else {
				returncontent += "!!!复制到浏览器打开链接，不用使用微信打开链接!!!，下方第三个按钮 可以选集" + url
			}
		}

		return returncontent
	} else if msg == "会员帮助" {
		return "链接: https://pan.baidu.com/s/19Q4q8Gh_2LqyJCGS1TG3GQ?pwd=skfn 提取码: skfn 复制这段内容后打开百度网盘手机App，操作更方便哦"
	} else if msg == "孙子宸" {
		return "链接: https://pan.baidu.com/s/1NNLlZ7XDJhvC27858lTltQ?pwd=ymyb 提取码: ymyb"
	}
	cfg := openai.DefaultConfig(s.token)
	cfg.BaseURL = "https://ai-yyds.com/v1"
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
	fmt.Println("content$$$$$$$:", content)
	msgs = append(msgs, openai.ChatCompletionMessage{Role: openai.ChatMessageRoleAssistant, Content: content})
	SaveMsgListWithDb(config.Bot_Type_Gpt, userID, msgs, s.toDbMsg)
	return content + "\n需要看电视，电影视频资源，输入如：tzs哈尔滨一九四四\n，\n"
}

func (s *SimpleGptChat) Chat(userID string, msg string) string {
	return WithTimeChat(userID, msg, s.chat)
}
