package main

import (
	"fmt"
	"os"

	error "discord_voice_assistant/error"
)

func ensureValue(key string) string {
	if value, available := os.LookupEnv(key); available {
		return value
	}

	error.FormatError(fmt.Errorf("failed to obtain environmental value using \"%s\" key", key))
	os.Exit(0)
	return "" // never reaches
}
