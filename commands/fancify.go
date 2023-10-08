package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/disgoorg/disgo/events"
)

var Command_fancify = Command{
	Name:        "fancify",
	Description: "Fancify your sentence",
	Execute: func(message *events.MessageCreate, args []string) {
		s := strings.Join(args, " ")
		if s == "" {
			return
		}
		message.Client().Rest().SendTyping(message.ChannelID)
		ctx := context.Background()
		prompt := fmt.Sprintf("Can you please make this sentence in a more fancier language: \"%s\"\n\nPlease only send back the fancier sentence without commas or anything else", s)
		res, err := openaiClient.SimpleSend(ctx, prompt)
		if err != nil {
			CreateMessage(message, err.Error(), true)
		}
		CreateMessage(message, res.Choices[0].Message.Content, true)
	},
}
