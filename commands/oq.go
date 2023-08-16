package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var oq = Command{
	Name:        "oq",
	Description: "oq",
	Permissions: discord.PermissionMuteMembers,
	Execute: func(message *events.MessageCreate, s []string) {
		CreateMessage(message, "oq")
	},
}
