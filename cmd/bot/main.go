package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SHA65536/wildfrostbot"
	"github.com/joho/godotenv"
)

func main() {
	lib, err := wildfrostbot.MakeLibrary()
	godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	discord, err := wildfrostbot.MakeDiscordHandler(os.Getenv("DiscordToken"), lib)
	if err := discord.Start(); err != nil {
		log.Fatal(err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	discord.Stop()
}
