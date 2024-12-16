// server - 2024/12/16
// Author: wangzx
// Description:

package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"log/slog"
	"seeker-bot/m/conf"
	"seeker-bot/m/domain/chat"
	"seeker-bot/m/domain/verify"
	"seeker-bot/m/middleware"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the server used by WeChat",
	Long:  `None`,

	Run: func(cmd *cobra.Command, args []string) {
		server()
	},
}

func server() {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(middleware.Auth)

	r.GET("/seeker", func(c *gin.Context) {
		verify.Verify(c)
	})

	r.POST("/seeker", func(c *gin.Context) {
		chat.Chat(c)
	})

	r.POST("/chat/result/:msgId", func(c *gin.Context) {
		chat.Result(c)
	})

	slog.Info("http server start")
	err := r.Run(conf.GvaConfig.App.Addr)
	if err != nil {
		return
	}
}
