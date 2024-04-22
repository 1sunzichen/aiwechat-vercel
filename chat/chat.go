package chat

import (
	"fmt"
	"os"
	"time"

	"github.com/pwh-pwh/aiwechat-vercel/db"
	"github.com/sashabaranov/go-openai"

	"github.com/pwh-pwh/aiwechat-vercel/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

type BaseChat interface {
	Chat(userID string, msg string) string
	HandleMediaMsg(msg *message.MixMessage) string
}
type SimpleChat struct {
}

func (s SimpleChat) Chat(userID string, msg string) string {
	panic("implement me")
}

func (s SimpleChat) HandleMediaMsg(msg *message.MixMessage) string {
	fmt.Println(msg.Content, "msgContent")
	switch msg.MsgType {
	case message.MsgTypeImage:
		return msg.PicURL
	case message.MsgTypeEvent:
		if msg.Event == message.EventSubscribe {
			subText := os.Getenv("subscribe")
			if subText == "" {
				subText = "哇，又有帅哥美女关注我啦😄,需要各大网站视频会员。电脑最好，手机体验不太好。" +
					" window 电脑 请输入 tzw 关键字，苹果电脑 输入 tzm,如果是 安卓手机使用请输入 tzs，苹果手机和ipad 输入tza"
			}
			return subText
		} else if msg.Content == "tzw" || msg.Content == "tzm" || msg.Content == "tza" || msg.Content == "tzs" {

			return msg.Content
		} else {
			a := msg.MsgType
			return string(a)
		}
	default:
		return "未支持的类型"
	}
}

// 加入超时控制
func WithTimeChat(userID, msg string, f func(userID, msg string) string) string {
	if _, ok := config.Cache.Load(userID); ok {
		rAny, _ := config.Cache.Load(userID)
		r := rAny.(string)
		config.Cache.Delete(userID)
		return r
	}
	resChan := make(chan string)
	go func() {
		resChan <- f(userID, msg)
	}()
	select {
	case res := <-resChan:
		return res
	case <-time.After(5 * time.Second):
		config.Cache.Store(userID, <-resChan)
		return ""
	}
}

type ErrorChat struct {
	errMsg string
}

func (e *ErrorChat) HandleMediaMsg(msg *message.MixMessage) string {
	return "20秒后重新尝试" + e.errMsg
}

func (e *ErrorChat) Chat(userID string, msg string) string {
	return "20秒后重新尝试" + e.errMsg
}

func GetChatBot(botType string) BaseChat {
	if botType == "" {
		botType = config.GetBotType()
	}
	var err error
	botType, err = config.CheckBotConfig(botType)
	if err != nil {
		return &ErrorChat{
			errMsg: err.Error(),
		}
	}

	switch botType {
	case config.Bot_Type_Gpt:
		url := os.Getenv("GPT_URL")
		if url == "" {
			url = "https://api.openai.com/v1/"
		}
		return &SimpleGptChat{
			token:      os.Getenv("GPT_TOKEN"),
			url:        url,
			SimpleChat: SimpleChat{},
		}
	case config.Bot_Type_Spark:
		config, _ := config.GetSparkConfig()
		return &SparkChat{
			BaseChat: SimpleChat{},
			Config:   config,
		}
	case config.Bot_Type_Qwen:
		config, _ := config.GetQwenConfig()
		return &QwenChat{
			BaseChat: SimpleChat{},
			Config:   config,
		}
	default:
		return &Echo{}
	}
}

type ChatMsg interface {
	openai.ChatCompletionMessage | QwenMessage | SparkMessage
}

func GetMsgListWithDb[T ChatMsg](botType, userId string, msg T, f func(msg T) db.Msg, f2 func(msg db.Msg) T) []T {
	if db.ChatDbInstance != nil {
		list, err := db.ChatDbInstance.GetMsgList(botType, userId)
		if err == nil {
			list = append(list, f(msg))
			r := make([]T, 0)
			for _, msg := range list {
				r = append(r, f2(msg))
			}
			return r
		}
	}
	return []T{msg}
}

func SaveMsgListWithDb[T ChatMsg](botType, userId string, msgList []T, f func(msg T) db.Msg) {
	if db.ChatDbInstance != nil {
		go func() {
			list := make([]db.Msg, 0)
			for _, msg := range msgList {
				list = append(list, f(msg))
			}
			db.ChatDbInstance.SetMsgList(botType, userId, list)
		}()
	}
}
