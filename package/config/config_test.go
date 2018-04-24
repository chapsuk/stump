package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetup(t *testing.T) {
	t.Run("first", func(t *testing.T) {
		opts := *DefaultOptions
		opts.Path = "./fixtures/first"

		conf, err := New(&opts)

		assert.NoError(t, err)
		assert.Equal(t, 1, conf.GetInt("a"))
		assert.Equal(t, "string", conf.GetString("b"))
		assert.Equal(t, 1, conf.GetInt("c.d"))
		assert.Equal(t, "string", conf.GetString("c.e"))
	})

	t.Run("second", func(t *testing.T) {
		opts := *DefaultOptions
		opts.Path = "./fixtures/second"

		conf, err := New(&opts)

		assert.NoError(t, err)
		assert.Equal(t, 1, conf.GetInt("a"))
		assert.Equal(t, "string", conf.GetString("b"))
		assert.Equal(t, 1, conf.GetInt("c.d"))
		assert.Equal(t, "string", conf.GetString("c.e"))
	})
}
