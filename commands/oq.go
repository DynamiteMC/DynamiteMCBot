package commands

import (
	"github.com/disgoorg/disgo/events"
)

var Command_oq = Command{
	Name:        "oq",
	Description: "oq",
	Execute: func(message *events.MessageCreate, args []string) {
		CreateMessage(message, Message{Content: "oq", Reply: true})
	},
}
