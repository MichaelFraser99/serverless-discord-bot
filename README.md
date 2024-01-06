# Discord Bot
This module provides the internals to run a discord bot on AWS lambda. The bot is written in go and uses the discord interactions API to receive and respond to commands.

## Requirements
- go 1.20+
- terraform 1.4.x+
- an AWS account

## Interactions Endpoint
This bot leverages the Interaction Endpoint feature of discord bots

https://discord.com/developers/docs/getting-started#adding-interaction-endpoint-url

The endpoint for this bot is your API Gateway endpoint + `/interactions`. For example, if your API Gateway endpoint is `https://example.com` then your interactions endpoint is `https://example.com/interactions`

## Getting Started
1. Run ```go get github.com/MichaelFraser99/serverless-discord-bot```
2. Create a main.go file with the following contents:

```go
package main

import (
	"context"
	"github.com/MichaelFraser99/serverless-discord-bot/handler"
	"github.com/MichaelFraser99/serverless-discord-bot/model"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog/log"
	"os"
)

var config model.BotConfig

func main() {
	publicKey, found := os.LookupEnv("PUBLIC_KEY")
	if !found {
		log.Ctx(context.Background()).Panic().Msg("unable to retrieve public key from environment")
		return
	}

	config = model.BotConfig{
		PublicKey: publicKey,
		ApplicationCommandHandlers: map[string]func(ctx context.Context, applicationCommand model.ApplicationCommand) (*model.InteractionResponse, error){
			"poke": func(ctx context.Context, applicationCommand model.ApplicationCommand) (*model.InteractionResponse, error) {
				return &model.InteractionResponse{
					Type: 4,
					Data: model.InteractionResponseData{
						Content: "Hello, world!",
						TTS:     false,
					},
				}, nil
			},
		},
	}

	lambda.Start(handler.New(config))
}
```
The above code snippet registers a command `poke` and provides a function which will be executed on command run.

Therefor, when a user runs the command `/poke` the bot will respond with "Hello, world!".

## Configuring the bot
This module exposes a basic configuration struct which is passed into the handler during the lambda init phase. This struct contains the following fields:
```go
type BotConfig struct {
	PublicKey                  string
	ApplicationCommandHandlers map[string]func(ctx context.Context, applicationCommand ApplicationCommand) (InteractionResponse, error)
}
```

## Deploying the bot
I'd recommend leveraging terraform to deploy the bot
A supporting terraform provider for configuring discord application commands can be found here: https://registry.terraform.io/providers/MichaelFraser99/discord-application/latest

This project packages a `./getting-started` directory which contains all you need to get started with deploying your bot

This directory contains an example swagger file to populate an API Gateway endpoint with the required routes to support the discord interactions API. This swagger file can be found at `./getting-started/api/swagger.yaml`

Finally, the `./getting-started/infrastructure` directory contains a terraform module which can be used to deploy the bot
