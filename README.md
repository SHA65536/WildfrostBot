# Wildfrost Bot
This is a discord bot for the [Wildfrost Wiki](https://wildfrostwiki.com)

## Setup
The bot takes the discord token through environment variables, `DiscordToken="yourtoken"`

Then run `go run ./cmd/bot` to start it!

## Commands
The main command of the bot is /search. It will search the query in the database and return the appropriate embed with a link to the wiki.

## Data
All the data is in the [Data](./Data) directory.

## W.I.P.
[ ] Admin Commands to set config for channel for example ephemeral responses or message commands
[ ] Better error handling
[ ] Fuzzy Search for empty Autocomplete