package commands

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"

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
		if _, e := url.ParseRequestURI(site); e != nil {
			CreateMessage(message, Message{
				Content: "Invalid URL!",
				Reply:   true,
			})
		}
		///Applications/Google Chrome.app/Contents/MacOS/Google Chrome
		exec.Command("chrome", "--headless", "--disable-gpu", "--screenshot", "--window-size=1920,1080", site).Output()
		f, err := os.Open("screenshot.png")
		if err != nil {
			CreateMessage(message, Message{
				Content: fmt.Sprintf("Failed to screenshot site: %s", err),
				Reply:   true,
			})
			return
		}
		CreateMessage(message, Message{
			Files: []*discord.File{
				{
					Name:   "screenshot.png",
					Reader: f,
				},
			},
			Reply: true,
		})
		f.Close()
		os.Remove("screenshot.png")
	},
}
