package daap

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/docker/docker/api/types/mount"
	. "github.com/otiai10/mint"
)

func TestProcess(t *testing.T) {
	proc := NewProcess("foo", Args{})
	Expect(t, proc).TypeOf("*daap.Process")
}

func TestProcess_Run(t *testing.T) {

	image := "otiai10/daap-test"
	args := Args{
		Machine: NewEnvMachine(),
		Env: []string{
			"DAAPTEST_FOO=FrostyJack",
			"DAAPTEST_BAR=BabyBlue",
		},
		Mounts: []mount.Mount{mount.Mount{Type: mount.TypeBind, Source: "/home/docker/data", Target: "/var/data"}},
	}

	proc := NewProcess(image, args)
	ctx := context.Background()
	err := proc.Run(ctx)
	Expect(t, err).ToBe(nil)

	b, err := ioutil.ReadAll(proc.Stdout)
	Expect(t, err).ToBe(nil)
	Expect(t, string(b)).ToBe("STDOUT: 'main.sh' started.STDOUT: 'main.sh' successfully finished.")
}
