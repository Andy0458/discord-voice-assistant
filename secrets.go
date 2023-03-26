package main

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"cloud.google.com/go/compute/metadata"
)

func getSecretNames() []string {
	return []string{"DISCORD_BOT_TOKEN"}
}

func getSecrets() (map[string]string, error) {
	secrets := make(map[string]string)
	ctx := context.Background()
	c, err := secretmanager.NewClient(ctx)
	if err != nil {
		return secrets, err
	}
	defer c.Close()

	metadataClient := metadata.NewClient(nil)
	projectID, err := metadataClient.ProjectID()
	if err != nil {
		return secrets, err
	}

	for _, secretName := range getSecretNames() {
		req := &secretmanagerpb.AccessSecretVersionRequest{
			Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, secretName),
		}
		resp, err := c.AccessSecretVersion(ctx, req)
		if err != nil {
			return secrets, err
		}
		secrets[secretName] = string(resp.GetPayload().GetData())
	}

	return secrets, nil
}