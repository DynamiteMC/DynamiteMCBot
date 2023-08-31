package commands

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func parseCodeBlock(data string) string {
	data = strings.TrimSpace(data)
	data = strings.TrimPrefix(data, "```")
	data = strings.TrimSuffix(data, "```")
	data = strings.TrimPrefix(data, "go")
	data = strings.TrimSpace(data)
	return data
}

func ease(src string) string {
	if !strings.HasPrefix(src, "package main") {
		src = "package main\n" + src
	}
	return src
}

var Command_go = Command{
	Name:        "go",
	Description: "Run Go code",
	Aliases:     []string{"golang"},
	Execute: func(message *events.MessageCreate, args []string) {
		code := ease(parseCodeBlock(strings.Join(args, " ")))
		ms, _ := CreateMessage(message, Message{
			Reply:   true,
			Content: "Running...",
		})
		data := struct {
			Body    string
			WithVet bool
		}{code, false}
		b, _ := json.Marshal(data)
		resp, _ := http.Post("https://play.golang.org/compile", "application/json", bytes.NewReader(b))
		d, _ := io.ReadAll(resp.Body)
		var resps map[string]interface{}
		json.Unmarshal(d, &resps)
		var stdout string
		var errors = resps["Errors"].(string)
		var stderr string
		if m, ok := resps["Events"].([]interface{}); ok {
			for _, msg := range m {
				if msg, ok := msg.(map[string]interface{}); ok {
					if msg["Kind"] == "stdout" {
						stdout += msg["Message"].(string)
					}
					if msg["Kind"] == "stderr" {
						stderr += msg["Message"].(string)
					}
				}
			}
		}
		embed := discord.NewEmbedBuilder().SetColor(color).SetTitle("Code")
		if errors != "" {
			embed.AddField("Compile Errors", "```"+errors+"```", false)
		}
		if stdout != "" {
			embed.AddField("Output", "```"+stdout+"```", false)
		}
		if stderr != "" {
			embed.AddField("Errors", "```"+stderr+"```", false)
		}
		EditMessage(message.Client(), message.ChannelID, ms.ID, Message{
			Embeds: []discord.Embed{embed.Build()},
		})
	},
}
