package internal

import (
	"context"
	"github.com/MichaelFraser99/serverless-discord-bot/model"
	"github.com/rs/zerolog/log"
)

func InteractionPing(ctx context.Context) model.InteractionResponse {
	log.Ctx(ctx).Info().Msg("Responding to ping")
	return model.InteractionResponse{
		Type: 1,
	}
}
