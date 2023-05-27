package main

import (
	"context"
	"github.com/MichaelFraser99/discord-bot/service/internal"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog/log"
	"os"
)

var config internal.BotConfig

func main() {
	publicKey, found := os.LookupEnv("PUBLIC_KEY")
	if !found {
		log.Ctx(context.Background()).Panic().Msg("unable to retrieve public key from environment")
		return
	}

	config = internal.BotConfig{
		PublicKey: publicKey,
		ApplicationCommandHandlers: map[string]func(ctx context.Context, applicationCommand internal.ApplicationCommand) (internal.InteractionResponse, error){
			"poke": func(ctx context.Context, applicationCommand internal.ApplicationCommand) (internal.InteractionResponse, error) {
				return internal.InteractionResponse{
					Type: 4,
					Data: internal.InteractionResponseData{
						Content: "Hello, world!",
						TTS:     false,
					},
				}, nil
			},
		},
	}

	lambda.Start(internal.NewHandler(config))
}
