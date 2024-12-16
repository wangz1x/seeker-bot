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
	"seeker-bot/m/db"
)

const baseURL = "https://api.deepseek.com/chat/completions"
const token = "sk-43675ee878c74c6d91ad5ba4500fcdb1"

func Chat(c *gin.Context) {

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	defer c.Request.Body.Close()

	var msg Message
	err = xml.Unmarshal(body, &msg)
	if err != nil {
		c.String(400, err.Error())
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
		c.String(500, result.Error.Error())
		return
	}

	var records []db.History
	result = db.DB.Where("user_id = ?", msg.FromUserName).Order("created_at desc").Limit(10).Find(&records)
	if result.Error != nil {
		c.String(500, result.Error.Error())
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
		c.String(500, err.Error())
		return
	}

	h = db.History{
		UserID:  msg.FromUserName,
		Role:    response.Choices[0].Message.Role,
		Content: response.Choices[0].Message.Content,
	}
	result = db.DB.Create(&h)
	if result.Error != nil {
		c.String(500, result.Error.Error())
		return
	}
	c.String(200, h.Content)
}
