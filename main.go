package main

import (
	"context"
	"fmt"
	harness "github.com/harness/ff-golang-server-sdk/client"
	"github.com/harness/ff-golang-server-sdk/dto"
	"log"
	"time"
)

const sdkKey = "6f09aba4-d4eb-4d15-bb9d-843a8af4f220"

const featureFlagKey = "SimpleBoolean"

func main() {
	target := dto.NewTargetBuilder("john").
		Name("John").
		Custom("customAttribute", 123).
		Build()

	client, err := harness.NewCfClient(sdkKey,
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
