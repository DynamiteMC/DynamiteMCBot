package commands

import (
	"fmt"
	"gobot/config"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

var Command_corner = Command{
	Name:        "corner",
	Description: "Sends a member to go sit in the corner",
	Permissions: discord.PermissionManageRoles,
	Aliases:     []string{"shame", "bully", "blind", "badboy", "lmfao", "didntask", "bruh", "ratio", "cope", "skillissue"},
	Execute: func(message *events.MessageCreate, args []string) {
		memberId := args[0]
		if memberId == "" {
			return
		}
		id := ParseMention(memberId)
		if id == 0 {
			CreateMessage(message, Message{Content: "Failed to parse member", Reply: true})
			return
		}
		err := message.Client().Rest().AddMemberRole(*message.GuildID, id, snowflake.ID(config.Config.MuteRole))
		if err != nil {
			CreateMessage(message, Message{Content: "Failed to mute member", Reply: true})
			return
		}
		var tag string
		member, err := message.Client().Rest().GetMember(*message.GuildID, id)
		if err != nil {
			tag = "Unknown#0000"
		} else {
			tag = member.User.Tag()
		}
		CreateMessage(message, Message{Content: fmt.Sprintf("Sent member %s to go sit in the corner.", tag), Reply: true})
	},
}
