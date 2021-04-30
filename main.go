package main

import (
	"context"
	"fmt"
	harness "github.com/drone/ff-golang-server-sdk/client"
	"github.com/drone/ff-golang-server-sdk/dto"
	"log"
	"time"
)

const sdkKey = "78e171de-fd0f-433b-862a-2fe5db69318e"

const featureFlagKey = "Dark_Mode"

func main() {
	target := dto.NewTargetBuilder("john").
		Firstname("John").
		Lastname("Doe").
		Email("john@doe.com").
		Country("USA").
		Custom("height", 186).
		Name("john").
		Build()

	client, err := harness.NewCfClient(sdkKey,
		harness.WithURL("http://34.82.119.242/api/1.0/"),
		harness.WithStreamEnabled(true),
		harness.WithTarget(target),
	)
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf("error while closing client err: %v", err)
		}
	}()

	if err != nil {
		log.Printf("could not connect to CF servers %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				showFeature, err := client.BoolVariation(featureFlagKey, &target, false)

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
