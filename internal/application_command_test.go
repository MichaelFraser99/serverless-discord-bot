package internal

import (
	"context"
	sdk "github.com/MichaelFraser99/discord-application-sdk/discord/model"
	"github.com/MichaelFraser99/serverless-discord-bot/model"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInteractionApplicationCommand(t *testing.T) {
	tests := []struct {
		name     string
		input    sdk.ApplicationCommand
		validate func(t *testing.T, response *model.InteractionResponse, err error)
	}{
		{
			name: "we can get an invalid command response when no command is registered",
			input: sdk.ApplicationCommand{
				ID:   "1234567890",
				Name: "poke",
				Type: 2,
			},
			validate: func(t *testing.T, response *model.InteractionResponse, err error) {
				require.NoError(t, err)
				assert.Equal(t, 4, response.Type)
				assert.Equal(t, "Unregistered command", response.Data.Content)
				assert.False(t, response.Data.TTS)
			},
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := interactionApplicationCommand(ctx, tt.input)
			tt.validate(t, response, err)
		})
	}
}
