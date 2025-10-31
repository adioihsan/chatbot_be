package cmd

import (
	"log"

	"cms-octo-chat-api/app"
	"cms-octo-chat-api/config"
	"cms-octo-chat-api/model"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Webserver CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		c := config.Initiate("api")

		c.DB.AutoMigrate(&model.User{}, &model.UserMatrix{}, &model.Conversation{}, &model.Message{})

		router := app.ApiRouter(model.Resources{
			Env:    c.Env,
			Logs:   c.Logs,
			DB:     c.DB,
			OpenAI: c.OpenAI,
		})

		log.Fatal(router.Listen(":" + c.Env.AppApiPort))
	},
}
