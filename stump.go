package stump

import (
	"github.com/m1ome/stump/core"
)

func Setup() (*core.Stump, error) {
	return core.New(nil)
}

func MustSetup() *core.Stump {
	s, err := Setup()
	if err != nil {
		panic(err)
	}

	return s
}