package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/MichaelFraser99/serverless-discord-bot/handler"
	"github.com/MichaelFraser99/serverless-discord-bot/model"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog/log"
)

var config model.BotConfig

func main() {
	publicKey, found := os.LookupEnv("PUBLIC_KEY")
	if !found {
		log.Ctx(context.Background()).Panic().Msg("unable to retrieve public key from environment")
		return
	}

	config = model.BotConfig{
		PublicKey: publicKey,
		ApplicationCommandHandlers: map[string]func(ctx context.Context, applicationCommand model.ApplicationCommand) (*model.InteractionResponse, error){
			"create": func(ctx context.Context, applicationCommand model.ApplicationCommand) (*model.InteractionResponse, error) {
				// Extract values from applicationCommand options
				modsEnabled := "false"
				mapValue := "gridmap_v2"
				maxPlayers := "5"
				maxCars := "1"
				private := "true"
				if applicationCommand.Options != nil {
					appCommands := applicationCommand.Options.([]interface{})
					for _, cmdInterface := range appCommands {
						cmd, ok := cmdInterface.(map[string]interface{})
						if !ok {
							log.Ctx(ctx).Error().Msgf("unable to cast application command option to model.ApplicationCommandOption %v", cmdInterface)
							continue
						}
						if cmd["value"] != nil {
							switch cmd["name"].(string) {
							case "modded":
								modsEnabled = cmd["value"].(string)
							case "map":
								mapValue = cmd["value"].(string)
							case "max_players":
								maxPlayers = cmd["value"].(string)
							case "max_cars":
								maxCars = cmd["value"].(string)
							case "private":
								private = cmd["value"].(string)
							default:
								log.Ctx(ctx).Warn().Msgf("unrecognised command option recieved %s", cmd["name"].(string))
							}
						} else {
							log.Ctx(ctx).Error().Msgf("command option recieved without value")
						}
					}
				} else {
					log.Ctx(ctx).Info().Msgf("application command recieved without options")
				}
				// Create the payload for the POST request
				payload := map[string]interface{}{
					"ref": "main",
					"inputs": map[string]interface{}{
						"action":      "apply",
						"map":         mapValue,
						"modded":      modsEnabled,
						"max_players": maxPlayers,
						"max_cars":    maxCars,
						"private":     private,
					},
				}

				// Convert the payload to JSON
				payloadBytes, err := json.Marshal(payload)
				if err != nil {
					return nil, err
				}

				// Make the POST request to GitHub
				err = ghPostRequest(payloadBytes)
				if err != nil {
					log.Ctx(ctx).Error().Msg("unable to make POST request to GitHub")
					return nil, err
				}

				// Return the interaction response
				return &model.InteractionResponse{
					Type: 4,
					Data: model.InteractionResponseData{
						Content: "Server Creation Initiated!",
						TTS:     false,
					},
				}, nil
			},
			"destroy": func(ctx context.Context, applicationCommand model.ApplicationCommand) (*model.InteractionResponse, error) {
				// Create the payload for the POST request
				payload := map[string]interface{}{
					"ref": "main",
					"inputs": map[string]interface{}{
						"action": "destroy",
					},
				}

				// Convert the payload to JSON
				payloadBytes, err := json.Marshal(payload)
				if err != nil {
					return nil, err
				}

				// Make the POST request to GitHub
				err = ghPostRequest(payloadBytes)
				if err != nil {
					log.Ctx(ctx).Error().Msg("unable to make POST request to GitHub")
					return nil, err
				}
				return &model.InteractionResponse{
					Type: 4,
					Data: model.InteractionResponseData{
						Content: "Server Destruction Initiated!",
						TTS:     false,
					},
				}, nil
			},
		},
	}
	lambda.Start(handler.New(config))
}

func ghPostRequest(payloadBytes []byte) error {
	req, err := http.NewRequest("POST", "https://api.github.com/repos/Harry-Moore-dev/tf-beamMP-server-deployment/actions/workflows/79678424/dispatches", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}
	token, found := os.LookupEnv("GITHUB_TOKEN")
	if !found {
		return fmt.Errorf("unable to retrieve github token from environment")
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected response status code: %d", resp.StatusCode)
	}
	return nil
}
