package daap

import (
	"os"
	"testing"

	"github.com/otiai10/dkmachine/v0/dkmachine"
	. "github.com/otiai10/mint"
	"github.com/otiai10/ternary"
)

var testmachine *MachineConfig

func TestMain(m *testing.M) {
	machine, err := dkmachine.Create(&dkmachine.CreateOptions{
		Driver: ternary.String(os.Getenv("DRIVER"))("virtualbox"),
	})
	if err != nil {
		panic(err)
	}
	testmachine = &MachineConfig{
		Host:     machine.Host(),
		CertPath: machine.CertPath(),
	}
	code := m.Run()
	if err := machine.Remove(); err != nil {
		panic(err)
	}
	os.Exit(code)
}

func TestNewContainer(t *testing.T) {
	container := NewContainer("debian:latest", Args{Machine: testmachine})
	Expect(t, container).TypeOf("*daap.Container")
	Expect(t, container.Args.Machine.CertPath).ToBe(testmachine.CertPath)
}
