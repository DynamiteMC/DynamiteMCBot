package commands

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"golang.org/x/net/html"
)

var Command_mcdoc = Command{
	Name:        "mcdoc",
	Description: "Searches wiki.vg",
	Aliases:     []string{"wiki.vg", "wikivg", "wvg"},
	Execute: func(message *events.MessageCreate, args []string) {
		if len(args) == 0 {
			return
		}
		typ := "-java"
		var search string
		if i, d := IsAny(args[len(args)-1], "-java", "-bedrock", "-pocket", "-classic"); i {
			typ = d
			args = args[:len(args)-1]
		}
		search = strings.Join(args, " ")
		var title string
		var url string
		var description string
		var keys []string
		var values []string
		var fields = [][]string{{"", "", ""}}

		switch typ {
		case "-bedrock":
			{

			}
		case "-pocket":
			{

			}
		case "-java":
			{
				resp, err := http.Get("https://wiki.vg/Protocol")
				if err != nil {
					CreateMessage(message, Message{
						Reply:   true,
						Content: fmt.Sprintf("Couldn't reach wiki.vg: %s", err),
					})
				}
				query, _ := goquery.NewDocumentFromReader(resp.Body)
				var q *html.Node
				for _, sel := range query.Find("div").Nodes {
					if len(sel.Attr) >= 1 {
						if sel.Attr[0].Val == "mw-parser-output" {
							q = sel
							break
						}
					}
				}

				for sel := q.FirstChild; sel != nil; sel = sel.NextSibling {
					if sel.Data == "h4" {
						if sel.FirstChild.FirstChild != nil && strings.EqualFold(sel.FirstChild.FirstChild.Data, search) {
							title = sel.FirstChild.FirstChild.Data
							for _, attr := range sel.FirstChild.Attr {
								if attr.Key == "id" {
									url = fmt.Sprintf("https://wiki.vg/Protocol#%s", attr.Val)
								}
							}
							if sel.NextSibling.NextSibling.Data == "p" {
								for a := sel.NextSibling.NextSibling.FirstChild; a != nil; a = a.NextSibling {
									if a.Data == "a" {
										var href string
										for _, attr := range a.Attr {
											if attr.Key == "href" {
												href = attr.Val
											}
										}
										if h, _ := HasAnyPrefix(href, "https://", "http://"); !h {
											href = "https://wiki.vg/Protocol" + href
										}
										description += fmt.Sprintf("[%s](%s)", a.FirstChild.Data, href)
									} else if a.Data == "kbd" {
										description += fmt.Sprintf("`%s`", a.FirstChild.Data)
									} else {
										description += a.Data
									}
								}
							}
							for i := sel.NextSibling; i != nil; i = i.NextSibling {
								if i.Data == "table" {
									for a := i.FirstChild; a != nil; a = a.NextSibling {
										if a.Data == "tbody" {
											for k := a.FirstChild.FirstChild; k != nil; k = k.NextSibling {
												if k.Data == "th" {
													keys = append(keys, k.FirstChild.Data)
												}
											}
											for k := a.FirstChild.NextSibling.NextSibling.FirstChild; k != nil; k = k.NextSibling {
												if k.Data == "td" {
													str := strings.TrimSpace(k.FirstChild.Data)
													if str == "a" {
														href := k.FirstChild.Attr[0].Val
														if h, _ := HasAnyPrefix(href, "https://", "http://"); !h {
															href = "https://wiki.vg/Protocol" + href
														}
														str = fmt.Sprintf("[%s](%s)", k.FirstChild.FirstChild.Data, href)
													}
													values = append(values, str)
												}
											}
											for k := a.FirstChild.NextSibling.NextSibling.NextSibling; k != nil; k = k.NextSibling {
												if k.Data == "tr" {
													var data []string
													for a := k.FirstChild; a != nil; a = a.NextSibling {
														if a.Data == "td" {
															str := strings.TrimSpace(a.FirstChild.Data)
															if str == "a" {
																href := a.FirstChild.Attr[0].Val
																if h, _ := HasAnyPrefix(href, "https://", "http://"); !h {
																	href = "https://wiki.vg/Protocol" + href
																}
																str = fmt.Sprintf("[%s](%s)", a.FirstChild.FirstChild.Data, href)
															}
															data = append(data, str)
														}
													}
													fields = append(fields, data)
												}
											}
										}
									}
									break
								}
							}
							break
						}
					}
				}
			}
		case "-classic":
			{

			}
		}
		if title == "" {
			CreateMessage(message, Message{
				Content: "Couldn't find packet!",
				Reply:   true,
			})
			return
		}

		var packetId string
		var state string
		var boundTo string
		for i, key := range keys {
			key = strings.TrimSpace(key)
			switch key {
			case "Packet ID":
				packetId = values[i]
				values[i] = ""
			case "State":
				state = values[i]
				values[i] = ""
			case "Bound To":
				boundTo = values[i]
				values[i] = ""
			case "Field Name":
				if values[i] == "i" {
					continue
				}
				fields[0] = []string{values[i], values[i+1], values[i+2]}
			}
		}
		description += "\n**Fields**:\n"
		for _, field := range fields {
			description += fmt.Sprintf("%s | %s | %s\n\n", field[0], field[1], field[2])
		}

		embed := discord.NewEmbedBuilder().
			SetTitle(title).
			SetURL(url).
			SetDescription(description).
			SetColor(color)
		embed.AddFields(
			discord.EmbedField{
				Name:   "Packet ID",
				Value:  packetId,
				Inline: &True,
			},
			discord.EmbedField{
				Name:   "State",
				Value:  state,
				Inline: &True,
			},
			discord.EmbedField{
				Name:   "Bound To",
				Value:  boundTo,
				Inline: &True,
			},
		)
		CreateMessage(message, Message{
			Embeds: []discord.Embed{embed.Build()},
			Reply:  true,
		})
	},
}
