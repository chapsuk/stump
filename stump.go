package stump

import (
	"github.com/m1ome/stump/lib"
)

func Setup() (*lib.Stump, error) {
	return lib.New(nil)
}

func MustSetup() *lib.Stump {
	s, err := Setup()
	if err != nil {
		panic(err)
	}

	return s
}