package responses

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog/log"
)

func InternalServerError(ctx context.Context, err error, cause string) (events.APIGatewayProxyResponse, error) {
	log.Ctx(ctx).Error().Err(err).Msg(cause)
	return events.APIGatewayProxyResponse{
		StatusCode: 500,
	}, err
}

func Ok(ctx context.Context, response string) (events.APIGatewayProxyResponse, error) {
	log.Ctx(ctx).Info().Msg("Request successful")
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       response,
	}, nil
}

func Unauthorized(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	log.Ctx(ctx).Info().Msg("Request unauthorized")
	return events.APIGatewayProxyResponse{
		StatusCode: 401,
	}, nil
}
