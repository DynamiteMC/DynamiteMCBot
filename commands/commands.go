package commands

import (
	"fmt"
	"gobot/config"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type Command struct {
	Name        string
	Description string
	Execute     func(*events.MessageCreate, []string)
	Permissions discord.Permissions
}

func CreateMessage(e *events.MessageCreate, content string) {
	e.Client().Rest().CreateMessage(e.ChannelID, discord.NewMessageCreateBuilder().SetContent(content).Build())
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
		return
	}
	fmt.Println(message.Client().Caches().MemberPermissions(*message.Message.Member))
	if !message.Client().Caches().MemberPermissions(*message.Message.Member).Has(command.Permissions) {
		return
	}
	command.Execute(message, args)
}

var commands = map[string]Command{"oq": oq}
