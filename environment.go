package main

import (
	"fmt"
	"os"
)

func ensureValue(key string) string {
	if value, available := os.LookupEnv(key); available {
		return value
	}

	FormatError(fmt.Errorf("failed to obtain environmental value using \"%s\" key", key))
	os.Exit(0)
	return "" // never reaches
}
