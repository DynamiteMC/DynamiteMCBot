package commands

import (
	"gobot/config"
	"strconv"
	"strings"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

var color = 0x9C182C
var startTime = time.Now()
var True = true

func Point[T any](data T) *T {
	return &data
}

func HasAnyPrefix(str string, prefixes ...string) (bool, string) {
	for _, s := range prefixes {
		if strings.HasPrefix(str, s) {
			return true, s
		}
	}
	return false, ""
}

func IsAny(str string, strs ...string) (bool, string) {
	for _, s := range strs {
		if str == s {
			return true, s
		}
	}
	return false, ""
}

func GetArgument(args []string, index int) string {
	if len(args) <= index {
		return ""
	}
	return args[index]
}

type Command struct {
	Name        string
	Description string
	Execute     func(*events.MessageCreate, []string)
	Aliases     []string
	Permissions discord.Permissions
}

type Message struct {
	Content string
	Reply   bool
	Embeds  []discord.Embed
	Files   []*discord.File
}

func EditMessage(client bot.Client, channelID snowflake.ID, id snowflake.ID, message Message) (*discord.Message, error) {
	builder := discord.NewMessageUpdateBuilder().
		SetContent(message.Content).
		SetEmbeds(message.Embeds...).
		SetFiles(message.Files...).
		SetAllowedMentions(&discord.AllowedMentions{RepliedUser: false})
	return client.Rest().UpdateMessage(channelID, id, builder.Build())
}

func CreateMessage(e *events.MessageCreate, message Message) (*discord.Message, error) {
	builder := discord.NewMessageCreateBuilder().
		SetContent(message.Content).
		SetEmbeds(message.Embeds...).
		SetFiles(message.Files...).
		SetAllowedMentions(&discord.AllowedMentions{RepliedUser: false})
	if message.Reply {
		builder.SetMessageReferenceByID(e.MessageID)
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
	if message.Message.Author.Bot {
		return
	}
	if strings.Contains(strings.ToLower(message.Message.Content), "mmm") {
		message.Client().Rest().AddReaction(message.ChannelID, message.MessageID, "✅")
	}
	args := strings.Split(message.Message.Content, " ")
	if !strings.HasPrefix(message.Message.Content, config.Config.InfoPrefix) {
		if strings.HasPrefix(message.Message.Content, config.Config.Prefix) {
			cmd := strings.Split(args[0], "\n")[0][len(config.Config.Prefix):]
			if len(strings.Split(args[0], "\n")) == 1 {
				args = args[1:]
			} else {
				args[0] = strings.Join(strings.Split(args[0], "\n")[1:], "\n")
			}
			command := commands[cmd]
			if command.Execute == nil {
				command = commands[aliases[cmd]]
				if command.Execute == nil {
					return
				}
			}
			message.Message.Member.GuildID = *message.GuildID
			if !message.Client().Caches().MemberPermissions(*message.Message.Member).Has(command.Permissions) {
				return
			}
			command.Execute(message, args)
		} else {
			return
		}
	} else {
		cmd := args[0][len(config.Config.InfoPrefix):]
		command := commands[cmd]
		if command.Execute == nil {
			command = commands[aliases[cmd]]
			if command.Execute == nil {
				return
			}
		}
		aliases := "None"
		if len(command.Aliases) > 0 {
			aliases = strings.Join(command.Aliases, ", ")
		}
		embed := discord.NewEmbedBuilder().
			SetTitle(command.Name).
			SetDescription(command.Description).
			AddFields(
				discord.EmbedField{
					Name:   "Aliases",
					Value:  aliases,
					Inline: &True,
				},
				discord.EmbedField{
					Name:   "Permissions Required",
					Value:  command.Permissions.String(),
					Inline: &True,
				},
			).SetColor(color).Build()
		CreateMessage(message, Message{Embeds: []discord.Embed{embed}, Reply: true})
	}
}

var commands = make(map[string]Command)
var aliases = make(map[string]string)

func RegisterCommands(cmds ...Command) {
	for _, command := range cmds {
		commands[command.Name] = command
		for _, alias := range command.Aliases {
			aliases[alias] = command.Name
		}
	}
}
