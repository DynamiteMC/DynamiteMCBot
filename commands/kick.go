package commands

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var Command_kick = Command{
	Name:        "kick",
	Description: "Kick a member",
	Permissions: discord.PermissionKickMembers,
	Aliases:     []string{"yeet", "kicc"},
	Execute: func(message *events.MessageCreate, args []string) {
		memberId := GetArgument(args, 0)
		if memberId == "" {
			return
		}
		id := ParseMention(memberId)
		if id == 0 {
			CreateMessage(message, Message{Content: "Failed to parse member", Reply: true})
			return
		}
		member, err := message.Client().Rest().GetMember(*message.GuildID, id)
		if err != nil {
			CreateMessage(message, Message{Content: "Member is not in the server.", Reply: true})
			return
		}
		tag := member.User.Tag()
		err = message.Client().Rest().RemoveMember(*message.GuildID, id)
		if err != nil {
			CreateMessage(message, Message{Content: "Failed to kick member", Reply: true})
			return
		}
		CreateMessage(message, Message{Content: fmt.Sprintf("Yeeted member %s.", tag), Reply: true})
	},
}
