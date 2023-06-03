package handler

import (
	"context"
	"github.com/MichaelFraser99/serverless-discord-bot/internal"
	"github.com/MichaelFraser99/serverless-discord-bot/model"
	"github.com/aws/aws-lambda-go/events"
)

func New(config model.BotConfig) func(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return internal.NewHandler(config)
}
