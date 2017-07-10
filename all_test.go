package daap

import (
	"context"
	"io/ioutil"
	"testing"

	. "github.com/otiai10/mint"
)

func TestProcess(t *testing.T) {
	proc := NewProcess("foo", nil)
	Expect(t, proc).TypeOf("*daap.Process")
}

func TestProcess_Run(t *testing.T) {
	proc := NewProcess("otiai10/foo", &Args{
		// Machine: &MachineConfig{
		// 	Host:     "tcp://192.168.99.100:2376",
		// 	CertPath: "/Users/otiai10/.docker/machine/machines/example",
		// },
		Machine: NewEnvMachine(),
	})
	ctx := context.Background()
	err := proc.Run(ctx)
	Expect(t, err).ToBe(nil)

	b, err := ioutil.ReadAll(proc.Stdout)
	Expect(t, err).ToBe(nil)
	Expect(t, string(b)).ToBe("hogeee")
}
