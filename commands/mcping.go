package commands

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func ParseIP(ip string) bool {
	if net.ParseIP(ip) != nil || len(strings.Split(ip, ".")) >= 2 {
		return true
	} else {
		sp := strings.Split(ip, ":")
		if len(sp) == 2 {
			if net.ParseIP(sp[0]) != nil {
				if _, err := strconv.Atoi(sp[1]); err == nil {
					return true
				}
			}
		}
	}
	return false
}

var Command_mcping = Command{
	Name:        "mcping",
	Description: "Ping a minecraft server",
	Aliases:     []string{"mcpi"},
	Execute: func(message *events.MessageCreate, args []string) {
		typ := "-java"
		ip := GetArgument(args, 0)
		if i, d := IsAny(args[len(args)-1], "-java", "-bedrock"); i {
			typ = d
		}
		if !ParseIP(ip) {
			CreateMessage(message, Message{Content: "Invalid IP address!", Reply: true})
		}
		var server map[string]interface{}
		if typ == "-bedrock" {
			resp, _ := http.Get(fmt.Sprintf("https://api.mcstatus.io/v2/status/bedrock/%s", ip))
			dec := json.NewDecoder(resp.Body)
			dec.Decode(&server)
		} else {
			resp, _ := http.Get(fmt.Sprintf("https://api.mcstatus.io/v2/status/java/%s", ip))
			dec := json.NewDecoder(resp.Body)
			dec.Decode(&server)
		}
		if _, ok := server["motd"]; !ok {
			CreateMessage(message, "Unknown server.", true)
			return
		}
		embed := discord.NewEmbedBuilder().
			SetAuthorName(ip).
			SetDescription(server["motd"].(map[string]interface{})["clean"].(string)).
			SetColor(color)
		for key, value := range server {
			switch key {
			case "version":
				{
					if version, ok := value.(map[string]interface{}); ok {
						name, ok := version["name_clean"].(string)
						if !ok {
							name, ok = version["name"].(string)
						}
						if !ok {
							continue
						}
						embed.AddField("Version", name, true)
					}
				}
			case "players":
				{
					if players, ok := value.(map[string]interface{}); ok {
						embed.AddField("Players", fmt.Sprintf("%v/%v", players["online"], players["max"]), true)
					}
				}
			}
		}
		CreateMessage(message, embed.Build(), true)
	},
}
