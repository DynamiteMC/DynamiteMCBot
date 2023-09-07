package commands

import (
	"gobot/store"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var Command_fetchdb = Command{
	Name:        "fetchdb",
	Description: "",
	Permissions: discord.PermissionAdministrator,
	Aliases:     []string{},
	Execute: func(message *events.MessageCreate, args []string) {
		db := GetArgument(args, 0)
		var ids []string
		if db == "" {
			return
		}
		switch db {
		case "corner":
			for i := range store.GetCorners() {
				ids = append(ids, i)
			}
		case "muted":
			for i := range store.GetMuted() {
				ids = append(ids, i)
			}
		}
		embed := discord.NewEmbedBuilder().SetTitlef("%d entries", len(ids)).SetDescription(strings.Join(ids, "\n")).SetColor(color)
		CreateMessage(message, embed.Build(), true)
	},
}
