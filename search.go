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
	var isLong bool

	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		// If command
		query = i.ApplicationCommandData().Options[0].Value.(string)
		log.Printf("[DISCORD] Command search invoked with \"%s\"", query)
	case discordgo.InteractionMessageComponent:
		// If button
		split := strings.Split(i.MessageComponentData().CustomID, "_")
		isLong = split[0] == "long"
		query = split[1]
		log.Printf("[DISCORD] Button search invoked with \"%s\"", query)
	}

	if isLong {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Flags: discordgo.MessageFlagsEphemeral},
		})
	} else {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})
	}

	// If embed found
	if embed, ok := h.Library.Articles[query]; ok {
		var message = &discordgo.WebhookParams{}
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

		// If article has short version
		if !isLong && embed.MakeMeta().HasShort {
			rows = append(rows, discordgo.Button{
				Emoji:    discordgo.ComponentEmoji{Name: "📜"},
				Label:    "Show More",
				Style:    discordgo.PrimaryButton,
				CustomID: "long_" + query,
			})
		}

		// Add buttons if any
		if len(rows) > 0 {
			message.Components = []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: rows,
				},
			}
		}

		// Embed Form
		if isLong {
			message.Embeds = []*discordgo.MessageEmbed{embed.MakeEmbed()}
			message.Flags = discordgo.MessageFlagsEphemeral
		} else {
			message.Embeds = []*discordgo.MessageEmbed{embed.MakeShort()}
		}

		// Send Message
		s.FollowupMessageCreate(i.Interaction, true, message)

	} else {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				FailEmbed,
			},
		})
	}
}

// SearchAutocomplete handles autocomplete interactions for the search command
func SearchAutocomplete(h *DiscordHandler, s *discordgo.Session, i *discordgo.InteractionCreate) {
	var options []*discordgo.ApplicationCommandOptionChoice
	var added = map[string]bool{}
	// Case insensitive
	query := strings.ToLower(i.ApplicationCommandData().Options[0].StringValue())
	// Find up to 25 results with same prefix
	for i := 0; i < len(h.Library.Commands) && len(options) < 25; i++ {
		if strings.HasPrefix(h.Library.Commands[i], query) {
			options = append(options, &discordgo.ApplicationCommandOptionChoice{
				Name:  strings.Title(h.Library.Commands[i]),
				Value: h.Library.Commands[i],
			})
			added[h.Library.Commands[i]] = true
		}
	}
	// Finding additional autocomplete with contains
	for i := 0; i < len(h.Library.Commands) && len(options) < 25; i++ {
		if !added[h.Library.Commands[i]] && strings.Contains(h.Library.Commands[i], query) {
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
