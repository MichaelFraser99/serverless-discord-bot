package internal

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateRequest(t *testing.T) {
	validBody := "validBody"
	validSignature, validPublicKey, validTimestamp, err := generateSignedArtifacts(validBody)
	require.Nil(t, err)

	tests := []struct {
		name           string
		inputPublicKey string
		inputBody      string
		inputSignature string
		inputTimestamp string
		validate       func(t *testing.T, result bool, err error)
	}{
		{
			name:           "we can validate a request",
			inputPublicKey: string(*validPublicKey),
			inputBody:      validBody,
			inputSignature: string(*validSignature),
			inputTimestamp: *validTimestamp,
			validate: func(t *testing.T, result bool, err error) {
				assert.NoError(t, err)
				assert.True(t, result)
			},
		},
		{
			name:           "we handle invalid signature",
			inputPublicKey: string(*validPublicKey),
			inputBody:      validBody,
			inputSignature: "invalidSignature",
			inputTimestamp: *validTimestamp,
			validate: func(t *testing.T, result bool, err error) {
				assert.Error(t, err)
				assert.False(t, result)
			},
		},
		{
			name:           "we handle invalid public key",
			inputPublicKey: "invalidPublicKey",
			inputBody:      validBody,
			inputSignature: string(*validSignature),
			inputTimestamp: *validTimestamp,
			validate: func(t *testing.T, result bool, err error) {
				assert.Error(t, err)
				assert.False(t, result)
			},
		},
		{
			name:           "we handle invalid timestamp",
			inputPublicKey: string(*validPublicKey),
			inputBody:      validBody,
			inputSignature: string(*validSignature),
			inputTimestamp: "invalidTimestamp",
			validate: func(t *testing.T, result bool, err error) {
				assert.Nil(t, err)
				assert.False(t, result)
			},
		},
		{
			name:           "we handle invalid body",
			inputPublicKey: string(*validPublicKey),
			inputBody:      "invalidBody",
			inputSignature: string(*validSignature),
			inputTimestamp: *validTimestamp,
			validate: func(t *testing.T, result bool, err error) {
				assert.Nil(t, err)
				assert.False(t, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ValidateRequest(context.Background(), tt.inputPublicKey, []byte(tt.inputBody), tt.inputSignature, tt.inputTimestamp)
			tt.validate(t, result, err)
		})
	}
}
