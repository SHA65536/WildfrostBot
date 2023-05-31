package wildfrostbot

import "github.com/bwmarrin/discordgo"

type Embeder interface {
	MakeEmbed() *discordgo.MessageEmbed
	MakeShort() *discordgo.MessageEmbed
	MakeMeta() Metadata
}

type Metadata struct {
	Title    string            `json:"title"`
	Aliases  []string          `json:"aliases"`
	Tags     []string          `json:"tags"`
	Keys     map[string]string `json:"keys"`
	HasShort bool
}

type Article struct {
	Metadata Metadata                `json:"metadata"`
	Embed    *discordgo.MessageEmbed `json:"embed"`
	Short    *discordgo.MessageEmbed
}

func (a *Article) MakeEmbed() *discordgo.MessageEmbed {
	return a.Embed
}

func (a *Article) MakeShort() *discordgo.MessageEmbed {
	return a.Short
}

func (a *Article) MakeMeta() Metadata {
	return a.Metadata
}
