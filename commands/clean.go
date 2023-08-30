package commands

import (
	"strconv"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

var Command_clean = Command{
	Name:        "clean",
	Description: "Clean a member",
	Permissions: discord.PermissionManageMessages,
	Aliases:     []string{"clear", "klean", "klear"},
	Execute: func(message *events.MessageCreate, args []string) {
		a := GetArgument(args, 0)
		if a == "" {
			return
		}
		amount, err := strconv.Atoi(a)
		if err != nil {
			return
		}
		if amount > 100 {
			CreateMessage(message, Message{
				Content: "You can only delete up to 100 messages at once",
				Reply:   true,
			})
			return
		}
		memberId := GetArgument(args, 1)
		id := ParseMention(memberId)
		var messages []snowflake.ID
		msgs, _ := message.Client().Rest().GetMessages(message.ChannelID, 0, message.MessageID, 0, amount)
		for _, m := range msgs {
			messages = append(messages, m.ID)
		}
		if id != 0 {
			messages = []snowflake.ID{}
			for _, m := range msgs {
				if m.Author.ID == id {
					messages = append(messages, m.ID)
				}
			}
		}
		message.Client().Rest().BulkDeleteMessages(message.ChannelID, messages)
		CreateMessage(message, Message{
			Content: "Successfully cleaned channel.",
			Reply:   true,
		})
	},
}
