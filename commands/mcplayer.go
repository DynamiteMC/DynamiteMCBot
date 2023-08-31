package commands

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/google/uuid"
)

func ParseUUID(uuid string) string {
	str := ""
	for i, char := range strings.Split(uuid, "") {
		str += char
		if i == 7 || i == 11 || i == 15 || i == 19 {
			str += "-"
		}
	}
	return str
}
func IsUUID(str string) bool {
	if _, e := uuid.Parse(str); e == nil {
		return true
	}
	if _, e := uuid.Parse(ParseUUID(str)); e == nil {
		return true
	}
	return false
}

var Command_mcplayer = Command{
	Name:        "mcplayer",
	Description: "Get info about a player",
	Aliases:     []string{"mcp"},
	Execute: func(message *events.MessageCreate, args []string) {
		id := GetArgument(args, 0)
		if id == "" {
			return
		}
		var profile map[string]interface{}
		var properties map[string]interface{}
		if IsUUID(id) {
			resp, _ := http.Get(fmt.Sprintf("https://sessionserver.mojang.com/session/minecraft/profile/%s", id))
			decoder := json.NewDecoder(resp.Body)
			decoder.Decode(&profile)
		} else {
			resp, _ := http.Get(fmt.Sprintf("https://api.mojang.com/users/profiles/minecraft/%s", id))
			var p map[string]interface{}
			decoder := json.NewDecoder(resp.Body)
			decoder.Decode(&p)

			if _, ok := p["errorMessage"]; ok {
				CreateMessage(message, Message{
					Content: "Unknown player.",
					Reply:   true,
				})
			}
			resp, _ = http.Get(fmt.Sprintf("https://sessionserver.mojang.com/session/minecraft/profile/%s", p["id"]))
			decoder = json.NewDecoder(resp.Body)
			decoder.Decode(&profile)
		}

		if _, ok := profile["errorMessage"]; ok {
			CreateMessage(message, Message{
				Content: "Unknown player.",
				Reply:   true,
			})
		}
		d, _ := base64.StdEncoding.DecodeString(profile["properties"].([]interface{})[0].(map[string]interface{})["value"].(string))
		json.Unmarshal(d, &properties)

		embed := discord.NewEmbedBuilder().
			SetAuthorName(profile["name"].(string)).
			SetAuthorIcon(fmt.Sprintf("https://crafatar.com/avatars/%s", profile["id"])).
			SetColor(color).
			AddFields(
				discord.EmbedField{
					Name:   "UUID",
					Value:  ParseUUID(profile["id"].(string)),
					Inline: &True,
				},
			)
		if s, ok := properties["textures"].(map[string]interface{})["SKIN"]; ok {
			embed.SetThumbnail(s.(map[string]interface{})["url"].(string))
		}

		CreateMessage(message, Message{
			Reply:  true,
			Embeds: []discord.Embed{embed.Build()},
		})
	},
}
