package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Discord struct {
	*discordgo.Session
}

type Color string

const (
	SuccessColor Color = "0x198754"
	FailedColor  Color = "0xDC3545"
	WarningColor Color = "0xFFC107"
	InfoColor    Color = "0x0D6EFD"
)

func NewDiscord(token string) *Discord {
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return nil
	}

	// Open the Discord session
	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
		return nil
	}

	return &Discord{discord}
}
