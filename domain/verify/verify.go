// verify - 2024/12/16
// Author: wangzx
// Description:

package verify

import (
	"crypto/sha1"
	"fmt"
	"github.com/gin-gonic/gin"
	"seeker-bot/m/conf"
	"sort"
	"strings"
)

var token = conf.GvaConfig.App.Token

func Verify(c *gin.Context) {
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	echoStr := c.Query("echostr")
	// 验证签名
	if !verifySignature(signature, timestamp, nonce) {
		c.String(403, "Invalid signature")
		return
	}
	c.String(200, echoStr)
}

func verifySignature(signature, timestamp, nonce string) bool {
	// 1. 将token、timestamp、nonce三个参数放入切片
	strs := []string{token, timestamp, nonce}

	// 2. 字典序排序
	sort.Strings(strs)

	// 3. 将三个参数拼接成一个字符串
	str := strings.Join(strs, "")

	// 4. 进行sha1加密
	h := sha1.New()
	h.Write([]byte(str))
	sha1Sum := fmt.Sprintf("%x", h.Sum(nil))

	// 5. 将加密后的字符串与signature进行对比
	return sha1Sum == signature
}
