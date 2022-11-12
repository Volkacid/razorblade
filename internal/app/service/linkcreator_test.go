package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGenerateShortString(t *testing.T) {
	SetCreatorSeed(time.Now().Unix())
	t.Run("String length check", func(t *testing.T) {
		result := GenerateShortString()
		assert.NotEmpty(t, result)
		assert.Equal(t, 6, len(result))
	})
}
