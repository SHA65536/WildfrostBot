package wildfrostbot

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const descLimit int = 4096
const fieldLimit int = 1024
const maxWords int = 100

var contentPath string = filepath.Join("Data", "Content")
var categoryPath string = filepath.Join("Data", "Categories")
var dictionaryPath string = filepath.Join("Data", "dictionary.json")

var replaceRegex = regexp.MustCompile(`<&([^<>&]*)&>`)

var replacementDict map[string]string

// MakeArticles loads the articles from disk
func (l *Library) MakeArticles() error {
	return filepath.WalkDir(contentPath, func(path string, d fs.DirEntry, ferr error) error {
		var art Article

		// Ignored files
		if d.IsDir() || !strings.HasSuffix(path, ".json") {
			return nil
		}

		// Open file
		jsonFile, err := os.Open(path)
		if err != nil {
			return err
		}

		// Unmarshal into article
		if err = json.NewDecoder(jsonFile).Decode(&art); err != nil {
			return err
		}

		// Replace the dictionary
		replaceEmbed(art.Embed)

		// Make Short Version
		art.Short, art.Metadata.HasShort = MakeShort(art.Embed)

		// Censor spoilers
		if art.Metadata.Keys["spoiler"] == "true" {
			censorEmbed(art.Embed)
			censorEmbed(art.Short)
		}

		// Check length limits
		if checkToolong(art.Embed) {
			return fmt.Errorf("embed too long: %s", art.Metadata.Title)
		}

		l.Articles[art.Metadata.Title] = &art
		return nil
	})
}

// MakeCategories loads categories from disk and generates
// embeds for them
func (l *Library) MakeCategories() error {
	return filepath.WalkDir(categoryPath, func(path string, d fs.DirEntry, ferr error) error {
		var art Article
		// Ignored files
		if d.IsDir() || !strings.HasSuffix(path, ".json") {
			return nil
		}

		// Open file
		jsonFile, err := os.Open(path)
		if err != nil {
			return err
		}

		// Unmarshal into category
		if err = json.NewDecoder(jsonFile).Decode(&art); err != nil {
			return err
		}

		// Replace the dictionary
		replaceEmbed(art.Embed)

		// Add articles to category
		var vals []string
		for _, other := range l.Articles {
			for _, t := range other.MakeMeta().Tags {
				if t == art.Metadata.Title {
					vals = append(vals, other.MakeMeta().Title)
					break
				}
			}
		}
		addListToCategory(art.Embed, vals)

		// Make Short Version
		art.Short, art.Metadata.HasShort = MakeShort(art.Embed)

		// Censor spoilers
		if art.Metadata.Keys["spoiler"] == "true" {
			censorEmbed(art.Embed)
			censorEmbed(art.Short)
		}

		// Check length limits
		if checkToolong(art.Embed) {
			return fmt.Errorf("embed too long: %s", art.Metadata.Title)
		}

		l.Articles[art.Metadata.Title] = &art
		return nil
	})
}

// MakeShort returns a short version of an embed, and a boolean
// representing whether the short version is different from the long one
func MakeShort(long *discordgo.MessageEmbed) (*discordgo.MessageEmbed, bool) {
	// Copying common fields
	var res = &discordgo.MessageEmbed{
		URL: long.URL, Title: long.Title, Color: long.Color,
		Footer: long.Footer, Image: long.Image, Thumbnail: long.Thumbnail,
		Author: long.Author,
	}

	var words int
	// Checking description
	longWords := strings.Split(long.Description, " ")
	words = len(longWords)
	if words >= maxWords {
		res.Description = strings.Join(longWords[:maxWords], " ") + "..."
		return res, true
	} else {
		res.Description = long.Description
	}

	// Checking embed fields
	for _, field := range long.Fields {
		shortField := &discordgo.MessageEmbedField{Name: field.Name, Inline: field.Inline}
		fieldWords := strings.Split(field.Value, " ")

		if words+len(fieldWords) >= maxWords {
			shortField.Value = strings.Join(fieldWords[:maxWords-words], " ") + "..."
			res.Fields = append(res.Fields, shortField)
			return res, true
		}
		shortField.Value = field.Value
		res.Fields = append(res.Fields, shortField)
		words += len(fieldWords)
	}

	return res, false
}

// CreateAliases generates aliases into the article list
func (l *Library) CreateAliases() {
	var newAliases = map[string]Embeder{}
	// Collect all aliases
	for _, v := range l.Articles {
		for _, alias := range v.MakeMeta().Aliases {
			newAliases[alias] = v
		}
	}
	// Add them to article list
	for k, v := range newAliases {
		l.Articles[k] = v
	}
}

// replaceEmbed replaces dictionary definitions in embed
func replaceEmbed(embed *discordgo.MessageEmbed) {
	embed.Description = replaceRegex.ReplaceAllStringFunc(embed.Description, replaceFunc)
	if embed.Footer != nil {
		embed.Footer.Text = replaceRegex.ReplaceAllStringFunc(embed.Footer.Text, replaceFunc)
	}
	for i := range embed.Fields {
		embed.Fields[i].Value = replaceRegex.ReplaceAllStringFunc(embed.Fields[i].Value, replaceFunc)
	}
}

// replaceFunc is a regex replace function for the dictionary
func replaceFunc(input string) string {
	stripped := input[2 : len(input)-2]
	if res, ok := replacementDict[stripped]; ok {
		return res
	}
	return input
}

// censorEmbed adds spoilers to an embed
func censorEmbed(embed *discordgo.MessageEmbed) {
	embed.Description = fmt.Sprintf("||%s||", embed.Description)
	for i := range embed.Fields {
		embed.Fields[i].Name = fmt.Sprintf("||%s||", embed.Fields[i].Name)
		embed.Fields[i].Value = fmt.Sprintf("||%s||", embed.Fields[i].Value)
	}
}

// addListToCategory adds the list of articles into a category embed
func addListToCategory(embed *discordgo.MessageEmbed, vals []string) {
	var fieldIdx int
	sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })

	for _, val := range vals {
		if len(embed.Fields[fieldIdx].Value)+len(val)+2 > fieldLimit {
			embed.Fields[fieldIdx].Value = embed.Fields[fieldIdx].Value[:len(embed.Fields[fieldIdx].Value)-2]
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   embed.Fields[0].Name,
				Value:  strings.Title(val) + ", ",
				Inline: embed.Fields[0].Inline,
			})
			fieldIdx++
		} else {
			embed.Fields[fieldIdx].Value += strings.Title(val) + ", "
		}
	}
	embed.Fields[fieldIdx].Value = embed.Fields[fieldIdx].Value[:len(embed.Fields[fieldIdx].Value)-2]
}

// checkTooLong returns true if the given embed is too long
func checkToolong(embed *discordgo.MessageEmbed) bool {
	if len(embed.Description) > descLimit {
		return true
	}
	for i := range embed.Fields {
		if len(embed.Fields[i].Value) > fieldLimit {
			return true
		}
	}
	return false
}
