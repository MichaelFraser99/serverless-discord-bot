package internal

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"github.com/rs/zerolog/log"
)

func ValidateRequest(ctx context.Context, publicKey string, body []byte, signature string, timestamp string) bool {
	log.Ctx(ctx).Info().Msg("validating message signature")
	hexDecodedPublicKey, err := hex.DecodeString(publicKey)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to hex decode public key")
		return false
	}

	hexDecodedSignature, err := hex.DecodeString(signature)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to hex decode signature")
		return false
	}

	message := append([]byte(timestamp), body...)
	//validate Ed25519 signature
	return ed25519.Verify(hexDecodedPublicKey, message, hexDecodedSignature)
}
