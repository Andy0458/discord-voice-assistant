package main

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

type Secret struct {
	name  string
	value string
}

func (secret *Secret) Name() string {
	return secret.name
}

func (secret *Secret) Value() string {
	return secret.value
}

func getSecretNames() []string {
	return []string{"DISCORD_BOT_TOKEN"}
}

func getSecrets() (map[string]Secret, error) {
	secrets := make(map[string]Secret)
	ctx := context.Background()
	c, err := secretmanager.NewClient(ctx)
	if err != nil {
		return secrets, err
	}
	defer c.Close()

	for _, secretName := range getSecretNames() {
		req := &secretmanagerpb.AccessSecretVersionRequest{
			Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", ensureValue("GOOGLE_CLOUD_PROJECT"), secretName),
		}
		resp, err := c.AccessSecretVersion(ctx, req)
		if err != nil {
			return secrets, err
		}
		secrets[secretName] = Secret{name: secretName, value: string(resp.GetPayload().GetData())}
	}

	return secrets, nil
}