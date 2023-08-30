package commands

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var Command_screenshot = Command{
	Name:        "screenshot",
	Description: "Screenshot a website",
	Permissions: discord.PermissionAdministrator,
	Aliases:     []string{"go", "goto"},
	Execute: func(message *events.MessageCreate, args []string) {
		site := GetArgument(args, 0)
		if site == "" {
			return
		}
		if !strings.HasPrefix(site, "https://") && !strings.HasPrefix(site, "http://") {
			site = "http://" + site
		}
		if _, e := url.ParseRequestURI(site); e != nil {
			CreateMessage(message, Message{
				Content: "Invalid URL!",
				Reply:   true,
			})
			return
		}
		msg, _ := CreateMessage(message, Message{
			Content: "Screenshotting...",
			Reply:   true,
		})
		exec.Command("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", "--headless", "--disable-gpu", "--screenshot", "--window-size=1366,768", site).Run()
		f, err := os.Open("screenshot.png")
		if err != nil {
			EditMessage(message.Client(), msg.ChannelID, msg.ID, Message{
				Content: fmt.Sprintf("Failed to screenshot site: %s", err),
			})
			return
		}
		EditMessage(message.Client(), msg.ChannelID, msg.ID, Message{
			Files: []*discord.File{
				{
					Name:   "screenshot.png",
					Reader: f,
				},
			},
		})
		f.Close()
		os.Remove("screenshot.png")
	},
}
