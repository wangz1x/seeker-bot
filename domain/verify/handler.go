// verify - 2024/12/16
// Author: wangzx
// Description:

package verify

import (
	"github.com/gin-gonic/gin"
)

func Verify(c *gin.Context) {
	echoStr := c.Query("echostr")
	c.String(200, echoStr)
}
