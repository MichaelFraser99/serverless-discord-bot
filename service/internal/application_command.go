package internal

import (
	"context"
	"github.com/rs/zerolog/log"
)

func interactionApplicationCommand(ctx context.Context, data ApplicationCommand) (InteractionResponse, error) {
	log.Ctx(ctx).Info().Interface("message", data).Msg("Processing application command")

	for command, commandHandler := range config.ApplicationCommandHandlers {
		if command == data.Name {
			return commandHandler(ctx, data)
		}
	}

	return InteractionResponse{
		Type: 4,
		Data: InteractionResponseData{
			Content: "Unregistered command",
			TTS:     false,
		},
	}, nil
}
