package internal

import (
	"context"
	"github.com/rs/zerolog/log"
)

func InteractionPing(ctx context.Context) InteractionResponse {
	log.Ctx(ctx).Info().Msg("Responding to ping")
	return InteractionResponse{
		Type: 1,
	}
}
