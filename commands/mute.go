package commands

import (
	"fmt"
	"gobot/config"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

var mute = Command{
	Name:        "mute",
	Description: "Mute a member",
	Permissions: discord.PermissionMuteMembers,
	Aliases:     []string{"silence", "shush"},
	Execute: func(message *events.MessageCreate, args []string) {
		memberId := args[0]
		if memberId == "" {
			return
		}
		id := ParseMention(memberId)
		if id == 0 {
			CreateMessage(message, "Failed to parse member", true)
			return
		}
		if HasRole(message.Client(), *message.GuildID, id, config.Config.MuteRole) {
			CreateMessage(message, "Member is already silenced.", true)
			return
		}
		err := message.Client().Rest().AddMemberRole(*message.GuildID, id, snowflake.ID(config.Config.MuteRole))
		if err != nil {
			CreateMessage(message, "Failed to mute member", true)
			return
		}
		var tag string
		member, err := message.Client().Rest().GetMember(*message.GuildID, id)
		if err != nil {
			tag = "Unknown#0000"
		}
		tag = member.User.Tag()
		CreateMessage(message, fmt.Sprintf("Silenced member %s.", tag), true)
	},
}
