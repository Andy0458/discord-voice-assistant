package main

import (
	"discord_voice_assistant/commands"
	"fmt"
	"log"
	"os"
	"time"

	tempest "github.com/Amatsagu/Tempest"
)

func main() {
	secrets, err := getSecrets()
	if err != nil {
		FormatError(fmt.Errorf("failed to obtain secrets: %v", err))
		os.Exit(0)
	}

	botTokenSecret, ok := secrets["DISCORD_BOT_TOKEN"]
	if !ok {
		FormatError(fmt.Errorf("Discord bot token secret missing"))
		os.Exit(0)
	}
	botToken := botTokenSecret.Value()

	client := tempest.CreateClient(tempest.ClientOptions{
		ApplicationID: tempest.StringToSnowflake(ensureValue("DISCORD_APP_ID")),
		PublicKey:     ensureValue("DISCORD_PUBLIC_KEY"),
		Token:         "Bot " + botToken,
		PreCommandExecutionHandler: func(itx tempest.CommandInteraction) *tempest.ResponseData {
			log.Printf("%v", itx)
			return nil
		},
		Cooldowns: &tempest.ClientCooldownOptions{
			Duration:  time.Second * 3,
			Ephemeral: true,
			CooldownResponse: func(user tempest.User, timeLeft time.Duration) tempest.ResponseData {
				return tempest.ResponseData{
					Content: fmt.Sprintf("You're still on cooldown! Try again in **%.2fs**.", timeLeft.Seconds()),
				}
			},
		},
	})

	addr := fmt.Sprintf("0.0.0.0:%s", ensureValue("PORT"))

	client.RegisterCommand(commands.Add)
	client.SyncCommands([]tempest.Snowflake{}, nil, false)

	log.Printf("Starting application at %s", addr)
	log.Printf("Latency: %dms", client.Ping().Milliseconds())

	if err := client.ListenAndServe("/", addr); err != nil {
		// Will happen in situation where normal std/http would panic so most likely never.
		FormatError(err)
		os.Exit(1)
	}
}