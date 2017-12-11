package daap

// Args represents argument for the process, representing machine (where to run), input (what to mount), output (where to output).
type Args struct {
	Machine *MachineConfig
	Env     []string
	Mounts  []Mount
	Name    string
}
