package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/MichaelFraser99/serverless-discord-bot/model"
	"github.com/MichaelFraser99/serverless-discord-bot/responses"
	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog/log"
)

var config model.BotConfig

func NewHandler(passedConfig model.BotConfig) func(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	config = passedConfig

	return func(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		ctx = log.Logger.WithContext(ctx)
		strEvent, err := json.Marshal(event)
		if err != nil {
			return responses.InternalServerError(ctx, err, "Failed to parse event as json")
		}

		log.Ctx(ctx).Info().Str("event", string(strEvent)).Msg("Processing request")

		bodyBytes := bytes.NewBufferString(event.Body).Bytes()

		interaction := model.Interaction{}
		err = json.Unmarshal(bodyBytes, &interaction)
		if err != nil {
			return responses.InternalServerError(ctx, err, "Failed to parse request body as json")
		}

		validRequest, err := ValidateRequest(ctx, config.PublicKey, bodyBytes, event.Headers["x-signature-ed25519"], event.Headers["x-signature-timestamp"])
		if err != nil {
			return responses.Unauthorized(ctx)
		}

		if !validRequest {
			log.Ctx(ctx).Error().Msg("Signature validation failed")
			return responses.Unauthorized(ctx)
		}
		log.Ctx(ctx).Info().Msg("Signature validation passed")

		var jsonBytes []byte

		switch interaction.Type {
		case 1:
			log.Ctx(ctx).Info().Msg("Ping interaction received")
			jsonBytes, err = json.Marshal(InteractionPing(ctx))
			break
		case 2:
			log.Ctx(ctx).Info().Msg("Application command interaction received")

			applicationCommand := model.ApplicationCommand{}
			if json.Unmarshal([]byte(interaction.Data), &applicationCommand) != nil {
				return responses.InternalServerError(ctx, err, "Failed to parse application command interaction to model")
			}

			interactionResponse, err := interactionApplicationCommand(ctx, applicationCommand)
			if err != nil {
				return responses.InternalServerError(ctx, err, "Failed to process application command interaction")
			}
			jsonBytes, err = json.Marshal(interactionResponse)
			break
		}

		if err != nil {
			return responses.InternalServerError(ctx, err, "Failed to marshal response object as string")
		}

		log.Ctx(ctx).Info().Str("response", string(jsonBytes)).Msg("Sending response")

		return responses.Ok(ctx, string(jsonBytes))
	}
}
