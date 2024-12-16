// handler - 2024/12/16
// Author: wangzx
// Description:

package chat

import (
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday/v2"
	"github.com/sashabaranov/go-openai"
	"html/template"
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

	slog.Info("receive message", "content", msg.Content, "user", msg.FromUserName, "msgId", msg.MsgId)
	var count int64
	result := db.DB.Model(&db.History{}).Where("user_id = ? and msg_id = ?", msg.FromUserName, msg.MsgId).Count(&count)
	if result.Error != nil {
		responseXML(c, result.Error.Error(), msg.FromUserName, msg.ToUserName)
		return
	}

	if count > 0 {
		// 说明是重试的，直接忽略，啥都不干
		return
	}

	h := db.History{
		UserID:  msg.FromUserName,
		Role:    "user",
		Content: msg.Content,
		MsgID:   msg.MsgId,
	}

	result = db.DB.Create(&h)
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

	resultURL := fmt.Sprintf("%s/chat/result/%s", "http://106.54.124.217:80", msg.MsgId)

	done := make(chan bool)

	var response openai.ChatCompletionResponse
	go func() {
		defaultConfig := openai.DefaultConfig(token)
		defaultConfig.BaseURL = baseURL
		client := openai.NewClientWithConfig(defaultConfig)
		response, err = client.CreateChatCompletion(c, openai.ChatCompletionRequest{
			Model:    "deepseek-chat",
			Messages: messages,
			Stream:   false,
		})
		if err == nil {
			h = db.History{
				UserID:  msg.FromUserName,
				Role:    response.Choices[0].Message.Role,
				Content: response.Choices[0].Message.Content,
				MsgID:   msg.MsgId,
			}
			_ = db.DB.Create(&h)
		}
		done <- true
	}()

	select {
	case <-done:
		responseXML(c, response.Choices[0].Message.Content, msg.FromUserName, msg.ToUserName)
	case <-time.After(time.Second * 3):
		responseXML(c, fmt.Sprintf("正在思考中，请稍后<a href='%s'>查看</a>", resultURL), msg.FromUserName, msg.ToUserName)
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

func Result(c *gin.Context) {
	msgId := c.Param("msgId")

	var history db.History
	result := db.DB.Where("msg_id = ? and role = ?", msgId, "assistant").First(&history)
	if result.Error != nil {
		// 如果还在处理中，显示等待页面
		c.HTML(http.StatusOK, "waiting.html", gin.H{})
		return
	}

	// 显示结果页面
	answerHTML := template.HTML(blackfriday.Run([]byte(history.Content)))
	c.HTML(http.StatusOK, "result.html", gin.H{
		"answer": answerHTMLg,
	})
}
