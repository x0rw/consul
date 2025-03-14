package main

import (
	"fmt"
	"time"

	"github.com/hashicorp/consul/api"
)

func main() {
	// Create a Consul client
	config := api.DefaultConfig()
	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}

	catalog := client.Catalog()

	_, err = catalog.Register(&api.CatalogRegistration{
		ID:      "158910b8-1946-45d1-b7c1-3511d419b9ac",
		Node:    "fake-esm",
		Address: "127.0.0.1",
	}, nil)

	if err != nil {
		panic(err)
	}

	fmt.Println("Registered Node")

	sessionClient := client.Session()
	sessionID, _, err := sessionClient.CreateNoChecks(&api.SessionEntry{
		Name: "esm-session",
		Node: "fake-esm",
		TTL:  "10s",
	}, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("Created Session")

	doneCh := make(chan struct{})
	go sessionClient.RenewPeriodic("1s", sessionID, nil, doneCh)

	_, err = catalog.Register(&api.CatalogRegistration{
		ID:             "158910b8-1946-45d1-b7c1-3511d419b9ac",
		Node:           "fake-esm",
		Address:        "127.0.0.1",
		SkipNodeUpdate: true,
		Checks: api.HealthChecks{
			{
				Node:    "fake-esm",
				CheckID: "esm-session",
				Name:    "esm-session",
				Status:  "passing",
				Type:    "session",
				Notes:   "check will go critical when the specified session is invalidated",
				Definition: api.HealthCheckDefinition{
					SessionID: sessionID,
				},
			},
		},
	}, nil)

	if err != nil {
		panic(err)
	}

	fmt.Println("Registered Check")

	time.Sleep(60 * time.Second)
	close(doneCh)
}
