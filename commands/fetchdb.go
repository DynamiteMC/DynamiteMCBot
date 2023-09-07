package commands

import (
	"fmt"
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
		if db == "" {
			return
		}
		if db == "-u" {
			db = GetArgument(args, 1)
			if db == "" {
				return
			}
			user := GetArgument(args, 2)
			if user == "" {
				return
			}
			id := ParseMention(user)
			if id == 0 {
				CreateMessage(message, Message{Content: "Failed to parse member", Reply: true})
				return
			}
			switch db {
			case "muted":
				CreateMessage(message, fmt.Sprint(store.IsMuted(int64(id))), true)
			case "corner":
				{
					is, corner := store.GetCorner(int64(id))
					if !is {
						CreateMessage(message, "false", true)
						return
					}
					CreateMessage(message, fmt.Sprint(corner), true)
				}
			}
		} else {
			var ids []string
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
		}
	},
}
