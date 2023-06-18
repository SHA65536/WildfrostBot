package wildfrostbot

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"go.etcd.io/bbolt"
)

type DiscordHandler struct {
	Session     *discordgo.Session
	Library     *Library
	Cmds        map[string]*Cmd
	CreatedCmds []*discordgo.ApplicationCommand
	ChannelDB   *bbolt.DB
}

// Cmd is a helper struct for registering commands
type Cmd struct {
	Cmd     *discordgo.ApplicationCommand
	Handler func(h *DiscordHandler, s *discordgo.Session, i *discordgo.InteractionCreate)
}

func MakeDiscordHandler(token string, lib *Library) (*DiscordHandler, error) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	var h = &DiscordHandler{
		Session: dg,
		Library: lib,
	}
	h.CreateDB("channels.db")
	dg.Identify.Intents = discordgo.IntentsGuildMessages
	dg.AddHandler(h.InteractionHandler)
	return h, nil
}

// InteractionHandler routes the interactions into the right handler
func (h *DiscordHandler) InteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var name string
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		name = i.ApplicationCommandData().Name
	case discordgo.InteractionMessageComponent:
		name = strings.Split(i.MessageComponentData().CustomID, "_")[0]
	case discordgo.InteractionApplicationCommandAutocomplete:
		SearchAutocomplete(h, s, i)
		return
	}

	if cmd, ok := h.Cmds[name]; ok {
		cmd.Handler(h, s, i)
	}
}

// RegisterCommands register the slash commands
func (h *DiscordHandler) RegisterCommands() error {
	var err error
	h.Cmds = map[string]*Cmd{
		"search": &SearchHandler,
	}
	commands := make([]*discordgo.ApplicationCommand, 0)
	for _, cmd := range h.Cmds {
		commands = append(commands, cmd.Cmd)
	}
	h.Cmds["long"] = &SearchHandler
	h.CreatedCmds, err = h.Session.ApplicationCommandBulkOverwrite(h.Session.State.User.ID, "", commands)
	return err
}

// Start starts the bot
func (h *DiscordHandler) Start() error {
	err := h.Session.Open()
	if err != nil {
		return err
	}
	log.Printf("Bot started!")
	return h.RegisterCommands()
}

// Stop stops the bot
func (h *DiscordHandler) Stop() error {
	log.Printf("Bot stopped!")
	return h.Session.Close()
}

// SendEmbed is a helper method to send an embed
func SendEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, m *discordgo.MessageEmbed) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				m,
			},
		},
	})
}

var FailEmbed = &discordgo.MessageEmbed{
	Title:       "Command failed! D:",
	Description: "Please try again! If this persists, contact someone important!",
	Color:       0xFF0000,
}
