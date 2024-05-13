package controllers

import (
	"auth-api/internal/pkg/tools"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOAuthController_StateValidation(t *testing.T) {
	testCases := []struct {
		name  string
		valid bool
		state string
	}{
		{
			name:  "invalid state",
			valid: false,
			state: "invalid_state",
		},
		{
			name:  "nil state",
			valid: false,
			state: "",
		},
		{
			name:  "invalid state",
			valid: false,
			state: "invalid_state",
		},
		{
			name:  "valid state: new uuid",
			valid: true,
			state: uuid.New().String(),
		},
		{
			name:  "valid state: copied uuid",
			valid: true,
			state: "ac575f74-bed1-4181-9ed2-1df726774044",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.valid, tools.IsUUID(testCase.state), testCase.name)
		})
	}
}
