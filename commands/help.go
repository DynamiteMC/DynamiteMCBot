package commands

import (
	"fmt"
	"gobot/config"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var Command_help = Command{
	Name:        "help",
	Description: "Show a list of commands",
	Execute: func(message *events.MessageCreate, args []string) {
		embed := discord.NewEmbedBuilder().
			SetTitle("Commands").
			SetFooterText(fmt.Sprintf("Use %s<command> to get info about a command\n", config.Config.InfoPrefix)).
			SetColor(color)
		for _, command := range commands {
			embed.SetDescription(embed.Description + fmt.Sprintf("%s**%s** - %s\n", config.Config.Prefix, command.Name, command.Description))
		}
		CreateMessage(message, Message{Embeds: []discord.Embed{embed.Build()}, Reply: true})
	},
}
