package wildfrostbot

import "github.com/bwmarrin/discordgo"

type Embeder interface {
	MakeEmbed() *discordgo.MessageEmbed
	MakeMeta() Metadata
}

type Metadata struct {
	Title   string            `json:"title"`
	Aliases []string          `json:"aliases"`
	Tags    []string          `json:"tags"`
	Keys    map[string]string `json:"keys"`
}

type Article struct {
	Metadata Metadata                `json:"metadata"`
	Embed    *discordgo.MessageEmbed `json:"embed"`
}

func (a *Article) MakeEmbed() *discordgo.MessageEmbed {
	return a.Embed
}

func (a *Article) MakeMeta() Metadata {
	return a.Metadata
}
