package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name   string
		value  string
		result bool
	}{
		{
			name:   "Valid URL",
			value:  "https://yandex.com",
			result: true,
		},
		{
			name:   "Invalid URL",
			value:  "http://",
			result: false,
		},
		{
			name:   "Empty string",
			value:  "",
			result: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualResult := ValidateURL(tt.value)
			assert.Equal(t, tt.result, actualResult)
		})
	}
}
