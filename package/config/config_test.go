package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetup(t *testing.T) {
	opts := *DefaultOptions
	opts.Path = "./fixtures"

	conf, err := New(&opts)

	assert.NoError(t, err)
	assert.Equal(t, 1, conf.GetInt("a"))
	assert.Equal(t, "string", conf.GetString("b"))
	assert.Equal(t, 1, conf.GetInt("c.d"))
	assert.Equal(t, "string", conf.GetString("c.e"))
}
