# Discord Bot
This is a basic template for getting started with a Discord bot. The library is designed to be deployed within AWS lambda behind an API Gateway. Remember these services cost money and scale to a very high degree so take care when deploying.

## Requirements
- go 1.20+
- an AWS account

## Getting Started
1. Clone the repo

## Configuring the bot
This project contains a basic configuration struct which is passed into the handler during the lambda init phase. This struct contains the following fields:
```go
type BotConfig struct {
	PublicKey                  string
	ApplicationCommandHandlers map[string]func(ctx context.Context, applicationCommand ApplicationCommand) (InteractionResponse, error)
}
```

## Adding custom commands
Configuring your discord bot with your own commands is very simple. Inside main.go you will see an example of adding a command to a bot:
```go
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
```
The above code snippet registers a command `poke` and provides a function which will be executed on command run.

Therefor, when a user runs the command `/poke` the bot will respond with "Hello, world!".