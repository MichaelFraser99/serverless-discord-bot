package internal

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	HEADER_TIMESTAMP = "x-signature-timestamp"
	HEADER_SIGNATURE = "x-signature-ed25519"
)

func generateSignedArtifacts(message string) (hexSignature, hexPublicKey *[]byte, timestamp *string, err error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, nil, err
	}

	hpk := make([]byte, hex.EncodedLen(len(publicKey)))
	hex.Encode(hpk, publicKey)

	ts := fmt.Sprint(time.Now().Unix())
	rawSignature := string(ed25519.Sign(privateKey, append([]byte(ts), []byte(message)...)))

	hs := make([]byte, hex.EncodedLen(len(rawSignature)))
	hex.Encode(hs, []byte(rawSignature))

	return &hs, &hpk, &ts, nil
}

func TestHandler(t *testing.T) {

	interactionTests := []struct {
		name                      string
		message                   string
		validateHappyPathResponse func(t *testing.T, response InteractionResponse)
	}{
		{
			name:    "ping",
			message: `{"type":1}`,
			validateHappyPathResponse: func(t *testing.T, response InteractionResponse) {
				assert.Equal(t, 1, response.Type)
			},
		},
		{
			name:    "application command:poke",
			message: `{"type":2,"data":{"id":"123456789","name":"poke","options":[{"name":"poke","value":"test"}]}}`,
			validateHappyPathResponse: func(t *testing.T, response InteractionResponse) {
				assert.Equal(t, 4, response.Type)
				assert.Equal(t, "Hello, world!", response.Data.Content)
				assert.False(t, response.Data.TTS)
			},
		},
	}

	for _, interaction := range interactionTests {

		hexSignature, hexPublicKey, timestamp, err := generateSignedArtifacts(interaction.message)
		require.Nil(t, err, "failed to generate signed artifacts")

		testConfig := BotConfig{
			PublicKey: string(*hexPublicKey),
			ApplicationCommandHandlers: map[string]func(ctx context.Context, applicationCommand ApplicationCommand) (InteractionResponse, error){
				"poke": func(ctx context.Context, applicationCommand ApplicationCommand) (InteractionResponse, error) {
					return InteractionResponse{
						Type: 4,
						Data: InteractionResponseData{
							Content: "Hello, world!",
							TTS:     false,
						},
					}, nil
				},
			},
		}

		tests := []struct {
			name     string
			input    events.APIGatewayProxyRequest
			validate func(t *testing.T, response events.APIGatewayProxyResponse, bodyValidation func(t *testing.T, response InteractionResponse), err error)
		}{
			{
				name: "valid request",
				input: events.APIGatewayProxyRequest{
					Body: interaction.message,
					Headers: map[string]string{
						HEADER_SIGNATURE: string(*hexSignature),
						HEADER_TIMESTAMP: *timestamp,
					},
				},
				validate: func(t *testing.T, response events.APIGatewayProxyResponse, bodyValidation func(t *testing.T, response InteractionResponse), err error) {
					assert.Nil(t, err)
					assert.NotNil(t, response.Body)
					assert.Equal(t, 200, response.StatusCode)

					interactionResponse := InteractionResponse{}
					assert.Nil(t, json.Unmarshal([]byte(response.Body), &interactionResponse))

					bodyValidation(t, interactionResponse)
				},
			},
			{
				name: "invalid signature",
				input: events.APIGatewayProxyRequest{
					Body: interaction.message,
					Headers: map[string]string{
						HEADER_SIGNATURE: "oogabooga",
						HEADER_TIMESTAMP: *timestamp,
					},
				},
				validate: func(t *testing.T, response events.APIGatewayProxyResponse, bodyValidation func(t *testing.T, response InteractionResponse), err error) {
					assert.Nil(t, err)
					assert.Empty(t, response.Body)
					assert.Equal(t, 401, response.StatusCode)
				},
			},
			{
				name: "invalid timestamp",
				input: events.APIGatewayProxyRequest{
					Body: interaction.message,
					Headers: map[string]string{
						HEADER_SIGNATURE: string(*hexSignature),
						HEADER_TIMESTAMP: "12345",
					},
				},
				validate: func(t *testing.T, response events.APIGatewayProxyResponse, bodyValidation func(t *testing.T, response InteractionResponse), err error) {
					assert.Nil(t, err)
					assert.Empty(t, response.Body)
					assert.Equal(t, 401, response.StatusCode)
				},
			},
			{
				name: "invalid body",
				input: events.APIGatewayProxyRequest{
					Body: "someutternonsense",
					Headers: map[string]string{
						HEADER_SIGNATURE: string(*hexSignature),
						HEADER_TIMESTAMP: "12345",
					},
				},
				validate: func(t *testing.T, response events.APIGatewayProxyResponse, bodyValidation func(t *testing.T, response InteractionResponse), err error) {
					assert.Error(t, err)
					assert.Empty(t, response.Body)
					assert.Equal(t, 500, response.StatusCode)
				},
			},
		}

		for _, tt := range tests {
			t.Run(fmt.Sprintf("%s-%s", interaction.name, tt.name), func(t2 *testing.T) {
				resp, err := NewHandler(testConfig)(context.Background(), tt.input)
				tt.validate(t, resp, interaction.validateHappyPathResponse, err)
			})
		}
	}
}
