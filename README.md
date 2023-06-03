# Discord Bot
This module provides the internal to run a discord bot on AWS lambda. The bot is written in go and uses the discord interactions API to receive and respond to commands.

## Requirements
- go 1.20+
- terraform 1.4.x+
- an AWS account

## Getting Started
todo

## Interactions Endpoint
This bot leverages the Interaction Endpoint feature of discord bots

https://discord.com/developers/docs/getting-started#adding-interaction-endpoint-url

The endpoint for this bot is your API Gateway endpoint + `/interactions`. For example, if your API Gateway endpoint is `https://example.com` then your interactions endpoint is `https://example.com/interactions`

## Configuring the bot
This project contains a basic configuration struct which is passed into the handler during the lambda init phase. This struct contains the following fields:
```go
type BotConfig struct {
	PublicKey                  string
	ApplicationCommandHandlers map[string]func(ctx context.Context, applicationCommand ApplicationCommand) (InteractionResponse, error)
}
```

## Adding custom commands
Configuring your discord bot with your own commands is very simple. Below is an example 'main.go' file which registers a single command called 'poke'
```go
package main

import (
	"context"
	"github.com/MichaelFraser99/discord-bot/internal"
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

```
The above code snippet registers a command `poke` and provides a function which will be executed on command run.

Therefor, when a user runs the command `/poke` the bot will respond with "Hello, world!".
