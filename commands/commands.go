package commands

import (
	"gobot/config"
	"strconv"
	"strings"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

type Command struct {
	Name        string
	Description string
	Execute     func(*events.MessageCreate, []string)
	Aliases     []string
	Permissions discord.Permissions
}

func FetchAlias(alias string) (bool, Command) {
	for _, command := range commands {
		for _, a := range command.Aliases {
			if a == alias {
				return true, command
			}
		}
	}
	return false, Command{}
}

func CreateMessage(e *events.MessageCreate, content string, reply bool) (*discord.Message, error) {
	builder := discord.NewMessageCreateBuilder().SetContent(content)
	if reply {
		builder.SetMessageReferenceByID(e.MessageID).SetAllowedMentions(&discord.AllowedMentions{RepliedUser: false})
	}
	return e.Client().Rest().CreateMessage(e.ChannelID, builder.Build())
}

func HasRole(client bot.Client, guildId snowflake.ID, memberId snowflake.ID, id int64) bool {
	member, _ := client.Rest().GetMember(guildId, memberId)
	for _, role := range member.RoleIDs {
		if role == snowflake.ID(id) {
			return true
		}
	}
	return false
}

func ParseMention(mention string) snowflake.ID {
	if strings.HasPrefix(mention, "<@") && strings.HasSuffix(mention, ">") && !strings.HasPrefix(mention, "<@&") {
		mention = strings.TrimPrefix(strings.TrimSuffix(mention, ">"), "<@")
	}
	id, err := strconv.ParseInt(mention, 10, 64)
	if err != nil {
		return 0
	}
	return snowflake.ID(id)
}

func Handle(message *events.MessageCreate) {
	if message.Message.Author.Bot || !strings.HasPrefix(message.Message.Content, config.Config.Prefix) {
		return
	}
	args := strings.Split(message.Message.Content, " ")
	cmd := args[0][1:]
	args = args[1:]
	command := commands[cmd]
	if command.Execute == nil {
		var exists bool
		exists, command = FetchAlias(cmd)
		if !exists || command.Execute == nil {
			return
		}
	}
	message.Message.Member.GuildID = *message.GuildID
	if !message.Client().Caches().MemberPermissions(*message.Message.Member).Has(command.Permissions) {
		return
	}
	command.Execute(message, args)
}

var commands = map[string]Command{"oq": oq, "mute": mute, "unmute": unmute}
