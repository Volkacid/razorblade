package service

import (
	"github.com/Volkacid/razorblade/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGenerateShortString(t *testing.T) {
	SetCreatorSeed(time.Now().Unix())
	db := storage.CreateTestStorage()
	t.Run("String length check", func(t *testing.T) {
		result := GenerateShortString(db)
		assert.NotEmpty(t, result)
		assert.Equal(t, 6, len(result))
	})
}
