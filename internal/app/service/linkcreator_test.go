package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateShortString(t *testing.T) {
	t.Run("String length check", func(t *testing.T) {
		result := GenerateShortString("https://ya.ru")
		assert.NotEmpty(t, result)
		assert.Equal(t, 7, len(result))
	})
}
