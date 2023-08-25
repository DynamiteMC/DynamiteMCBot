package commands

import (
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var Command_ban = Command{
	Name:        "ban",
	Description: "Ban a member",
	Permissions: discord.PermissionBanMembers,
	Aliases:     []string{"banish", "bab", "bam"},
	Execute: func(message *events.MessageCreate, args []string) {
		memberId := GetArgument(args, 0)
		durationString := GetArgument(args, 1)
		reason := GetArgument(args, 2)
		var duration time.Duration
		if memberId == "" {
			return
		}
		if durationString == "" {
			duration = time.Duration(0)
		} else {
			d, err := time.ParseDuration(durationString)
			if err != nil {
				duration = time.Duration(0)
			} else {
				duration = d
			}
		}
		id := ParseMention(memberId)
		if id == 0 {
			CreateMessage(message, Message{Content: "Failed to parse member", Reply: true})
			return
		}
		member, err := message.Client().Rest().GetMember(*message.GuildID, id)
		var tag string
		if err != nil {
			tag = "Unknown#0000"
		} else {
			tag = member.User.Tag()
		}
		err = message.Client().Rest().AddBan(*message.GuildID, id, duration)
		if err != nil {
			CreateMessage(message, Message{Content: "Failed to ban member", Reply: true})
			return
		}
		msg := fmt.Sprintf("Banished member %s.", tag)
		if reason != "" {
			msg += fmt.Sprintf("\nReason: `%s`", reason)
		}
		CreateMessage(message, Message{Content: msg, Reply: true})
	},
}
