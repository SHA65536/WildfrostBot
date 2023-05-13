package wildfrostbot

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// SearchHandler is a registration for the search slash command
var SearchHandler Cmd = Cmd{
	Cmd: &discordgo.ApplicationCommand{
		Name:        "search",
		Description: "Search the wiki!",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:         discordgo.ApplicationCommandOptionString,
				Name:         "query",
				Description:  "Query",
				Required:     true,
				Autocomplete: true,
			},
		},
	},
	Handler: SearchCmd,
}

// SearchCmd handles the search slash command and button press
func SearchCmd(h *DiscordHandler, s *discordgo.Session, i *discordgo.InteractionCreate) {
	var query string

	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		// If command
		query = i.ApplicationCommandData().Options[0].Value.(string)
		log.Printf("[DISCORD] Command search invoked with \"%s\"", query)
	case discordgo.InteractionMessageComponent:
		// If button
		query = strings.Split(i.MessageComponentData().CustomID, "_")[1]
		log.Printf("[DISCORD] Button search invoked with \"%s\"", query)
	}

	// If embed found
	if embed, ok := h.Library.Articles[query]; ok {
		// Generate buttons from tags
		var rows []discordgo.MessageComponent
		for _, v := range embed.MakeMeta().Tags {
			rows = append(rows, discordgo.Button{
				Emoji:    discordgo.ComponentEmoji{Name: "🔎"},
				Label:    strings.Title(v),
				Style:    discordgo.SuccessButton,
				CustomID: "search_" + v,
			})
		}
		// Only send componenets if there are any
		if len(rows) > 0 {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						embed.MakeEmbed(),
					},
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: rows,
						},
					},
				},
			})
		} else {
			SendEmbed(s, i, embed.MakeEmbed())
		}

	} else {
		SendEmbed(s, i, FailEmbed)
	}
}

// SearchAutocomplete handles autocomplete interactions for the search command
func SearchAutocomplete(h *DiscordHandler, s *discordgo.Session, i *discordgo.InteractionCreate) {
	var options []*discordgo.ApplicationCommandOptionChoice
	// Case insensitive
	query := strings.ToLower(i.ApplicationCommandData().Options[0].StringValue())
	// Find up to 25 results with same prefix
	for i := 0; i < len(h.Library.Commands) && len(options) < 25; i++ {
		if strings.HasPrefix(h.Library.Commands[i], query) {
			options = append(options, &discordgo.ApplicationCommandOptionChoice{
				Name:  strings.Title(h.Library.Commands[i]),
				Value: h.Library.Commands[i],
			})
		}
	}
	// Send autocomplete
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: options,
		},
	})
}
