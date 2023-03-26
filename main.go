package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	error "discord_voice_assistant/error"
	commands "discord_voice_assistant/commands"
	discordgo "github.com/bwmarrin/discordgo"
)

func main() {
	// Get bot token
	secrets, err := getSecrets()
	if err != nil {
		error.FormatError(fmt.Errorf("failed to obtain secrets: %v", err))
		os.Exit(1)
	}

	botToken, ok := secrets["DISCORD_BOT_TOKEN"]
	if !ok {
		error.FormatError(fmt.Errorf("Discord bot token secret missing"))
		os.Exit(1)
	}

	// Start discord session
	s, err := discordgo.New("Bot " + botToken)
	if err != nil {
		error.FormatError(err)
		os.Exit(1)
	}

	// Set up commands
	cmds := []*discordgo.ApplicationCommand{
		&commands.Add,
	}
	commandHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		commands.Add.Name: commands.AddHandler,
	}


	// Set up intents -- we only want to receive events for Guild and direct interactions
	s.Identify.Intents = discordgo.IntentGuildVoiceStates | discordgo.IntentGuildMessages | discordgo.IntentDirectMessages

	// Add handlers
	s.AddHandlerOnce(func(s *discordgo.Session, r *discordgo.Ready) { log.Println("Bot is up!") })
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	// Open session
	err = s.Open()
	if err != nil {
		error.FormatError(err)
		os.Exit(0)
	}
	defer s.Close()

	// Register commands - override
	// set DISCORD_EXPERIMENTAL_SERVER_ID to cache commands only with a test discord server, otherwise ensure it is empty
	experimentalGuildId, _ := os.LookupEnv("DISCORD_EXPERIMENTAL_SERVER_ID")
	_, err = s.ApplicationCommandBulkOverwrite(s.State.User.ID, experimentalGuildId, cmds)
	if err != nil {
		error.FormatError(fmt.Errorf("Cannot register commands: %v", err))
	}

	// Dummy http server to keep cloudrun happy on deployment
	svr := http.Server{
		Addr: ":8080",
	}
	svr.ListenAndServe()

	// Wait for stop notification and gracefully close
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Gracefully shutting down")
}
