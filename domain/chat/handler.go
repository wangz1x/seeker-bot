// handler - 2024/12/16
// Author: wangzx
// Description:

package chat

import (
	"encoding/xml"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"io"
	"log/slog"
	"net/http"
	"seeker-bot/m/db"
	"time"
)

const baseURL = "https://api.deepseek.com"
const token = "sk-43675ee878c74c6d91ad5ba4500fcdb1"

func Chat(c *gin.Context) {

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		responseXML(c, err.Error(), "", "")
		return
	}
	defer c.Request.Body.Close()

	var msg Message
	err = xml.Unmarshal(body, &msg)
	if err != nil {
		responseXML(c, err.Error(), "", "")
		return
	}

	slog.Info("receive message: %s\nfrom user: %s", msg.Content, msg.FromUserName)
	h := db.History{
		UserID:  msg.FromUserName,
		Role:    "user",
		Content: msg.Content,
	}

	result := db.DB.Create(&h)
	if result.Error != nil {
		responseXML(c, result.Error.Error(), msg.FromUserName, msg.ToUserName)
		return
	}

	var records []db.History
	result = db.DB.Where("user_id = ?", msg.FromUserName).Order("created_at desc").Limit(10).Find(&records)
	if result.Error != nil {
		responseXML(c, result.Error.Error(), msg.FromUserName, msg.ToUserName)
		return
	}

	// 查的时候是倒叙的，也就是新的在最前边，但是传话的话需要按时间顺序，所以再转一下
	var messages []openai.ChatCompletionMessage
	for i := len(records) - 1; i >= 0; i-- {
		record := records[i]
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    record.Role,
			Content: record.Content,
		})
	}

	defaultConfig := openai.DefaultConfig(token)
	defaultConfig.BaseURL = baseURL
	client := openai.NewClientWithConfig(defaultConfig)
	response, err := client.CreateChatCompletion(c, openai.ChatCompletionRequest{
		Model:    "deepseek-chat",
		Messages: []openai.ChatCompletionMessage{},
		Stream:   false,
	})
	if err != nil {
		responseXML(c, result.Error.Error(), msg.FromUserName, msg.ToUserName)
		return
	}

	h = db.History{
		UserID:  msg.FromUserName,
		Role:    response.Choices[0].Message.Role,
		Content: response.Choices[0].Message.Content,
	}
	result = db.DB.Create(&h)
	if result.Error != nil {
		responseXML(c, result.Error.Error(), msg.FromUserName, msg.ToUserName)
		return
	}
	responseXML(c, h.Content, msg.FromUserName, msg.ToUserName)
}

func responseXML(c *gin.Context, content, toUserName, fromUserName string) {
	msg := Message{
		ToUserName:   toUserName,
		FromUserName: fromUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      content,
	}
	c.XML(http.StatusOK, msg)
}
