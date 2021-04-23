package main

import (
	"context"
	"fmt"
	"log"
	"time"


        harness "github.com/drone/ff-golang-server-sdk/client"
        "github.com/drone/ff-golang-server-sdk/dto"
)

const sdkKey = "7bc7df92-16bc-4c62-9ad4-b78a497690e7"

const featureFlagKey = "toggle"

func main() {

	client, err := harness.NewCfClient(sdkKey,
		harness.WithUrl("http://34.82.119.242/api/1.0/"),
		harness.WithEventSourceUrl("http://34.82.119.242/api/1.0/stream/environments/%s"),
		harness.WithStreamEnabled(true),
	)
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf("error while closing client err: %v", err)
		}
	}()

	if err != nil {
		log.Printf("could not connect to CF servers %v", err)
	}

	target := dto.NewTargetBuilder("john").
		Firstname("John").
		Lastname("Doe").
		Email("john@doe.com").
		Country("USA").
		Custom("height", 186).
		Build()

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				showFeature, err := client.BoolVariation(featureFlagKey, target, false)

				if err != nil {
					fmt.Printf("Error getting value: %v", err)
				}

				fmt.Printf("KeyFeature flag '%s' is %t for this user\n", featureFlagKey, showFeature)
				time.Sleep(10 * time.Second)
			}
		}
	}()

	time.Sleep(5 * time.Minute)
	cancel()
}
