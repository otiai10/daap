package daap

import (
	"context"
	"testing"

	. "github.com/otiai10/mint"
)

func TestNewContainer(t *testing.T) {
	container := NewContainer("debian:latest", Args{
		Machine: NewEnvMachine(),
	})
	Expect(t, container).TypeOf("*daap.Container")
}

func TestContainer_PullImage(t *testing.T) {
	container := NewContainer("debian:latest", Args{
		Machine: NewEnvMachine(),
	})
	pulllog, err := container.PullImage(context.Background())
	Expect(t, err).ToBe(nil)
	for log := range pulllog {
		Expect(t, log).TypeOf("daap.ImagePullResponsePayload")
	}
}
