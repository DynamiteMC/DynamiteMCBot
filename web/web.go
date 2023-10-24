package web

import (
	"encoding/json"
	"fmt"
	"gobot/config"
	"strings"

	"github.com/aimjel/minecraft/chat"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/gorilla/websocket"
)

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type Log struct {
	Type    string `json:"type"`
	Time    string `json:"time"`
	Message string `json:"message"`
}

func Connect(client bot.Client, config config.ConfigS) {
	//p, err := url.Parse(config.DynamiteServer)
	//if err != nil {
	//	return
	//}
	ws, _, err := websocket.DefaultDialer.Dial(config.DynamiteServer, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	ws.WriteJSON(Message{
		Type: "auth",
		Data: config.DynamitePassword,
	})
	var msgid snowflake.ID
	var msg string
	for {
		var m Message
		err := ws.ReadJSON(&m)
		if err != nil {
			fmt.Println(err)
			return
		}
		switch m.Type {
		case "sync":
			{
				data := m.Data.(map[string]interface{})
				logs := parseLogs(data["log"].(string))
				for _, log := range logs {
					switch log.Type {
					case "info":
						msg += fmt.Sprintf("\u001b[0;30m%s \u001b[0;34mINFO\u001b[37;0m: %s\n", log.Time, log.Message)
					case "debug":
						msg += fmt.Sprintf("\u001b[0;30m%s \u001b[0;36mDEBUG\u001b[37;0m: %s\n", log.Time, log.Message)
					case "warn":
						msg += fmt.Sprintf("\u001b[0;30m%s \u001b[0;33mDEBUG\u001b[37;0m: %s\n", log.Time, log.Message)
					case "error":
						msg += fmt.Sprintf("\u001b[0;30m%s \u001b[0;31mDEBUG\u001b[37;0m: %s\n", log.Time, log.Message)
					case "chat":
						var ms chat.Message
						json.Unmarshal([]byte(log.Message), &ms)
						msg += parseChat(ms) + "\n"
					}
				}
				m, err := client.Rest().CreateMessage(snowflake.ID(config.DynamiteLogChannel), discord.MessageCreate{
					Content: fmt.Sprintf("```ansi\n%s\n```", msg),
				})
				if err != nil {
					return
				}
				msgid = m.ID
			}
		case "auth":
			return
		case "log":
			var log Log
			json.Unmarshal([]byte(m.Data.(string)), &log)
			var typ string

			switch log.Type {
			case "info":
				msg += fmt.Sprintf("\u001b[0;30m%s \u001b[0;34mINFO\u001b[37;0m: %s\n", log.Time, log.Message)
			case "debug":
				msg += fmt.Sprintf("\u001b[0;30m%s \u001b[0;36mDEBUG\u001b[37;0m: %s\n", log.Time, log.Message)
			case "warn":
				msg += fmt.Sprintf("\u001b[0;30m%s \u001b[0;33mDEBUG\u001b[37;0m: %s\n", log.Time, log.Message)
			case "error":
				msg += fmt.Sprintf("\u001b[0;30m%s \u001b[0;31mDEBUG\u001b[37;0m: %s\n", log.Time, log.Message)
			case "chat":
				var ms chat.Message
				json.Unmarshal([]byte(log.Message), &ms)
				msg += parseChat(ms) + "\n"
			}

			msg += fmt.Sprintf("\u001b[0;30m%s %s\u001b[37;0m: %s\n", log.Time, typ, log.Message)
			d := fmt.Sprintf("```ansi\n%s\n```", msg)
			client.Rest().UpdateMessage(snowflake.ID(config.DynamiteLogChannel), msgid, discord.MessageUpdate{
				Content: &d,
			})
		}
	}
}

func parseLogs(l string) (logs []Log) {
	lo := strings.Split(l, "\n")
	for _, l := range lo {
		var log Log
		json.Unmarshal([]byte(l), &log)
		logs = append(logs, log)
	}
	return
}

var colors = map[string]string{
	"black":        "30",
	"dark_blue":    "34",
	"dark_green":   "32",
	"dark_aqua":    "36",
	"dark_red":     "31",
	"dark_purple":  "35",
	"gold":         "33",
	"gray":         "30",
	"dark_gray":    "30",
	"blue":         "34",
	"green":        "32",
	"aqua":         "36",
	"red":          "31",
	"light_purple": "35",
	"yellow":       "33",
	"white":        "37",
}

func parseChat(msg chat.Message) string {
	var str string
	texts := []chat.Message{msg}
	texts = append(texts, msg.Extra...)

	for _, text := range texts {
		s := fmt.Sprintf("\u001b[0;%sm", colors[text.Color])
		if text.Bold {
			s = fmt.Sprintf("\u001b[1;%sm", colors[text.Color])
		}
		if text.Underlined {
			s = fmt.Sprintf("\u001b[4;%sm", colors[text.Color])
		}
		str += s + text.Text
	}

	return str
}
