package internal

import (
	"context"
	sdk "github.com/MichaelFraser99/discord-application-sdk/discord/model"
	"github.com/MichaelFraser99/serverless-discord-bot/model"
	"github.com/rs/zerolog/log"
)

func interactionApplicationCommand(ctx context.Context, data sdk.ApplicationCommand) (*model.InteractionResponse, error) {
	log.Ctx(ctx).Info().Interface("message", data).Msg("Processing application command")

	for command, commandHandler := range config.ApplicationCommandHandlers {
		if command == data.Name {
			return commandHandler(ctx, data)
		}
	}

	return &model.InteractionResponse{
		Type: 4,
		Data: model.InteractionResponseData{
			Content: "Unregistered command",
			TTS:     false,
		},
	}, nil
}
