package jobs

import (
	"github.com/kondohiroki/go-boilerplate/internal/app"
	"go.uber.org/zap"
)

func SendNotificationViaDiscord(c *app.AppContext) {
	// Send a message to the channel
	_, err := c.Discord.ChannelMessageSend(c.Config.Discord.ChannelID, "Hello, World!")
	if err != nil {
		c.Logger.Error("Error sending message: ", zap.Error(err))
		return
	}

	c.Logger.Info("Message sent successfully.") // Use logger from app context
	// fmt.Println("Message sent successfully.")
}
