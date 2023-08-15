package commands

import (
	"github.com/disgoorg/disgo/events"
)

var oq = Command{
	Name:        "oq",
	Description: "oq",
	Execute: func(message *events.MessageCreate, s []string) {
		CreateMessage(message, "oq")
	},
}
