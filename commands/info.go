package commands

import (
	"fmt"
	"runtime"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var Command_info = Command{
	Name:        "info",
	Description: "Shows info about the bot",
	Execute: func(message *events.MessageCreate, args []string) {
		var stats runtime.MemStats
		runtime.ReadMemStats(&stats)
		uptime := time.Since(startTime)
		embed := discord.NewEmbedBuilder().
			SetTitle("GoBot").
			SetURL("https://github.com/DynamiteMC/GoBot").
			SetColor(color).
			AddFields(
				discord.EmbedField{
					Name:   "Developer",
					Inline: &True,
					Value:  "oq_x (<@755487116615745607>)",
				},
				discord.EmbedField{
					Name:   "Uptime",
					Inline: &True,
					Value:  uptime.String(),
				},
				discord.EmbedField{
					Name:   "Ram Usage",
					Value:  fmt.Sprintf("%dMiB", stats.HeapInuse/1024/1024),
					Inline: &True,
				})
		CreateMessage(message, Message{
			Reply:  true,
			Embeds: []discord.Embed{embed.Build()},
		})
	},
}
