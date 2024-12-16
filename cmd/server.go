// server - 2024/12/16
// Author: wangzx
// Description:

package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"log/slog"
	"seeker-bot/m/conf"
	"seeker-bot/m/domain/verify"
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

	r.GET("/seeker", func(c *gin.Context) {
		verify.Verify(c)
	})

	slog.Info("http server start")
	err := r.Run(conf.GvaConfig.App.Addr)
	if err != nil {
		return
	}
}
