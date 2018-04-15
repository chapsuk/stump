package logger

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("Basic logger (Development)", func(t *testing.T) {
		logger, err := New(nil)
		assert.NoError(t, err)
		assert.NotNil(t, logger)
	})

	t.Run("Basic logger (Production)", func(t *testing.T) {
		logger, err := New(&Options{
			Level: LoggerLevelProduction,
		})
		assert.NoError(t, err)
		assert.NotNil(t, logger)
	})

	t.Run("Unknown level", func(t *testing.T) {
		logger, err := New(&Options{
			Level: 3,
		})

		assert.Equal(t, ErrUnknownLevel, err)
		assert.Nil(t, logger)
	})

	t.Run("It write down logs in proper way", func(t *testing.T) {
		logger, err := New(nil)
		assert.NoError(t, err)
		assert.NotNil(t, logger)
	})
}
