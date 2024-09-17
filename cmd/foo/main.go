package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/replicate/replicate-go"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	// You can also provide a token directly with
	// `replicate.NewClient(replicate.WithToken("r8_..."))`
	r8, err := replicate.NewClient(replicate.WithTokenFromEnv())
	if err != nil {
		log.Fatal(err)
	}

	input := replicate.PredictionInput{
		"prompt": "an astronaut riding a horse on mars, hd, dramatic lighting",
	}

	webhook := replicate.Webhook{
		URL:    "https://webhook.site/214305cd-0416-4f6b-890e-69cb47f43c3a",
		Events: []replicate.WebhookEventType{"start", "completed"},
	}

	// Run a model and wait for its output
	output, err := r8.CreatePrediction(ctx, "d70beb400d223e6432425a5299910329c6050c6abcf97b8c70537d6a1fcb269a", input, &webhook, false)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("output: ", output)

	// // Create a prediction
	// prediction, err := r8.CreatePrediction(ctx, version, input, &webhook, false)
	// if err != nil {
	// 	// handle error
	// }
}
