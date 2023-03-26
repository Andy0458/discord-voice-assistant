package commands

import (
	"fmt"

	error "discord_voice_assistant/error"
	discordgo "github.com/bwmarrin/discordgo"
)

var Add discordgo.ApplicationCommand = discordgo.ApplicationCommand{
	Name: "add",
	Description: "Adds 2 numbers.",
	Type: discordgo.ChatApplicationCommand,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name: "first",
			Description: "First number to add.",
			Type: discordgo.ApplicationCommandOptionNumber,
			Required: true,
		},
		{
			Name: "second",
			Description: "Second number to add.",
			Type: discordgo.ApplicationCommandOptionNumber,
			Required: true,
		},
	},
}

func AddHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	a := optionMap["first"].FloatValue()
	b := optionMap["second"].FloatValue()
	// ^ There's no need to check second bool value if option exists because we set them as required on lines 15 & 21.

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Result: %.2f", a+b),
		},
	})
	if err != nil {
		error.FormatError(err)
	}
}