package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetup(t *testing.T) {
	conf, err := New(&Options{
		Path:    "./fixtures",
		AutoEnv: true,
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, conf.GetInt("a"))
	assert.Equal(t, "string", conf.GetString("b"))
	assert.Equal(t, 1, conf.GetInt("c.d"))
	assert.Equal(t, "string", conf.GetString("c.e"))
}
